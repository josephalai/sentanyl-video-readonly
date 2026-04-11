package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description"`
		Kind        string   `json:"kind"`
		Tags        []string `json:"tags"`
		Folder      string   `json:"folder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	// TODO: Implement full media creation via queries layer
	c.JSON(http.StatusCreated, gin.H{
		"status":    "created",
		"title":     req.Title,
		"tenant_id": tenantID,
	})
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

	err := h.Queries.DeleteMedia(tenantID, publicId)
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
