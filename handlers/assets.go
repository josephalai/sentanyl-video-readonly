package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/josephalai/sentanyl/pkg/storage"
)

// AssetsHandler handles generic tenant asset uploads (media source files).
// Files are streamed to GCS under uploads/{tenantID}/{objectID}-{filename}
// and the returned public URL is used as the Media.source_url.
type AssetsHandler struct {
	Storage storage.StorageProvider
	Bucket  string
}

// HandleUpload accepts a multipart file upload and streams it to object storage.
func (h *AssetsHandler) HandleUpload(c *gin.Context) {
	if h.Storage == nil || h.Bucket == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object storage not configured"})
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read upload"})
		return
	}
	defer f.Close()

	safeName := strings.ReplaceAll(fh.Filename, " ", "_")
	ext := strings.ToLower(filepath.Ext(safeName))
	objectPath := fmt.Sprintf("uploads/%s/%s-%s", tenantID, bson.NewObjectId().Hex(), safeName)

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mimeForExt(ext)
	}

	publicURL, err := h.Storage.UploadObject(h.Bucket, objectPath, contentType, f)
	if err != nil {
		log.Printf("HandleUploadAsset: GCS upload failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":          publicURL,
		"file_url":     publicURL,
		"object_path":  objectPath,
		"content_type": contentType,
	})
}

func mimeForExt(ext string) string {
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	}
	return "application/octet-stream"
}
