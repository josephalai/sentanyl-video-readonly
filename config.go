package main

import (
	"fmt"
	"log"
	"os"
)

// Config holds configuration for the Video Service.
type Config struct {
	// Server
	Port string

	// Google Cloud
	GCPProjectID     string
	GCSBucket        string
	GCSOutputBucket  string
	TranscoderRegion string

	// Pub/Sub
	PubSubTopic string

	// Sentanyl callback
	SentanylBaseURL    string
	SentanylCallbackKey string

	// CDN / Playback
	CDNBaseURL string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	cfg := &Config{
		Port:             getEnv("PORT", "8090"),
		GCPProjectID:     getEnv("GCP_PROJECT_ID", "sendhero"),
		GCSBucket:        getEnv("GCS_BUCKET", "sendhero-videos"),
		GCSOutputBucket:  getEnv("GCS_OUTPUT_BUCKET", "sendhero-videos"),
		TranscoderRegion: getEnv("TRANSCODER_REGION", "us-central1"),
		PubSubTopic:      getEnv("PUBSUB_TOPIC", "video-events"),
		SentanylBaseURL:  getEnv("SENTANYL_BASE_URL", "http://localhost:8080"),
		SentanylCallbackKey: getEnv("SENTANYL_CALLBACK_KEY", ""),
		CDNBaseURL:       getEnv("CDN_BASE_URL", "https://storage.googleapis.com/sendhero-videos"),
	}

	log.Printf("Video Service Config: project=%s bucket=%s region=%s",
		cfg.GCPProjectID, cfg.GCSBucket, cfg.TranscoderRegion)

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func (c *Config) GCSInputPrefix() string {
	return fmt.Sprintf("gs://%s/uploads", c.GCSBucket)
}

func (c *Config) GCSOutputPrefix() string {
	return fmt.Sprintf("gs://%s/outputs", c.GCSOutputBucket)
}

func (c *Config) PlaybackBaseURL() string {
	if c.CDNBaseURL != "" {
		return c.CDNBaseURL
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s", c.GCSOutputBucket)
}
