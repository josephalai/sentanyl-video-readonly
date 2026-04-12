package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/josephalai/sentanyl/video-service/handlers"
	"github.com/josephalai/sentanyl/video-service/queries"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	videoQueries := queries.NewVideoQueries()

	intelHandler := &handlers.IntelligenceHandler{
		Queries: videoQueries,
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "video-service"})
	})

	// Video intelligence routes (creator-facing) under /api/video/*
	// matching the Caddy gateway prefix.
	video := r.Group("/api/video")
	{
		video.GET("/media", intelHandler.HandleListMedia)
		video.POST("/media", intelHandler.HandleCreateMediaIntel)
		video.GET("/media/:id", intelHandler.HandleGetMedia)
		video.PUT("/media/:id", intelHandler.HandleUpdateMedia)
		video.DELETE("/media/:id", intelHandler.HandleDeleteMedia)

		video.GET("/presets", intelHandler.HandleListPlayerPresets)
		video.POST("/presets", intelHandler.HandleCreatePlayerPreset)

		video.GET("/channels", intelHandler.HandleListChannels)
		video.POST("/channels", intelHandler.HandleCreateChannel)

		video.GET("/webhooks", intelHandler.HandleListWebhooks)
		video.POST("/webhooks", intelHandler.HandleCreateWebhook)

		video.GET("/analytics/overview", intelHandler.HandleGetAnalyticsOverview)
		video.GET("/media/:id/viewers", intelHandler.HandleListViewersByMedia)
		video.GET("/media/:id/sessions", intelHandler.HandleListSessionsByMedia)
	}

	// Tracking
	r.GET("/api/track/click/:token", handlers.HandleClickTracking)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Video Service (intelligence) starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
