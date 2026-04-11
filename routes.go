package main

import (
	"github.com/gin-gonic/gin"

	"github.com/josephalai/sentanyl/video-service/handlers"
)

// RegisterRoutes sets up all Video Service routes.
func RegisterRoutes(r *gin.Engine, h *handlers.VideoHandler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "video-service"})
	})

	// Media lifecycle
	r.POST("/media", h.HandleCreateMedia)
	r.POST("/media/:id/upload", h.HandleUpload)
	r.POST("/media/:id/process", h.HandleProcess)
	r.GET("/media/:id/status", h.HandleGetStatus)

	// Playback
	r.GET("/media/:id/playback", h.HandleGetPlayback)

	// Event ingest
	r.POST("/event", h.HandleEvent)

	// Viewer identification
	r.POST("/identify", h.HandleIdentify)

	// Transcoder callback
	r.POST("/callbacks/transcoder", h.HandleTranscoderCallback)
}
