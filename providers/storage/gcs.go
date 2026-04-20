package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

// StorageProvider abstracts cloud object storage operations.
type StorageProvider interface {
	// GenerateUploadURL returns a signed URL for direct upload.
	GenerateUploadURL(bucket, objectPath, contentType string) (signedURL string, err error)

	// GeneratePlaybackURL returns a (potentially signed) URL for playback.
	GeneratePlaybackURL(bucket, objectPath string) (string, error)

	// ObjectExists checks if an object exists in the bucket.
	ObjectExists(bucket, objectPath string) (bool, error)

	// DeleteObject removes an object from storage.
	DeleteObject(bucket, objectPath string) error

	// UploadObject uploads data to the specified object path.
	UploadObject(bucket, objectPath, contentType string, data io.Reader) (string, error)
}

// GCSProvider implements StorageProvider for Google Cloud Storage.
type GCSProvider struct {
	ProjectID string
	client    *storage.Client
}

// NewGCSProvider creates a new GCS storage provider.
// Uses Application Default Credentials (ADC) for authentication.
//
// The client forces HTTP/1.1 because some network paths (certain residential
// ISPs, dev containers behind NAT) reset long-lived HTTP/2 streams mid-upload
// with "http2: client connection lost", which the GCS SDK cannot retry past.
func NewGCSProvider(projectID string) (*GCSProvider, error) {
	ctx := context.Background()

	// Base transport with HTTP/2 disabled.
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.ForceAttemptHTTP2 = false
	base.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)

	// Wrap with Google auth (ADC) so the HTTP client carries bearer tokens.
	authed, err := htransport.NewTransport(ctx, base, option.WithScopes(storage.ScopeFullControl))
	if err != nil {
		return nil, fmt.Errorf("failed to build authed GCS transport: %w", err)
	}

	httpClient := &http.Client{Transport: authed}

	client, err := storage.NewClient(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return &GCSProvider{ProjectID: projectID, client: client}, nil
}

// Close releases the GCS client resources.
func (g *GCSProvider) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}

func (g *GCSProvider) GenerateUploadURL(bucket, objectPath, contentType string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:      storage.SigningSchemeV4,
		Method:      "PUT",
		ContentType: contentType,
		Expires:     time.Now().Add(15 * time.Minute),
	}

	signedURL, err := g.client.Bucket(bucket).SignedURL(objectPath, opts)
	if err != nil {
		// Fall back to public URL if signing fails (e.g., public bucket without service account)
		log.Printf("[GCS] SignedURL failed, using public URL: %v", err)
		return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectPath), nil
	}

	log.Printf("[GCS] Generated signed upload URL for %s/%s", bucket, objectPath)
	return signedURL, nil
}

func (g *GCSProvider) GeneratePlaybackURL(bucket, objectPath string) (string, error) {
	ctx := context.Background()

	// Check if the object exists first
	_, err := g.client.Bucket(bucket).Object(objectPath).Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("object %s/%s not found: %w", bucket, objectPath, err)
	}

	// For public buckets, return direct URL
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectPath), nil
}

func (g *GCSProvider) ObjectExists(bucket, objectPath string) (bool, error) {
	ctx := context.Background()

	_, err := g.client.Bucket(bucket).Object(objectPath).Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

func (g *GCSProvider) DeleteObject(bucket, objectPath string) error {
	ctx := context.Background()

	obj := g.client.Bucket(bucket).Object(objectPath)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete %s/%s: %w", bucket, objectPath, err)
	}

	log.Printf("[GCS] Deleted %s/%s", bucket, objectPath)
	return nil
}

// UploadObject uploads data from a reader to the specified GCS object path.
// Returns the public URL of the uploaded object.
func (g *GCSProvider) UploadObject(bucket, objectPath, contentType string, data io.Reader) (string, error) {
	ctx := context.Background()

	wc := g.client.Bucket(bucket).Object(objectPath).NewWriter(ctx)
	wc.ContentType = contentType
	wc.CacheControl = "public, max-age=31536000"
	// Force chunked resumable uploads so transient HTTP/2 disconnects are retried
	// per-chunk instead of failing the whole upload. Default is 16 MiB single-shot
	// for files below that threshold, which is brittle on flaky networks.
	wc.ChunkSize = 4 << 20 // 4 MiB

	if _, err := io.Copy(wc, data); err != nil {
		wc.Close()
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize GCS upload: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectPath)
	log.Printf("[GCS] Uploaded %s (%s)", publicURL, contentType)
	return publicURL, nil
}
