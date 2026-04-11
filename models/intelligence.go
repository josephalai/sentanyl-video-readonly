package models

import "time"

// Collection name constants for video intelligence entities.
const (
	MediaCollectionName               = "media"
	PlayerPresetCollectionName        = "player_presets"
	MediaChannelCollectionName        = "media_channels"
	MediaWebhookCollectionName        = "media_webhooks"
	ViewerIdentityCollectionName      = "viewer_identities"
	ViewingSessionCollectionName      = "viewing_sessions"
	MediaEventCollectionName          = "media_events"
	MediaLeadCaptureCollectionName    = "media_lead_captures"
	MediaDailyAggregateCollectionName = "media_daily_aggregates"
)

// Video event name constants.
const (
	VideoEventPlay            = "play"
	VideoEventPause           = "pause"
	VideoEventProgress        = "progress"
	VideoEventComplete        = "complete"
	VideoEventTurnstileSubmit = "turnstile_submit"
	VideoEventCTAClick        = "cta_click"
	VideoEventChapterClick    = "chapter_click"
	VideoEventIdentify        = "identify"
	VideoEventRewatch         = "rewatch"
)

// Media represents a first-class video/audio asset in the platform.
type Media struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	PublicId     string `json:"public_id,omitempty" bson:"public_id"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id"`

	Title       string `json:"title,omitempty" bson:"title"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Kind        string `json:"kind,omitempty" bson:"kind,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`

	SourceURL string `json:"source_url,omitempty" bson:"source_url,omitempty"`
	PosterURL string `json:"poster_url,omitempty" bson:"poster_url,omitempty"`

	DurationSec   int64  `json:"duration_sec,omitempty" bson:"duration_sec,omitempty"`
	FileSizeBytes int64  `json:"file_size_bytes,omitempty" bson:"file_size_bytes,omitempty"`
	MimeType      string `json:"mime_type,omitempty" bson:"mime_type,omitempty"`

	ProcessingProvider string `json:"processing_provider,omitempty" bson:"processing_provider,omitempty"`
	ProcessingStatus   string `json:"processing_status,omitempty" bson:"processing_status,omitempty"`
	GCSInputURI        string `json:"gcs_input_uri,omitempty" bson:"gcs_input_uri,omitempty"`
	GCSOutputPrefix    string `json:"gcs_output_prefix,omitempty" bson:"gcs_output_prefix,omitempty"`
	TranscoderJobName  string `json:"transcoder_job_name,omitempty" bson:"transcoder_job_name,omitempty"`

	PlaybackHLSURL  string `json:"playback_hls_url,omitempty" bson:"playback_hls_url,omitempty"`
	PlaybackDASHURL string `json:"playback_dash_url,omitempty" bson:"playback_dash_url,omitempty"`
	PlaybackMP4URL  string `json:"playback_mp4_url,omitempty" bson:"playback_mp4_url,omitempty"`

	Renditions   []*MediaRendition    `json:"renditions,omitempty" bson:"renditions,omitempty"`
	Transcript   *MediaTranscript     `json:"transcript,omitempty" bson:"transcript,omitempty"`
	Captions     []*MediaCaptionTrack `json:"captions,omitempty" bson:"captions,omitempty"`
	Chapters     []*MediaChapter      `json:"chapters,omitempty" bson:"chapters,omitempty"`
	Interactions []*MediaInteraction   `json:"interactions,omitempty" bson:"interactions,omitempty"`
	BadgeRules   []*MediaBadgeRule    `json:"badge_rules,omitempty" bson:"badge_rules,omitempty"`
	Access       *MediaAccessPolicy   `json:"access,omitempty" bson:"access,omitempty"`

	PlayerPresetID string `json:"player_preset_id,omitempty" bson:"player_preset_id,omitempty"`

	Tags   []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Folder string   `json:"folder,omitempty" bson:"folder,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// MediaRendition represents a single encoded variant.
type MediaRendition struct {
	Kind        string `json:"kind,omitempty" bson:"kind"`
	Label       string `json:"label,omitempty" bson:"label,omitempty"`
	URL         string `json:"url,omitempty" bson:"url"`
	Width       int    `json:"width,omitempty" bson:"width,omitempty"`
	Height      int    `json:"height,omitempty" bson:"height,omitempty"`
	BitrateKbps int    `json:"bitrate_kbps,omitempty" bson:"bitrate_kbps,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty" bson:"is_default,omitempty"`
}

// MediaTranscript holds transcript data.
type MediaTranscript struct {
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
	Text      string `json:"text,omitempty" bson:"text,omitempty"`
	Generated bool   `json:"generated,omitempty" bson:"generated,omitempty"`
}

// MediaCaptionTrack represents a caption/subtitle track.
type MediaCaptionTrack struct {
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	Label     string `json:"label,omitempty" bson:"label,omitempty"`
	URL       string `json:"url,omitempty" bson:"url,omitempty"`
	IsDefault bool   `json:"is_default,omitempty" bson:"is_default,omitempty"`
}

// MediaChapter represents a navigable chapter within the media.
type MediaChapter struct {
	PublicId string `json:"public_id,omitempty" bson:"public_id"`
	Title    string `json:"title,omitempty" bson:"title"`
	StartSec int    `json:"start_sec,omitempty" bson:"start_sec"`
	EndSec   int    `json:"end_sec,omitempty" bson:"end_sec,omitempty"`
}

// MediaInteraction represents an interactive element overlaid on the media.
type MediaInteraction struct {
	PublicId string      `json:"public_id,omitempty" bson:"public_id"`
	Kind     string      `json:"kind,omitempty" bson:"kind"`
	StartSec int         `json:"start_sec,omitempty" bson:"start_sec"`
	EndSec   int         `json:"end_sec,omitempty" bson:"end_sec,omitempty"`
	Status   string      `json:"status,omitempty" bson:"status,omitempty"`
	Config   interface{} `json:"config,omitempty" bson:"config,omitempty"`
}

// MediaBadgeRule connects media events to badge grants.
type MediaBadgeRule struct {
	PublicId        string `json:"public_id,omitempty" bson:"public_id"`
	EventName       string `json:"event_name,omitempty" bson:"event_name"`
	Operator        string `json:"operator,omitempty" bson:"operator"`
	Threshold       int    `json:"threshold,omitempty" bson:"threshold,omitempty"`
	BadgePublicId   string `json:"badge_public_id,omitempty" bson:"badge_public_id"`
	AchievementText string `json:"achievement_text,omitempty" bson:"achievement_text,omitempty"`
	OncePerViewer   bool   `json:"once_per_viewer,omitempty" bson:"once_per_viewer,omitempty"`
	Enabled         bool   `json:"enabled,omitempty" bson:"enabled"`
}

// MediaAccessPolicy controls who can access the media.
type MediaAccessPolicy struct {
	PrivacyMode     string   `json:"privacy_mode,omitempty" bson:"privacy_mode,omitempty"`
	AllowedDomains  []string `json:"allowed_domains,omitempty" bson:"allowed_domains,omitempty"`
	DownloadEnabled bool     `json:"download_enabled,omitempty" bson:"download_enabled,omitempty"`
	ShareEnabled    bool     `json:"share_enabled,omitempty" bson:"share_enabled,omitempty"`
}

// PlayerPreset defines a reusable branded player configuration.
type PlayerPreset struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	PublicId     string `json:"public_id,omitempty" bson:"public_id"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id"`

	Name              string `json:"name,omitempty" bson:"name"`
	PlayerColor       string `json:"player_color,omitempty" bson:"player_color,omitempty"`
	ShowControls      bool   `json:"show_controls,omitempty" bson:"show_controls"`
	ShowRewind        bool   `json:"show_rewind,omitempty" bson:"show_rewind"`
	ShowFastForward   bool   `json:"show_fast_forward,omitempty" bson:"show_fast_forward"`
	ShowSkip          bool   `json:"show_skip,omitempty" bson:"show_skip"`
	ShowDownload      bool   `json:"show_download,omitempty" bson:"show_download"`
	HideProgressBar   bool   `json:"hide_progress_bar,omitempty" bson:"hide_progress_bar"`
	ShowBigPlayButton bool   `json:"show_big_play_button,omitempty" bson:"show_big_play_button"`
	AllowFullscreen   bool   `json:"allow_fullscreen,omitempty" bson:"allow_fullscreen"`
	AllowPlaybackRate bool   `json:"allow_playback_rate,omitempty" bson:"allow_playback_rate"`
	Autoplay          bool   `json:"autoplay,omitempty" bson:"autoplay"`
	MutedDefault      bool   `json:"muted_default,omitempty" bson:"muted_default"`
	DisablePause      bool   `json:"disable_pause,omitempty" bson:"disable_pause"`
	EndBehavior       string `json:"end_behavior,omitempty" bson:"end_behavior"`
	RoundedPlayer     bool   `json:"rounded_player,omitempty" bson:"rounded_player"`
	AllowSeeking      bool   `json:"allow_seeking" bson:"allow_seeking"`

	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// MediaChannel groups media items into a series, playlist, or channel.
