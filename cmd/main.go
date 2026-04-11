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

	// Intelligence routes (creator-facing)
	intel := r.Group("/api/v1/intel")
	{
		intel.GET("/media", intelHandler.HandleListMedia)
		intel.POST("/media", intelHandler.HandleCreateMediaIntel)
		intel.GET("/media/:id", intelHandler.HandleGetMedia)
		intel.PUT("/media/:id", intelHandler.HandleUpdateMedia)
		intel.DELETE("/media/:id", intelHandler.HandleDeleteMedia)

		intel.GET("/presets", intelHandler.HandleListPlayerPresets)
		intel.POST("/presets", intelHandler.HandleCreatePlayerPreset)

		intel.GET("/channels", intelHandler.HandleListChannels)
		intel.POST("/channels", intelHandler.HandleCreateChannel)

		intel.GET("/webhooks", intelHandler.HandleListWebhooks)
		intel.POST("/webhooks", intelHandler.HandleCreateWebhook)

		intel.GET("/analytics/overview", intelHandler.HandleGetAnalyticsOverview)
		intel.GET("/media/:id/viewers", intelHandler.HandleListViewersByMedia)
		intel.GET("/media/:id/sessions", intelHandler.HandleListSessionsByMedia)
	}

	// Tracking
	r.GET("/api/track/click/:token", handlers.HandleClickTracking)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Video Service (intelligence) starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
