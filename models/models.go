package models

import "time"

// MediaStatus represents the processing state of a media item.
type MediaStatus struct {
	MediaID          string    `json:"media_id"`
	TenantID         string    `json:"tenant_id"`
	Status           string    `json:"status"` // "pending", "uploading", "processing", "ready", "failed"
	ProcessingStatus string    `json:"processing_status"`
	GCSInputURI      string    `json:"gcs_input_uri"`
	GCSOutputPrefix  string    `json:"gcs_output_prefix"`
	TranscoderJob    string    `json:"transcoder_job_name"`
	DurationSec      int64     `json:"duration_sec"`
	PosterURL        string    `json:"poster_url"`
	PlaybackHLSURL   string    `json:"playback_hls_url"`
	PlaybackDASHURL  string    `json:"playback_dash_url"`
	PlaybackMP4URL   string    `json:"playback_mp4_url"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PlaybackResponse is returned to players requesting media playback info.
type PlaybackResponse struct {
	MediaID    string            `json:"media_id"`
	Status     string            `json:"status"`
	Sources    []PlaybackSource  `json:"sources"`
	Poster     string            `json:"poster,omitempty"`
	Duration   int64             `json:"duration_sec,omitempty"`
	Captions   []CaptionTrack    `json:"captions,omitempty"`
}

// PlaybackSource represents a single playback variant.
type PlaybackSource struct {
	Type string `json:"type"` // "application/x-mpegURL", "video/mp4"
	URL  string `json:"url"`
	Label string `json:"label,omitempty"` // "1080p", "720p"
}

// CaptionTrack for playback.
type CaptionTrack struct {
	Language string `json:"language"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	Default  bool   `json:"default"`
}

// PlayerEvent represents a raw event from the video player.
type PlayerEvent struct {
	TenantID        string                 `json:"tenant_id"`
	MediaID         string                 `json:"media_id"`
	ViewerID        string                 `json:"viewer_id,omitempty"`
	SessionID       string                 `json:"session_id,omitempty"`
	EventName       string                 `json:"event_name"` // play, pause, progress, complete, etc.
	CurrentSecond   int                    `json:"current_second"`
	ProgressPercent int                    `json:"progress_percent"`
	PageURL         string                 `json:"page_url,omitempty"`
	Domain          string                 `json:"domain,omitempty"`
	Referrer        string                 `json:"referrer,omitempty"`
	DeviceType      string                 `json:"device_type,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
}

// IdentifyRequest links a viewer session to an identity.
type IdentifyRequest struct {
	TenantID   string `json:"tenant_id"`
	SessionID  string `json:"session_id"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Source     string `json:"source"` // "turnstile", "logged_in", "url_param", "api"
}

// TranscoderCallback represents data from Transcoder API completion.
type TranscoderCallback struct {
	JobName  string `json:"job_name"`
	State    string `json:"state"` // "SUCCEEDED", "FAILED"
	MediaID  string `json:"media_id"`
	TenantID string `json:"tenant_id"`
}

// UploadRequest for initiating media upload.
type UploadRequest struct {
	TenantID    string `json:"tenant_id"`
	MediaID     string `json:"media_id"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
}

// UploadResponse returns signed upload URLs.
type UploadResponse struct {
	UploadURL    string `json:"upload_url"`
	GCSInputURI  string `json:"gcs_input_uri"`
	ExpiresAt    string `json:"expires_at"`
}

// ProcessRequest to start transcoding.
type ProcessRequest struct {
	TenantID string `json:"tenant_id"`
	MediaID  string `json:"media_id"`
}