type MediaChannel struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	PublicId     string `json:"public_id,omitempty" bson:"public_id"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id"`

	Title          string              `json:"title,omitempty" bson:"title"`
	Description    string              `json:"description,omitempty" bson:"description,omitempty"`
	Layout         string              `json:"layout,omitempty" bson:"layout,omitempty"`
	Theme          string              `json:"theme,omitempty" bson:"theme,omitempty"`
	AccessMode     string              `json:"access_mode,omitempty" bson:"access_mode,omitempty"`
	PlayerPresetID string              `json:"player_preset_id,omitempty" bson:"player_preset_id,omitempty"`
	Items          []*MediaChannelItem `json:"items,omitempty" bson:"items,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// MediaChannelItem is an ordered reference to a media item within a channel.
type MediaChannelItem struct {
	MediaPublicId string `json:"media_public_id,omitempty" bson:"media_public_id"`
	Order         int    `json:"order,omitempty" bson:"order"`
	TitleOverride string `json:"title_override,omitempty" bson:"title_override,omitempty"`
}

// MediaWebhook configures outbound webhooks for video events.
type MediaWebhook struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	PublicId     string `json:"public_id,omitempty" bson:"public_id"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id"`

	Name       string   `json:"name,omitempty" bson:"name"`
	URL        string   `json:"url,omitempty" bson:"url"`
	Secret     string   `json:"secret,omitempty" bson:"secret,omitempty"`
	EventTypes []string `json:"event_types,omitempty" bson:"event_types,omitempty"`
	Enabled    bool     `json:"enabled,omitempty" bson:"enabled"`

	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// ViewerIdentity tracks a known or anonymous viewer across sessions.
