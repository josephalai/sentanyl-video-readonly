package queries

import (
"log"
"time"

"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"

"github.com/josephalai/sentanyl/pkg/db"
"github.com/josephalai/sentanyl/video-service/models"
)

// VideoQueries provides the query layer for video intelligence entities.
type VideoQueries struct{}

// NewVideoQueries creates a new VideoQueries instance.
func NewVideoQueries() *VideoQueries {
return &VideoQueries{}
}

// ---------- Media CRUD ----------

func (q *VideoQueries) CreateMedia(media *models.Media) (*models.Media, error) {
media.PublicId = bson.NewObjectId().Hex()
now := time.Now()
media.CreatedAt = &now
if media.Status == "" {
media.Status = "draft"
}
if media.Kind == "" {
media.Kind = "video"
}
err := db.GetCollection(models.MediaCollectionName).Insert(media)
if err != nil {
log.Println("CreateMedia error:", err)
return nil, err
}
return media, nil
}

func (q *VideoQueries) GetMediaByPublicId(tenantID, publicId string) (*models.Media, error) {
result := models.Media{}
err := db.GetCollection(models.MediaCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
log.Println("GetMediaByPublicId error:", err)
return nil, err
}
return &result, nil
}

func (q *VideoQueries) ListMedia(tenantID, status string, skip, limit int) ([]*models.Media, error) {
result := []*models.Media{}
query := bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}
if status != "" {
query["status"] = status
}
mgoQ := db.GetCollection(models.MediaCollectionName).Find(query).Sort("-created_at")
if skip > 0 {
mgoQ = mgoQ.Skip(skip)
}
if limit > 0 {
mgoQ = mgoQ.Limit(limit)
}
if err := mgoQ.All(&result); err != nil {
log.Println("ListMedia error:", err)
return nil, err
}
return result, nil
}

func (q *VideoQueries) UpdateMedia(tenantID, publicId string, update map[string]interface{}) (*models.Media, error) {
update["timestamps.updated_at"] = time.Now()
err := db.GetCollection(models.MediaCollectionName).Update(
bson.M{"tenant_id": tenantID, "public_id": publicId, "deleted_at": nil},
bson.M{"$set": update},
)
if err != nil {
log.Println("UpdateMedia error:", err)
return nil, err
}
return q.GetMediaByPublicId(tenantID, publicId)
}

