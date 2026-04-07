package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"video-service/models"
	"video-service/providers/events"
	"video-service/providers/storage"
	"video-service/providers/transcoder"
)

// VideoHandler holds dependencies for video service endpoints.
type VideoHandler struct {
	Storage    storage.StorageProvider
	Transcoder transcoder.TranscoderProvider
	Publisher  events.EventPublisher
	Config     *HandlerConfig
}

// HandlerConfig provides configuration to handlers.
type HandlerConfig struct {
	GCSBucket       string
	GCSOutputBucket string
	CDNBaseURL      string
}

// HandleCreateMedia initializes a media entry in the video service.
func (h *VideoHandler) HandleCreateMedia(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		MediaID  string `json:"media_id" binding:"required"`
		Title    string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id and media_id are required"})
		return
	}

	status := &models.MediaStatus{
		MediaID:    req.MediaID,
		TenantID:   req.TenantID,
		Status:     "created",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	c.JSON(http.StatusCreated, status)
}

// HandleUpload generates a signed upload URL for direct client upload.
func (h *VideoHandler) HandleUpload(c *gin.Context) {
	mediaID := c.Param("id")

	var req models.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.MediaID = mediaID

	// Generate GCS object path
	objectPath := fmt.Sprintf("uploads/%s/%s/%s", req.TenantID, req.MediaID, req.FileName)
	gcsURI := fmt.Sprintf("gs://%s/%s", h.Config.GCSBucket, objectPath)

	// Generate signed upload URL
	uploadURL, err := h.Storage.GenerateUploadURL(h.Config.GCSBucket, objectPath, req.ContentType)
	if err != nil {
		log.Printf("Failed to generate upload URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, &models.UploadResponse{
		UploadURL:   uploadURL,
		GCSInputURI: gcsURI,
		ExpiresAt:   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	})
}

// HandleProcess starts a transcoding job for the media.
func (h *VideoHandler) HandleProcess(c *gin.Context) {
	mediaID := c.Param("id")

	var req models.ProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.MediaID = mediaID

	inputURI := fmt.Sprintf("gs://%s/uploads/%s/%s", h.Config.GCSBucket, req.TenantID, req.MediaID)
	outputPrefix := fmt.Sprintf("gs://%s/outputs/%s/%s/", h.Config.GCSOutputBucket, req.TenantID, req.MediaID)

	job, err := h.Transcoder.CreateJob(transcoder.TranscodeInput{
		InputURI:     inputURI,
		OutputPrefix: outputPrefix,
		MediaID:      req.MediaID,
		TenantID:     req.TenantID,
		EnableHLS:    true,
		EnableMP4:    true,
		Resolutions:  []string{"1080p", "720p", "480p"},
	})
	if err != nil {
		log.Printf("Failed to create transcoding job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start processing"})
		return
	}

	// Notify Sentanyl of processing start
	h.Publisher.PublishProcessingUpdate(events.ProcessingUpdate{
		TenantID:         req.TenantID,
		MediaPublicId:    req.MediaID,
		ProcessingStatus: "submitted",
		Status:           "processing",
	})

	c.JSON(http.StatusAccepted, gin.H{
		"status":    "processing",
		"job_name":  job.JobName,
		"input_uri": inputURI,
	})
}

// HandleGetStatus returns the current processing status of a media item.
func (h *VideoHandler) HandleGetStatus(c *gin.Context) {
	mediaID := c.Param("id")

	// In production, this would look up the job status from Transcoder API
	// and/or a local state store

	c.JSON(http.StatusOK, &models.MediaStatus{
		MediaID:          mediaID,
		Status:           "ready",
		ProcessingStatus: "succeeded",
		UpdatedAt:        time.Now(),
	})
}

// HandleGetPlayback returns playback information for a media item.
func (h *VideoHandler) HandleGetPlayback(c *gin.Context) {
	mediaID := c.Param("id")
	tenantID := c.Query("tenant_id")

	// Build playback URLs from CDN/storage
	baseURL := h.Config.CDNBaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://storage.googleapis.com/%s", h.Config.GCSOutputBucket)
	}

	outputBase := fmt.Sprintf("%s/outputs/%s/%s", baseURL, tenantID, mediaID)

	playback := &models.PlaybackResponse{
		MediaID: mediaID,
		Status:  "ready",
		Sources: []models.PlaybackSource{
			{
				Type:  "application/x-mpegURL",
				URL:   fmt.Sprintf("%s/manifest.m3u8", outputBase),
				Label: "Auto",
			},
			{
				Type:  "video/mp4",
				URL:   fmt.Sprintf("%s/output.mp4", outputBase),
				Label: "1080p",
			},
		},
		Poster: fmt.Sprintf("%s/poster.jpg", outputBase),
	}

	c.JSON(http.StatusOK, playback)
}

// HandleEvent ingests a player event and forwards to Sentanyl.
func (h *VideoHandler) HandleEvent(c *gin.Context) {
	var req models.PlayerEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event"})
		return
	}

	// Forward normalized event to Sentanyl control plane
	err := h.Publisher.PublishEvent(events.VideoEvent{
		TenantID:        req.TenantID,
		MediaPublicId:   req.MediaID,
		ViewerPublicId:  req.ViewerID,
		SessionPublicId: req.SessionID,
		EventName:       req.EventName,
		CurrentSecond:   req.CurrentSecond,
		ProgressPercent: req.ProgressPercent,
		PageURL:         req.PageURL,
		Domain:          req.Domain,
		Data:            req.Data,
		OccurredAt:      time.Now(),
	})

	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		// Don't fail the request — event ingest should be fire-and-forget for the player
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// HandleIdentify links a viewer session to an identity (email, external ID).
func (h *VideoHandler) HandleIdentify(c *gin.Context) {
	var req models.IdentifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identify request"})
		return
	}

	// Forward identify event to Sentanyl
	h.Publisher.PublishEvent(events.VideoEvent{
		TenantID:        req.TenantID,
		SessionPublicId: req.SessionID,
		EventName:       "identify",
		Data: map[string]interface{}{
			"email":       req.Email,
			"external_id": req.ExternalID,
			"source":      req.Source,
		},
		OccurredAt: time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"status": "identified"})
}

// HandleTranscoderCallback handles Transcoder API job completion callbacks.
func (h *VideoHandler) HandleTranscoderCallback(c *gin.Context) {
	var req models.TranscoderCallback
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid callback"})
		return
	}

	baseURL := h.Config.CDNBaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://storage.googleapis.com/%s", h.Config.GCSOutputBucket)
	}

	update := events.ProcessingUpdate{
		TenantID:      req.TenantID,
		MediaPublicId: req.MediaID,
	}

	switch req.State {
	case "SUCCEEDED":
		update.ProcessingStatus = "succeeded"
		update.Status = "ready"
		outputBase := fmt.Sprintf("%s/outputs/%s/%s", baseURL, req.TenantID, req.MediaID)
		update.PlaybackHLSURL = fmt.Sprintf("%s/manifest.m3u8", outputBase)
		update.PlaybackMP4URL = fmt.Sprintf("%s/output.mp4", outputBase)
		update.PosterURL = fmt.Sprintf("%s/poster.jpg", outputBase)
	case "FAILED":
		update.ProcessingStatus = "failed"
		update.Status = "failed"
		update.ErrorMessage = "Transcoding job failed"
	default:
		update.ProcessingStatus = "running"
		update.Status = "processing"
	}

	h.Publisher.PublishProcessingUpdate(update)

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}
