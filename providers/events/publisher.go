package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// EventPublisher abstracts event publication for the Video Service.
// Events can be forwarded to Sentanyl via HTTP callback or Pub/Sub.
type EventPublisher interface {
	// PublishEvent sends a normalized video event to the control plane.
	PublishEvent(event VideoEvent) error

	// PublishProcessingUpdate sends a media processing status update.
	PublishProcessingUpdate(update ProcessingUpdate) error
}

// VideoEvent is a normalized event from the Video Service.
type VideoEvent struct {
	TenantID        string                 `json:"tenant_id"`
	MediaPublicId   string                 `json:"media_public_id"`
	ViewerPublicId  string                 `json:"viewer_public_id,omitempty"`
	SessionPublicId string                 `json:"session_public_id,omitempty"`
	EventName       string                 `json:"event_name"`
	CurrentSecond   int                    `json:"current_second"`
	ProgressPercent int                    `json:"progress_pct"`
	PageURL         string                 `json:"page_url,omitempty"`
	Domain          string                 `json:"domain,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	OccurredAt      time.Time              `json:"occurred_at"`
}

// ProcessingUpdate represents a media processing state change.
type ProcessingUpdate struct {
	TenantID         string `json:"tenant_id"`
	MediaPublicId    string `json:"media_public_id"`
	ProcessingStatus string `json:"processing_status"`
	Status           string `json:"status"`
	DurationSec      int64  `json:"duration_sec,omitempty"`
	PosterURL        string `json:"poster_url,omitempty"`
	PlaybackHLSURL   string `json:"playback_hls_url,omitempty"`
	PlaybackDASHURL  string `json:"playback_dash_url,omitempty"`
	PlaybackMP4URL   string `json:"playback_mp4_url,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

// HTTPCallbackPublisher implements EventPublisher via HTTP webhook to Sentanyl.
type HTTPCallbackPublisher struct {
	SentanylBaseURL string
	CallbackKey     string
	Client          *http.Client
}

// NewHTTPCallbackPublisher creates a new HTTP callback publisher.
func NewHTTPCallbackPublisher(baseURL, callbackKey string) *HTTPCallbackPublisher {
	return &HTTPCallbackPublisher{
		SentanylBaseURL: baseURL,
		CallbackKey:     callbackKey,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *HTTPCallbackPublisher) PublishEvent(event VideoEvent) error {
	url := fmt.Sprintf("%s/api/tenant/video/events", p.SentanylBaseURL)
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if p.CallbackKey != "" {
		req.Header.Set("Authorization", p.CallbackKey)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("HTTPCallback: Failed to publish event to %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sentanyl returned status %d", resp.StatusCode)
	}

	log.Printf("HTTPCallback: Published %s event for media %s", event.EventName, event.MediaPublicId)
	return nil
}

func (p *HTTPCallbackPublisher) PublishProcessingUpdate(update ProcessingUpdate) error {
	url := fmt.Sprintf("%s/api/internal/video/processing-update", p.SentanylBaseURL)
	body, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if p.CallbackKey != "" {
		req.Header.Set("Authorization", p.CallbackKey)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("HTTPCallback: Failed to publish processing update: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("HTTPCallback: Published processing update for media %s (status=%s)",
		update.MediaPublicId, update.ProcessingStatus)
	return nil
}

// PubSubPublisher implements EventPublisher via Google Cloud Pub/Sub.
// This is the recommended production path for async event processing.
type PubSubPublisher struct {
	ProjectID string
	TopicID   string
}

// NewPubSubPublisher creates a new Pub/Sub publisher.
func NewPubSubPublisher(projectID, topicID string) *PubSubPublisher {
	return &PubSubPublisher{
		ProjectID: projectID,
		TopicID:   topicID,
	}
}

func (p *PubSubPublisher) PublishEvent(event VideoEvent) error {
	// TODO: Implement with cloud.google.com/go/pubsub
	log.Printf("PubSub: Would publish %s event to topic %s", event.EventName, p.TopicID)
	return nil
}

func (p *PubSubPublisher) PublishProcessingUpdate(update ProcessingUpdate) error {
	// TODO: Implement with cloud.google.com/go/pubsub
	log.Printf("PubSub: Would publish processing update to topic %s", p.TopicID)
	return nil
}