type ViewerIdentity struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id,omitempty"`

	PublicId             string `json:"public_id,omitempty" bson:"public_id"`
	ContactID            string `json:"contact_id,omitempty" bson:"contact_id,omitempty"`
	Email                string `json:"email,omitempty" bson:"email,omitempty"`
	ExternalID           string `json:"external_id,omitempty" bson:"external_id,omitempty"`
	SessionKey           string `json:"session_key,omitempty" bson:"session_key,omitempty"`
	IdentificationSource string `json:"identification_source,omitempty" bson:"identification_source,omitempty"`

	FirstSeenAt *time.Time `json:"first_seen_at,omitempty" bson:"first_seen_at,omitempty"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty" bson:"last_seen_at,omitempty"`
}

// ViewingSession records a single watch session for a media item.
type ViewingSession struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id,omitempty"`

	PublicId       string `json:"public_id,omitempty" bson:"public_id"`
	MediaPublicId  string `json:"media_public_id,omitempty" bson:"media_public_id"`
	ViewerPublicId string `json:"viewer_public_id,omitempty" bson:"viewer_public_id"`

	PageURL    string `json:"page_url,omitempty" bson:"page_url,omitempty"`
	Domain     string `json:"domain,omitempty" bson:"domain,omitempty"`
	Referrer   string `json:"referrer,omitempty" bson:"referrer,omitempty"`
	DeviceType string `json:"device_type,omitempty" bson:"device_type,omitempty"`
	Country    string `json:"country,omitempty" bson:"country,omitempty"`

	DurationWatchedSec int  `json:"duration_watched_sec,omitempty" bson:"duration_watched_sec,omitempty"`
	MaxProgressPct     int  `json:"max_progress_pct,omitempty" bson:"max_progress_pct,omitempty"`
	Completed          bool `json:"completed,omitempty" bson:"completed,omitempty"`

	StartedAt *time.Time `json:"started_at,omitempty" bson:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
}

// MediaEvent records a single player event (append-only).
type MediaEvent struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id,omitempty"`

	MediaPublicId   string `json:"media_public_id,omitempty" bson:"media_public_id"`
	ViewerPublicId  string `json:"viewer_public_id,omitempty" bson:"viewer_public_id,omitempty"`
	SessionPublicId string `json:"session_public_id,omitempty" bson:"session_public_id,omitempty"`

	EventName   string `json:"event_name,omitempty" bson:"event_name"`
	CurrentSec  int    `json:"current_second,omitempty" bson:"current_second,omitempty"`
	ProgressPct int    `json:"progress_pct,omitempty" bson:"progress_pct,omitempty"`

	PageURL string `json:"page_url,omitempty" bson:"page_url,omitempty"`
	Domain  string `json:"domain,omitempty" bson:"domain,omitempty"`

	Data       map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	OccurredAt time.Time              `json:"occurred_at,omitempty" bson:"occurred_at"`
}

