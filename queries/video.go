package queries

import (
	"errors"
	"log"
	"time"

	"github.com/josephalai/sentanyl/video-service/models"
)

// VideoQueries provides the query layer for video intelligence entities.
// In the monolith these used mgo directly; here they are stubs that will
// be wired to the shared pkg/db layer.
type VideoQueries struct {
	// TODO: inject *mgo.Database or pkg/db reference
}

// NewVideoQueries creates a new VideoQueries instance.
func NewVideoQueries() *VideoQueries {
	return &VideoQueries{}
}

// ---------- Media ----------

func (q *VideoQueries) ListMedia(tenantID, status string, skip, limit int) ([]*models.Media, error) {
	// TODO: implement with pkg/db
	log.Printf("ListMedia stub: tenant=%s status=%s skip=%d limit=%d", tenantID, status, skip, limit)
	return []*models.Media{}, nil
}

func (q *VideoQueries) GetMediaByPublicId(tenantID, publicId string) (*models.Media, error) {
	// TODO: implement with pkg/db
	log.Printf("GetMediaByPublicId stub: tenant=%s publicId=%s", tenantID, publicId)
	return nil, errors.New("not implemented")
}

func (q *VideoQueries) UpdateMedia(tenantID, publicId string, update map[string]interface{}) (*models.Media, error) {
	// TODO: implement with pkg/db
	log.Printf("UpdateMedia stub: tenant=%s publicId=%s", tenantID, publicId)
	return nil, errors.New("not implemented")
}

func (q *VideoQueries) DeleteMedia(tenantID, publicId string) error {
	// TODO: implement with pkg/db
	log.Printf("DeleteMedia stub: tenant=%s publicId=%s", tenantID, publicId)
	return nil
}

// ---------- PlayerPreset ----------

func (q *VideoQueries) ListPlayerPresets(tenantID string) ([]*models.PlayerPreset, error) {
	log.Printf("ListPlayerPresets stub: tenant=%s", tenantID)
	return []*models.PlayerPreset{}, nil
}

// ---------- MediaChannel ----------

func (q *VideoQueries) ListMediaChannels(tenantID string) ([]*models.MediaChannel, error) {
	log.Printf("ListMediaChannels stub: tenant=%s", tenantID)
	return []*models.MediaChannel{}, nil
}

// ---------- MediaWebhook ----------

func (q *VideoQueries) ListMediaWebhooks(tenantID string) ([]*models.MediaWebhook, error) {
	log.Printf("ListMediaWebhooks stub: tenant=%s", tenantID)
	return []*models.MediaWebhook{}, nil
}

// ---------- Viewer / Session ----------

func (q *VideoQueries) ListViewersByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*models.ViewerIdentity, error) {
	log.Printf("ListViewersByMedia stub: tenant=%s media=%s", tenantID, mediaPublicId)
	return []*models.ViewerIdentity{}, nil
}

func (q *VideoQueries) ListSessionsByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*models.ViewingSession, error) {
	log.Printf("ListSessionsByMedia stub: tenant=%s media=%s", tenantID, mediaPublicId)
	return []*models.ViewingSession{}, nil
}

// ---------- Analytics ----------

func (q *VideoQueries) GetMediaAnalyticsOverview(tenantID string) (map[string]interface{}, error) {
	log.Printf("GetMediaAnalyticsOverview stub: tenant=%s", tenantID)
	return map[string]interface{}{
		"total_media":         0,
		"total_plays":         0,
		"total_viewers":       0,
		"total_lead_captures": 0,
	}, nil
}

// ---------- Badge Rule Evaluation ----------

// EvaluateMediaBadgeRules checks badge rules against an event.
func EvaluateMediaBadgeRules(media *models.Media, event *models.MediaEvent) []string {
	if media == nil || media.BadgeRules == nil {
		return nil
	}
	var badges []string
	for _, rule := range media.BadgeRules {
		if !rule.Enabled {
			continue
		}
		if !matchesEventRule(rule, event) {
			continue
		}
		badges = append(badges, rule.BadgePublicId)
	}
	return badges
}

func matchesEventRule(rule *models.MediaBadgeRule, event *models.MediaEvent) bool {
	if rule.EventName == "progress" && event.EventName == models.VideoEventProgress {
		return evaluateThreshold(rule.Operator, event.ProgressPct, rule.Threshold)
	}
	if rule.EventName != event.EventName {
		return false
	}
	if rule.EventName == "progress" {
		return evaluateThreshold(rule.Operator, event.ProgressPct, rule.Threshold)
	}
	return true
}

func evaluateThreshold(operator string, actual, threshold int) bool {
	switch operator {
	case ">":
		return actual > threshold
	case ">=":
		return actual >= threshold
	case "<":
		return actual < threshold
	case "<=":
		return actual <= threshold
	case "==":
		return actual == threshold
	case "exists":
		return true
	default:
		return actual >= threshold
	}
}

// AppendMediaEvent is a stub for inserting a media event.
func (q *VideoQueries) AppendMediaEvent(event *models.MediaEvent) (*models.MediaEvent, error) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	// TODO: implement with pkg/db
	log.Printf("AppendMediaEvent stub: media=%s event=%s", event.MediaPublicId, event.EventName)
	return event, nil
}