func (q *VideoQueries) DeleteMedia(tenantID, publicId string) (*models.Media, error) {
media, err := q.GetMediaByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
err = db.GetCollection(models.MediaCollectionName).Update(
bson.M{"_id": media.ID},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
if err != nil {
log.Println("DeleteMedia error:", err)
return nil, err
}
media.DeletedAt = &now
return media, nil
}

// ---------- PlayerPreset CRUD ----------

func (q *VideoQueries) CreatePlayerPreset(preset *models.PlayerPreset) (*models.PlayerPreset, error) {
preset.ID = bson.NewObjectId().Hex()
preset.PublicId = bson.NewObjectId().Hex()
now := time.Now()
preset.CreatedAt = &now
if err := db.GetCollection(models.PlayerPresetCollectionName).Insert(preset); err != nil {
log.Println("CreatePlayerPreset error:", err)
return nil, err
}
return preset, nil
}

func (q *VideoQueries) GetPlayerPresetByPublicId(tenantID, publicId string) (*models.PlayerPreset, error) {
result := models.PlayerPreset{}
err := db.GetCollection(models.PlayerPresetCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) ListPlayerPresets(tenantID string) ([]*models.PlayerPreset, error) {
result := []*models.PlayerPreset{}
err := db.GetCollection(models.PlayerPresetCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) DeletePlayerPreset(tenantID, publicId string) (*models.PlayerPreset, error) {
preset, err := q.GetPlayerPresetByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(models.PlayerPresetCollectionName).Update(
bson.M{"_id": preset.ID},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
preset.DeletedAt = &now
return preset, nil
}

// ---------- MediaChannel CRUD ----------

func (q *VideoQueries) CreateMediaChannel(channel *models.MediaChannel) (*models.MediaChannel, error) {
channel.ID = bson.NewObjectId().Hex()
channel.PublicId = bson.NewObjectId().Hex()
now := time.Now()
channel.CreatedAt = &now
if err := db.GetCollection(models.MediaChannelCollectionName).Insert(channel); err != nil {
log.Println("CreateMediaChannel error:", err)
return nil, err
}
return channel, nil
}

func (q *VideoQueries) ListMediaChannels(tenantID string) ([]*models.MediaChannel, error) {
result := []*models.MediaChannel{}
err := db.GetCollection(models.MediaChannelCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) GetMediaChannelByPublicId(tenantID, publicId string) (*models.MediaChannel, error) {
result := models.MediaChannel{}
err := db.GetCollection(models.MediaChannelCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) DeleteMediaChannel(tenantID, publicId string) (*models.MediaChannel, error) {
ch, err := q.GetMediaChannelByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(models.MediaChannelCollectionName).Update(
bson.M{"_id": ch.ID},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
ch.DeletedAt = &now
return ch, nil
}

// ---------- MediaWebhook CRUD ----------

func (q *VideoQueries) CreateMediaWebhook(webhook *models.MediaWebhook) (*models.MediaWebhook, error) {
webhook.ID = bson.NewObjectId().Hex()
webhook.PublicId = bson.NewObjectId().Hex()
now := time.Now()
webhook.CreatedAt = &now
if err := db.GetCollection(models.MediaWebhookCollectionName).Insert(webhook); err != nil {
log.Println("CreateMediaWebhook error:", err)
return nil, err
}
return webhook, nil
}

func (q *VideoQueries) ListMediaWebhooks(tenantID string) ([]*models.MediaWebhook, error) {
result := []*models.MediaWebhook{}
err := db.GetCollection(models.MediaWebhookCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) GetMediaWebhookByPublicId(tenantID, publicId string) (*models.MediaWebhook, error) {
result := models.MediaWebhook{}
err := db.GetCollection(models.MediaWebhookCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) DeleteMediaWebhook(tenantID, publicId string) (*models.MediaWebhook, error) {
wh, err := q.GetMediaWebhookByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(models.MediaWebhookCollectionName).Update(
bson.M{"_id": wh.ID},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
wh.DeletedAt = &now
return wh, nil
}

// ---------- ViewerIdentity Queries ----------

func (q *VideoQueries) FindOrCreateViewer(tenantID, email, sessionKey, source string) (*models.ViewerIdentity, error) {
result := models.ViewerIdentity{}
if email != "" {
err := db.GetCollection(models.ViewerIdentityCollectionName).Find(bson.M{
"tenant_id": tenantID,
"email":     email,
}).One(&result)
if err == nil {
now := time.Now()
db.GetCollection(models.ViewerIdentityCollectionName).Update(
bson.M{"_id": result.ID},
bson.M{"$set": bson.M{"last_seen_at": now}},
)
result.LastSeenAt = &now
return &result, nil
}
if err != mgo.ErrNotFound {
return nil, err
}
}
now := time.Now()
viewer := &models.ViewerIdentity{
ID:                   bson.NewObjectId().Hex(),
TenantID:             tenantID,
PublicId:             bson.NewObjectId().Hex(),
Email:                email,
SessionKey:           sessionKey,
IdentificationSource: source,
FirstSeenAt:          &now,
LastSeenAt:           &now,
}
if err := db.GetCollection(models.ViewerIdentityCollectionName).Insert(viewer); err != nil {
return nil, err
}
return viewer, nil
}

func (q *VideoQueries) ListViewersByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*models.ViewerIdentity, error) {
var viewerIds []string
db.GetCollection(models.ViewingSessionCollectionName).Find(bson.M{
"tenant_id":       tenantID,
"media_public_id": mediaPublicId,
}).Distinct("viewer_public_id", &viewerIds)

if len(viewerIds) == 0 {
return []*models.ViewerIdentity{}, nil
}
result := []*models.ViewerIdentity{}
mgoQ := db.GetCollection(models.ViewerIdentityCollectionName).Find(bson.M{
"tenant_id": tenantID,
"public_id": bson.M{"$in": viewerIds},
}).Sort("-last_seen_at")
if skip > 0 {
mgoQ = mgoQ.Skip(skip)
}
if limit > 0 {
mgoQ = mgoQ.Limit(limit)
}
mgoQ.All(&result)
return result, nil
}

// ---------- ViewingSession Queries ----------

func (q *VideoQueries) CreateViewingSession(session *models.ViewingSession) (*models.ViewingSession, error) {
session.ID = bson.NewObjectId().Hex()
if session.PublicId == "" {
session.PublicId = bson.NewObjectId().Hex()
}
now := time.Now()
session.StartedAt = &now
if err := db.GetCollection(models.ViewingSessionCollectionName).Insert(session); err != nil {
log.Println("CreateViewingSession error:", err)
return nil, err
}
return session, nil
}

func (q *VideoQueries) UpdateViewingSession(sessionID string, update map[string]interface{}) error {
return db.GetCollection(models.ViewingSessionCollectionName).Update(
bson.M{"_id": sessionID},
bson.M{"$set": update},
)
}

func (q *VideoQueries) ListSessionsByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*models.ViewingSession, error) {
result := []*models.ViewingSession{}
mgoQ := db.GetCollection(models.ViewingSessionCollectionName).Find(bson.M{
"tenant_id":       tenantID,
"media_public_id": mediaPublicId,
}).Sort("-started_at")
if skip > 0 {
mgoQ = mgoQ.Skip(skip)
}
if limit > 0 {
mgoQ = mgoQ.Limit(limit)
}
mgoQ.All(&result)
return result, nil
}

// ---------- MediaEvent Queries ----------

func (q *VideoQueries) AppendMediaEvent(event *models.MediaEvent) (*models.MediaEvent, error) {
event.ID = bson.NewObjectId().Hex()
if event.OccurredAt.IsZero() {
event.OccurredAt = time.Now()
}
if err := db.GetCollection(models.MediaEventCollectionName).Insert(event); err != nil {
log.Println("AppendMediaEvent error:", err)
return nil, err
}
return event, nil
}

// ---------- MediaLeadCapture Queries ----------

func (q *VideoQueries) CreateMediaLeadCapture(capture *models.MediaLeadCapture) (*models.MediaLeadCapture, error) {
capture.ID = bson.NewObjectId().Hex()
if capture.SubmittedAt.IsZero() {
capture.SubmittedAt = time.Now()
}
if err := db.GetCollection(models.MediaLeadCaptureCollectionName).Insert(capture); err != nil {
log.Println("CreateMediaLeadCapture error:", err)
return nil, err
}
return capture, nil
}

// ---------- Analytics ----------

func (q *VideoQueries) UpsertMediaDailyAggregate(tenantID, mediaPublicId, date string, inc map[string]interface{}) error {
query := bson.M{
"tenant_id":       tenantID,
"media_public_id": mediaPublicId,
"date":            date,
}
update := bson.M{
"$inc": inc,
"$setOnInsert": bson.M{
"_id":             bson.NewObjectId().Hex(),
"tenant_id":       tenantID,
"media_public_id": mediaPublicId,
"date":            date,
},
}
_, err := db.GetCollection(models.MediaDailyAggregateCollectionName).Upsert(query, update)
return err
}

func (q *VideoQueries) GetMediaAnalyticsOverview(tenantID string) (map[string]interface{}, error) {
result := map[string]interface{}{}

mediaCount, err := db.GetCollection(models.MediaCollectionName).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Count()
if err != nil {
return nil, err
}
result["total_media"] = mediaCount

playCount, err := db.GetCollection(models.MediaEventCollectionName).Find(bson.M{
"tenant_id":  tenantID,
"event_name": models.VideoEventPlay,
}).Count()
if err != nil {
return nil, err
}
result["total_plays"] = playCount

viewerCount, err := db.GetCollection(models.ViewerIdentityCollectionName).Find(bson.M{
"tenant_id": tenantID,
}).Count()
if err != nil {
return nil, err
}
result["total_viewers"] = viewerCount

captureCount, err := db.GetCollection(models.MediaLeadCaptureCollectionName).Find(bson.M{
"tenant_id": tenantID,
}).Count()
if err != nil {
return nil, err
}
result["total_lead_captures"] = captureCount

return result, nil
}

// ---------- Badge Rule Evaluation ----------

func (q *VideoQueries) EvaluateMediaBadgeRules(media *models.Media, event *models.MediaEvent) []string {
if media == nil || media.BadgeRules == nil {
return nil
}
var badgesToGrant []string
for _, rule := range media.BadgeRules {
if !rule.Enabled {
continue
}
if matchesEventRule(rule, event) {
badgesToGrant = append(badgesToGrant, rule.BadgePublicId)
}
}
return badgesToGrant
}

func matchesEventRule(rule *models.MediaBadgeRule, event *models.MediaEvent) bool {
if rule.EventName != event.EventName {
if rule.EventName == "progress" && event.EventName == models.VideoEventProgress {
return evaluateThreshold(rule.Operator, event.ProgressPct, rule.Threshold)
}
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
default:
return actual >= threshold
}
}
