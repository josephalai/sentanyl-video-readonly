package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/josephalai/sentanyl/video-service/handlers"
	"github.com/josephalai/sentanyl/video-service/providers/storage"
	"github.com/josephalai/sentanyl/video-service/queries"
	"github.com/josephalai/sentanyl/pkg/auth"
	"github.com/josephalai/sentanyl/pkg/db"
	httputil "github.com/josephalai/sentanyl/pkg/http"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Initialize MongoDB connection.
	db.MongoHost = envOrDefault("MONGO_HOST", "localhost")
	db.MongoPort = envOrDefault("MONGO_PORT", "27017")
	db.MongoDB = envOrDefault("MONGO_DB", "sentanyl")
	db.MongoDefaultCollectionName = "media"
	db.UsingLocalMongo = true
	db.InitMongoConnection()

	videoQueries := queries.NewVideoQueries()

	intelHandler := &handlers.IntelligenceHandler{
		Queries: videoQueries,
	}

	// Initialize object storage for tenant asset uploads.
	// If credentials are missing the service still starts — upload handler 503s.
	bucket := envOrDefault("GCS_BUCKET", "sendhero-videos")
	projectID := envOrDefault("GCS_PROJECT_ID", envOrDefault("GCP_PROJECT_ID", "sendhero"))
	assetsHandler := &handlers.AssetsHandler{Bucket: bucket}
	if storageProvider, serr := storage.NewGCSProvider(projectID); serr != nil {
		log.Printf("WARN: GCS storage provider init failed, asset uploads disabled: %v", serr)
	} else {
		defer storageProvider.Close()
		assetsHandler.Storage = storageProvider
	}

	r := gin.Default()
	r.Use(httputil.CORSMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "video-service"})
	})

	// Video intelligence routes (creator-facing) under /api/video/*
	// matching the Caddy gateway prefix.
	// RequireTenantAuth validates the JWT; the bridge middleware then forwards
	// the tenant_id from the gin context into the X-Tenant-ID header so the
	// existing handler code (which reads that header) continues to work unchanged.
	video := r.Group("/api/video")
	video.Use(auth.RequireTenantAuth())
	video.Use(func(c *gin.Context) {
		if tid := auth.GetTenantID(c); tid != "" {
			c.Request.Header.Set("X-Tenant-ID", tid)
		}
		c.Next()
	})
	{
		video.GET("/media", intelHandler.HandleListMedia)
		video.POST("/media", intelHandler.HandleCreateMediaIntel)
		video.GET("/media/:id", intelHandler.HandleGetMedia)
		video.PUT("/media/:id", intelHandler.HandleUpdateMedia)
		video.DELETE("/media/:id", intelHandler.HandleDeleteMedia)

		video.GET("/presets", intelHandler.HandleListPlayerPresets)
		video.POST("/presets", intelHandler.HandleCreatePlayerPreset)
		video.PUT("/presets/:id", intelHandler.HandleUpdatePlayerPreset)
		video.DELETE("/presets/:id", intelHandler.HandleDeletePlayerPreset)

		video.POST("/assets/upload", assetsHandler.HandleUpload)

		video.GET("/channels", intelHandler.HandleListChannels)
		video.POST("/channels", intelHandler.HandleCreateChannel)

		video.GET("/webhooks", intelHandler.HandleListWebhooks)
		video.POST("/webhooks", intelHandler.HandleCreateWebhook)

		video.GET("/analytics/overview", intelHandler.HandleGetAnalyticsOverview)
		video.GET("/media/:id/viewers", intelHandler.HandleListViewersByMedia)
		video.GET("/media/:id/sessions", intelHandler.HandleListSessionsByMedia)
	}

	// Player event ingestion — no auth, called from visitor browsers.
	// Returns viewer_id, session_id, and any badge_public_ids to grant.
	r.POST("/api/video/events", intelHandler.HandleIngestEvent)

	// Tracking
	r.GET("/api/track/click/:token", handlers.HandleClickTracking)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Video Service (intelligence) starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
