package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/josephalai/sentanyl/video-service/handlers"
	"github.com/josephalai/sentanyl/video-service/providers/events"
	"github.com/josephalai/sentanyl/pkg/storage"
	"github.com/josephalai/sentanyl/video-service/providers/transcoder"
)

func main() {
	cfg := LoadConfig()

	// Initialize providers
	storageProvider, err := storage.NewGCSProvider(cfg.GCPProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize GCS storage provider: %v", err)
	}
	defer storageProvider.Close()

	transcoderProvider, err := transcoder.NewGCPTranscoderProvider(cfg.GCPProjectID, cfg.TranscoderRegion)
	if err != nil {
		log.Fatalf("Failed to initialize GCP Transcoder provider: %v", err)
	}
	defer transcoderProvider.Close()

	// Initialize event publisher
	// Use HTTP callback by default; switch to Pub/Sub via config
	var publisher events.EventPublisher
	if cfg.PubSubTopic != "" {
		pubsubPublisher, pubErr := events.NewPubSubPublisher(cfg.GCPProjectID, cfg.PubSubTopic)
		if pubErr != nil {
			log.Printf("Failed to initialize Pub/Sub publisher, falling back to HTTP callback: %v", pubErr)
			publisher = events.NewHTTPCallbackPublisher(cfg.SentanylBaseURL, cfg.SentanylCallbackKey)
		} else {
			defer pubsubPublisher.Close()
			publisher = pubsubPublisher
		}
	} else {
		publisher = events.NewHTTPCallbackPublisher(cfg.SentanylBaseURL, cfg.SentanylCallbackKey)
	}

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
