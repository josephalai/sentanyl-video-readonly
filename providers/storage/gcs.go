package storage

import (
	"fmt"
	"log"
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
}

// GCSProvider implements StorageProvider for Google Cloud Storage.
type GCSProvider struct {
	ProjectID string
}

// NewGCSProvider creates a new GCS storage provider.
func NewGCSProvider(projectID string) *GCSProvider {
	return &GCSProvider{ProjectID: projectID}
}

func (g *GCSProvider) GenerateUploadURL(bucket, objectPath, contentType string) (string, error) {
	// TODO: Implement with cloud.google.com/go/storage SignedURL
	// For now, return a constructed public URL (bucket is public)
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectPath)
	log.Printf("GCS: Generated upload URL for %s/%s", bucket, objectPath)
	return url, nil
}

func (g *GCSProvider) GeneratePlaybackURL(bucket, objectPath string) (string, error) {
	// For public buckets, use direct URL
	// For private buckets, this would generate a signed URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectPath)
	return url, nil
}

func (g *GCSProvider) ObjectExists(bucket, objectPath string) (bool, error) {
	// TODO: Implement with GCS client
	log.Printf("GCS: Checking existence of %s/%s", bucket, objectPath)
	return false, nil
}

func (g *GCSProvider) DeleteObject(bucket, objectPath string) error {
	// TODO: Implement with GCS client
	log.Printf("GCS: Deleting %s/%s", bucket, objectPath)
	return nil
}
