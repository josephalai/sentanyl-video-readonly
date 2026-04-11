package queries

import (
"log"
"time"

"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"

"github.com/josephalai/sentanyl/pkg/db"
pkgmodels "github.com/josephalai/sentanyl/pkg/models"
)

// VideoQueries provides the query layer for video intelligence entities.
type VideoQueries struct{}

// NewVideoQueries creates a new VideoQueries instance.
func NewVideoQueries() *VideoQueries {
return &VideoQueries{}
}

// ---------- Media CRUD ----------

func (q *VideoQueries) CreateMedia(media *pkgmodels.Media) (*pkgmodels.Media, error) {
media.PublicId = bson.NewObjectId().Hex()
now := time.Now()
media.CreatedAt = &now
if media.Status == "" {
media.Status = "draft"
}
if media.Kind == "" {
media.Kind = "video"
}
err := db.GetCollection(pkgmodels.MediaCollection).Insert(media)
if err != nil {
log.Println("CreateMedia error:", err)
return nil, err
}
return media, nil
}

func (q *VideoQueries) GetMediaByPublicId(tenantID, publicId string) (*pkgmodels.Media, error) {
result := pkgmodels.Media{}
err := db.GetCollection(pkgmodels.MediaCollection).Find(bson.M{
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

func (q *VideoQueries) ListMedia(tenantID, status string, skip, limit int) ([]*pkgmodels.Media, error) {
result := []*pkgmodels.Media{}
query := bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}
if status != "" {
query["status"] = status
}
mgoQ := db.GetCollection(pkgmodels.MediaCollection).Find(query).Sort("-created_at")
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

func (q *VideoQueries) UpdateMedia(tenantID, publicId string, update map[string]interface{}) (*pkgmodels.Media, error) {
update["timestamps.updated_at"] = time.Now()
err := db.GetCollection(pkgmodels.MediaCollection).Update(
bson.M{"tenant_id": tenantID, "public_id": publicId, "deleted_at": nil},
bson.M{"$set": update},
)
if err != nil {
log.Println("UpdateMedia error:", err)
return nil, err
}
return q.GetMediaByPublicId(tenantID, publicId)
}

func (q *VideoQueries) DeleteMedia(tenantID, publicId string) (*pkgmodels.Media, error) {
media, err := q.GetMediaByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
err = db.GetCollection(pkgmodels.MediaCollection).Update(
bson.M{"_id": media.Id},
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

func (q *VideoQueries) CreatePlayerPreset(preset *pkgmodels.PlayerPreset) (*pkgmodels.PlayerPreset, error) {
preset.Id = bson.NewObjectId()
preset.PublicId = bson.NewObjectId().Hex()
now := time.Now()
preset.CreatedAt = &now
if err := db.GetCollection(pkgmodels.PlayerPresetCollection).Insert(preset); err != nil {
log.Println("CreatePlayerPreset error:", err)
return nil, err
}
return preset, nil
}

func (q *VideoQueries) GetPlayerPresetByPublicId(tenantID, publicId string) (*pkgmodels.PlayerPreset, error) {
result := pkgmodels.PlayerPreset{}
err := db.GetCollection(pkgmodels.PlayerPresetCollection).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) ListPlayerPresets(tenantID string) ([]*pkgmodels.PlayerPreset, error) {
result := []*pkgmodels.PlayerPreset{}
err := db.GetCollection(pkgmodels.PlayerPresetCollection).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) DeletePlayerPreset(tenantID, publicId string) (*pkgmodels.PlayerPreset, error) {
preset, err := q.GetPlayerPresetByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(pkgmodels.PlayerPresetCollection).Update(
bson.M{"_id": preset.Id},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
preset.DeletedAt = &now
return preset, nil
}

// ---------- MediaChannel CRUD ----------

func (q *VideoQueries) CreateMediaChannel(channel *pkgmodels.MediaChannel) (*pkgmodels.MediaChannel, error) {
channel.Id = bson.NewObjectId()
channel.PublicId = bson.NewObjectId().Hex()
now := time.Now()
channel.CreatedAt = &now
if err := db.GetCollection(pkgmodels.MediaChannelCollection).Insert(channel); err != nil {
log.Println("CreateMediaChannel error:", err)
return nil, err
}
return channel, nil
}

func (q *VideoQueries) ListMediaChannels(tenantID string) ([]*pkgmodels.MediaChannel, error) {
result := []*pkgmodels.MediaChannel{}
err := db.GetCollection(pkgmodels.MediaChannelCollection).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) GetMediaChannelByPublicId(tenantID, publicId string) (*pkgmodels.MediaChannel, error) {
result := pkgmodels.MediaChannel{}
err := db.GetCollection(pkgmodels.MediaChannelCollection).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) DeleteMediaChannel(tenantID, publicId string) (*pkgmodels.MediaChannel, error) {
ch, err := q.GetMediaChannelByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(pkgmodels.MediaChannelCollection).Update(
bson.M{"_id": ch.Id},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
ch.DeletedAt = &now
return ch, nil
}

// ---------- MediaWebhook CRUD ----------

func (q *VideoQueries) CreateMediaWebhook(webhook *pkgmodels.MediaWebhook) (*pkgmodels.MediaWebhook, error) {
webhook.Id = bson.NewObjectId()
webhook.PublicId = bson.NewObjectId().Hex()
now := time.Now()
webhook.CreatedAt = &now
if err := db.GetCollection(pkgmodels.MediaWebhookCollection).Insert(webhook); err != nil {
log.Println("CreateMediaWebhook error:", err)
return nil, err
}
return webhook, nil
}

func (q *VideoQueries) ListMediaWebhooks(tenantID string) ([]*pkgmodels.MediaWebhook, error) {
result := []*pkgmodels.MediaWebhook{}
err := db.GetCollection(pkgmodels.MediaWebhookCollection).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Sort("-created_at").All(&result)
if err != nil {
return nil, err
}
return result, nil
}

func (q *VideoQueries) GetMediaWebhookByPublicId(tenantID, publicId string) (*pkgmodels.MediaWebhook, error) {
result := pkgmodels.MediaWebhook{}
err := db.GetCollection(pkgmodels.MediaWebhookCollection).Find(bson.M{
"tenant_id":             tenantID,
"public_id":             publicId,
"deleted_at": nil,
}).One(&result)
if err != nil {
return nil, err
}
return &result, nil
}

func (q *VideoQueries) DeleteMediaWebhook(tenantID, publicId string) (*pkgmodels.MediaWebhook, error) {
wh, err := q.GetMediaWebhookByPublicId(tenantID, publicId)
if err != nil {
return nil, err
}
now := time.Now()
db.GetCollection(pkgmodels.MediaWebhookCollection).Update(
bson.M{"_id": wh.Id},
bson.M{"$set": bson.M{"timestamps.deleted_at": now}},
)
wh.DeletedAt = &now
return wh, nil
}

// ---------- ViewerIdentity Queries ----------

func (q *VideoQueries) FindOrCreateViewer(tenantID, email, sessionKey, source string) (*pkgmodels.ViewerIdentity, error) {
result := pkgmodels.ViewerIdentity{}
if email != "" {
err := db.GetCollection(pkgmodels.ViewerIdentityCollection).Find(bson.M{
"tenant_id": tenantID,
"email":     email,
}).One(&result)
if err == nil {
now := time.Now()
db.GetCollection(pkgmodels.ViewerIdentityCollection).Update(
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
viewer := &pkgmodels.ViewerIdentity{
ID:                   bson.NewObjectId().Hex(),
TenantID:             tenantID,
PublicId:             bson.NewObjectId().Hex(),
Email:                email,
SessionKey:           sessionKey,
IdentificationSource: source,
FirstSeenAt:          &now,
LastSeenAt:           &now,
}
if err := db.GetCollection(pkgmodels.ViewerIdentityCollection).Insert(viewer); err != nil {
return nil, err
}
return viewer, nil
}

func (q *VideoQueries) ListViewersByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*pkgmodels.ViewerIdentity, error) {
var viewerIds []string
db.GetCollection(pkgmodels.ViewingSessionCollection).Find(bson.M{
"tenant_id":       tenantID,
"media_public_id": mediaPublicId,
}).Distinct("viewer_public_id", &viewerIds)

if len(viewerIds) == 0 {
return []*pkgmodels.ViewerIdentity{}, nil
}
result := []*pkgmodels.ViewerIdentity{}
mgoQ := db.GetCollection(pkgmodels.ViewerIdentityCollection).Find(bson.M{
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

func (q *VideoQueries) CreateViewingSession(session *pkgmodels.ViewingSession) (*pkgmodels.ViewingSession, error) {
session.ID = bson.NewObjectId().Hex()
if session.PublicId == "" {
session.PublicId = bson.NewObjectId().Hex()
}
now := time.Now()
session.StartedAt = &now
if err := db.GetCollection(pkgmodels.ViewingSessionCollection).Insert(session); err != nil {
log.Println("CreateViewingSession error:", err)
return nil, err
}
return session, nil
}

func (q *VideoQueries) UpdateViewingSession(sessionID string, update map[string]interface{}) error {
return db.GetCollection(pkgmodels.ViewingSessionCollection).Update(
bson.M{"_id": sessionID},
bson.M{"$set": update},
)
}

func (q *VideoQueries) ListSessionsByMedia(tenantID, mediaPublicId string, skip, limit int) ([]*pkgmodels.ViewingSession, error) {
result := []*pkgmodels.ViewingSession{}
mgoQ := db.GetCollection(pkgmodels.ViewingSessionCollection).Find(bson.M{
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

func (q *VideoQueries) AppendMediaEvent(event *pkgmodels.MediaEvent) (*pkgmodels.MediaEvent, error) {
event.ID = bson.NewObjectId().Hex()
if event.OccurredAt.IsZero() {
event.OccurredAt = time.Now()
}
if err := db.GetCollection(pkgmodels.MediaEventCollection).Insert(event); err != nil {
log.Println("AppendMediaEvent error:", err)
return nil, err
}
return event, nil
}

// ---------- MediaLeadCapture Queries ----------

func (q *VideoQueries) CreateMediaLeadCapture(capture *pkgmodels.MediaLeadCapture) (*pkgmodels.MediaLeadCapture, error) {
capture.ID = bson.NewObjectId().Hex()
if capture.SubmittedAt.IsZero() {
capture.SubmittedAt = time.Now()
}
if err := db.GetCollection(pkgmodels.MediaLeadCaptureCollection).Insert(capture); err != nil {
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
_, err := db.GetCollection(pkgmodels.MediaDailyAggregateCollection).Upsert(query, update)
return err
}

func (q *VideoQueries) GetMediaAnalyticsOverview(tenantID string) (map[string]interface{}, error) {
result := map[string]interface{}{}

mediaCount, err := db.GetCollection(pkgmodels.MediaCollection).Find(bson.M{
"tenant_id":             tenantID,
"deleted_at": nil,
}).Count()
if err != nil {
return nil, err
}
result["total_media"] = mediaCount

playCount, err := db.GetCollection(pkgmodels.MediaEventCollection).Find(bson.M{
"tenant_id":  tenantID,
"event_name": pkgmodels.VideoEventPlay,
}).Count()
if err != nil {
return nil, err
}
result["total_plays"] = playCount

viewerCount, err := db.GetCollection(pkgmodels.ViewerIdentityCollection).Find(bson.M{
"tenant_id": tenantID,
}).Count()
if err != nil {
return nil, err
}
result["total_viewers"] = viewerCount

captureCount, err := db.GetCollection(pkgmodels.MediaLeadCaptureCollection).Find(bson.M{
"tenant_id": tenantID,
}).Count()
if err != nil {
return nil, err
}
result["total_lead_captures"] = captureCount

return result, nil
}

// ---------- Badge Rule Evaluation ----------

func (q *VideoQueries) EvaluateMediaBadgeRules(media *pkgmodels.Media, event *pkgmodels.MediaEvent) []string {
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

func matchesEventRule(rule *pkgmodels.MediaBadgeRule, event *pkgmodels.MediaEvent) bool {
if rule.EventName != event.EventName {
if rule.EventName == "progress" && event.EventName == pkgmodels.VideoEventProgress {
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
