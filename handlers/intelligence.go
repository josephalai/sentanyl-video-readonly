package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	pkgmodels "github.com/josephalai/sentanyl/pkg/models"
	"github.com/josephalai/sentanyl/video-service/queries"
)

// IntelligenceHandler holds dependencies for creator-facing video intelligence endpoints.
type IntelligenceHandler struct {
	Queries *queries.VideoQueries
}

// ---------- Media CRUD Handlers ----------

func (h *IntelligenceHandler) HandleListMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status := c.Query("status")
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	media, err := h.Queries.ListMedia(tenantID, status, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": media})
}

func (h *IntelligenceHandler) HandleGetMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	publicId := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	media, err := h.Queries.GetMediaByPublicId(tenantID, publicId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *IntelligenceHandler) HandleCreateMediaIntel(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Title       string                       `json:"title" binding:"required"`
		Description string                       `json:"description"`
		Kind        string                       `json:"kind"`
		SourceURL   string                       `json:"source_url"`
		PosterURL   string                       `json:"poster_url"`
		BadgeRules  []*pkgmodels.MediaBadgeRule  `json:"badge_rules"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	media := &pkgmodels.Media{
		Title:        req.Title,
		Description:  req.Description,
		Kind:         req.Kind,
		SourceURL:    req.SourceURL,
		PosterURL:    req.PosterURL,
		BadgeRules:   req.BadgeRules,
		SubscriberId: tenantID,
	}
	if bson.IsObjectIdHex(tenantID) {
		media.TenantID = bson.ObjectIdHex(tenantID)
	}

	created, err := h.Queries.CreateMedia(media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create media"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"media": created})
}

func (h *IntelligenceHandler) HandleUpdateMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	publicId := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	media, err := h.Queries.UpdateMedia(tenantID, publicId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *IntelligenceHandler) HandleDeleteMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	publicId := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	_, err := h.Queries.DeleteMedia(tenantID, publicId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ---------- PlayerPreset Handlers ----------

func (h *IntelligenceHandler) HandleListPlayerPresets(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	presets, err := h.Queries.ListPlayerPresets(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list presets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"presets": presets})
}

func (h *IntelligenceHandler) HandleCreatePlayerPreset(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: Implement full preset creation via queries layer
	c.JSON(http.StatusCreated, gin.H{"status": "created", "tenant_id": tenantID})
}

// ---------- MediaChannel Handlers ----------

func (h *IntelligenceHandler) HandleListChannels(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	channels, err := h.Queries.ListMediaChannels(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"channels": channels})
}

func (h *IntelligenceHandler) HandleCreateChannel(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: Implement full channel creation via queries layer
	c.JSON(http.StatusCreated, gin.H{"status": "created", "tenant_id": tenantID})
}

// ---------- MediaWebhook Handlers ----------

func (h *IntelligenceHandler) HandleListWebhooks(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	webhooks, err := h.Queries.ListMediaWebhooks(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list webhooks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"webhooks": webhooks})
}

func (h *IntelligenceHandler) HandleCreateWebhook(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: Implement full webhook creation via queries layer
	c.JSON(http.StatusCreated, gin.H{"status": "created", "tenant_id": tenantID})
}

// ---------- Analytics Handlers ----------

func (h *IntelligenceHandler) HandleGetAnalyticsOverview(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	overview, err := h.Queries.GetMediaAnalyticsOverview(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analytics"})
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *IntelligenceHandler) HandleListViewersByMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	mediaPublicId := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	viewers, err := h.Queries.ListViewersByMedia(tenantID, mediaPublicId, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list viewers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"viewers": viewers})
}

func (h *IntelligenceHandler) HandleListSessionsByMedia(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	mediaPublicId := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	sessions, err := h.Queries.ListSessionsByMedia(tenantID, mediaPublicId, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// HandleIngestEvent is the player-side event ingestion endpoint.
//
// POST /api/video/events
//
// The player sends this for every play, pause, progress, complete, identify,
// turnstile_submit, and cta_click event.  The handler:
//   1. Finds or creates a ViewerIdentity (by email or anonymous session key)
//   2. Creates a ViewingSession on the first "play" event for a session key,
//      or updates an existing one on progress/complete
//   3. Appends a raw MediaEvent record
//   4. Evaluates the media's BadgeRules and returns any badges to grant so
//      the player can call back to the badge service
//   5. Updates the daily aggregate counters
//
// No auth middleware — the player runs in a visitor's browser.
func (h *IntelligenceHandler) HandleIngestEvent(c *gin.Context) {
	var req pkgmodels.PlayerEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event payload"})
		return
	}
	if req.TenantID == "" || req.MediaID == "" || req.EventName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id, media_id, and event_name are required"})
		return
	}

	// 1. Resolve / create the viewer identity.
	// ViewerID may be an email address (from an identify call) or a random session key.
	// FindOrCreateViewer deduplicates by email when it looks like one.
	viewerEmail := ""
	sessionKey := req.SessionID
	if req.ViewerID != "" {
		// Treat as email if it contains "@", otherwise use as session key.
		if len(req.ViewerID) > 0 {
			for _, ch := range req.ViewerID {
				if ch == '@' {
					viewerEmail = req.ViewerID
					break
				}
			}
			if viewerEmail == "" {
				sessionKey = req.ViewerID
			}
		}
	}
	viewer, err := h.Queries.FindOrCreateViewer(req.TenantID, viewerEmail, sessionKey, "player")
	if err != nil {
		log.Printf("HandleIngestEvent: FindOrCreateViewer error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "viewer resolution failed"})
		return
	}

	// 2. Manage the viewing session.
	sessionPublicId := req.SessionID
	switch req.EventName {
	case pkgmodels.VideoEventPlay:
		// Start a new viewing session on play, using the player's session key
		// as PublicId so subsequent progress/complete events can find it.
		session := &pkgmodels.ViewingSession{
			TenantID:       req.TenantID,
			PublicId:       req.SessionID, // store player-defined key
			MediaPublicId:  req.MediaID,
			ViewerPublicId: viewer.PublicId,
			PageURL:        req.PageURL,
			Domain:         req.Domain,
			Referrer:       req.Referrer,
		}
		created, serr := h.Queries.CreateViewingSession(session)
		if serr == nil {
			sessionPublicId = created.PublicId
		}

	case pkgmodels.VideoEventProgress:
		// Update max progress and duration on progress events.
		if sessionPublicId != "" {
			h.Queries.UpdateViewingSession(sessionPublicId, map[string]interface{}{
				"max_progress_pct":     req.ProgressPercent,
				"duration_watched_sec": req.CurrentSecond,
			})
		}

	case pkgmodels.VideoEventComplete:
		// Mark session completed.
		if sessionPublicId != "" {
			now := time.Now()
			h.Queries.UpdateViewingSession(sessionPublicId, map[string]interface{}{
				"completed":            true,
				"max_progress_pct":     100,
				"duration_watched_sec": req.CurrentSecond,
				"ended_at":             now,
			})
		}
	}

	// 3. Append the raw event record.
	event := &pkgmodels.MediaEvent{
		TenantID:        req.TenantID,
		MediaPublicId:   req.MediaID,
		ViewerPublicId:  viewer.PublicId,
		SessionPublicId: sessionPublicId,
		EventName:       req.EventName,
		CurrentSec:      req.CurrentSecond,
		ProgressPct:     req.ProgressPercent,
		PageURL:         req.PageURL,
		Domain:          req.Domain,
		OccurredAt:      time.Now(),
	}
	h.Queries.AppendMediaEvent(event)

	// 4. Evaluate badge rules on the media.
	media, merr := h.Queries.GetMediaByPublicId(req.TenantID, req.MediaID)
	var badgesToGrant []string
	if merr == nil && media != nil {
		badgesToGrant = h.Queries.EvaluateMediaBadgeRules(media, event)
		if len(badgesToGrant) > 0 {
			log.Printf("HandleIngestEvent: media=%s event=%s pct=%d → badges %v",
				req.MediaID, req.EventName, req.ProgressPercent, badgesToGrant)
		}
	}

	// 5. Update daily aggregate counters.
	date := time.Now().Format("2006-01-02")
	inc := map[string]interface{}{}
	switch req.EventName {
	case pkgmodels.VideoEventPlay:
		inc["plays"] = 1
	case pkgmodels.VideoEventComplete:
		inc["completions"] = 1
	case pkgmodels.VideoEventTurnstileSubmit:
		inc["turnstile_submits"] = 1
	case pkgmodels.VideoEventCTAClick:
		inc["cta_clicks"] = 1
	}
	if len(inc) > 0 {
		h.Queries.UpsertMediaDailyAggregate(req.TenantID, req.MediaID, date, inc)
	}

	c.JSON(http.StatusOK, gin.H{
		"viewer_id":      viewer.PublicId,
		"session_id":     sessionPublicId,
		"badges_to_grant": badgesToGrant,
	})
}
