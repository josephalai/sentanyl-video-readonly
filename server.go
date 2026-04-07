package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"video-service/handlers"
	"video-service/providers/events"
	"video-service/providers/storage"
	"video-service/providers/transcoder"
)

func main() {
	cfg := LoadConfig()

	// Initialize providers
	storageProvider := storage.NewGCSProvider(cfg.GCPProjectID)
	transcoderProvider := transcoder.NewGCPTranscoderProvider(cfg.GCPProjectID, cfg.TranscoderRegion)

	// Initialize event publisher
	// Use HTTP callback by default; switch to Pub/Sub via config
	var publisher events.EventPublisher
	publisher = events.NewHTTPCallbackPublisher(cfg.SentanylBaseURL, cfg.SentanylCallbackKey)

	// Create handler with dependencies
	handler := &handlers.VideoHandler{
		Storage:    storageProvider,
		Transcoder: transcoderProvider,
		Publisher:  publisher,
		Config: &handlers.HandlerConfig{
			GCSBucket:       cfg.GCSBucket,
			GCSOutputBucket: cfg.GCSOutputBucket,
			CDNBaseURL:      cfg.CDNBaseURL,
		},
	}

	// Setup Gin router
	r := gin.Default()
	RegisterRoutes(r, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Video Service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