// MediaLeadCapture records a turnstile form submission during playback.
type MediaLeadCapture struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	SubscriberId string `json:"subscriber_id,omitempty" bson:"subscriber_id,omitempty"`

	MediaPublicId   string `json:"media_public_id,omitempty" bson:"media_public_id"`
	ViewerPublicId  string `json:"viewer_public_id,omitempty" bson:"viewer_public_id,omitempty"`
	SessionPublicId string `json:"session_public_id,omitempty" bson:"session_public_id,omitempty"`

	FormPublicId string    `json:"form_public_id,omitempty" bson:"form_public_id,omitempty"`
	Email        string    `json:"email,omitempty" bson:"email,omitempty"`
	FirstName    string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	SubmittedAt  time.Time `json:"submitted_at,omitempty" bson:"submitted_at"`
}

// MediaDailyAggregate stores pre-computed daily analytics for a media item.
type MediaDailyAggregate struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID      string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	MediaPublicId string `json:"media_public_id,omitempty" bson:"media_public_id"`
	Date          string `json:"date,omitempty" bson:"date"`

	Plays            int     `json:"plays,omitempty" bson:"plays"`
	UniqueViewers    int     `json:"unique_viewers,omitempty" bson:"unique_viewers"`
	Completions      int     `json:"completions,omitempty" bson:"completions"`
	AvgProgressPct   float64 `json:"avg_progress_pct,omitempty" bson:"avg_progress_pct"`
	AvgWatchTimeSec  float64 `json:"avg_watch_time_sec,omitempty" bson:"avg_watch_time_sec"`
	TurnstileSubmits int     `json:"turnstile_submits,omitempty" bson:"turnstile_submits"`
	CTAClicks        int     `json:"cta_clicks,omitempty" bson:"cta_clicks"`
}
