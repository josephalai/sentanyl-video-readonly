# video-service

Video hosting, transcoding, and intelligence service. Manages the full media lifecycle — upload, GCP transcoding, playback URL delivery, viewer analytics, badge rules, and click tracking.

**Ports:** `8084` (intelligence) · `8090` (processing)

## Responsibilities

- Creator-facing media management (CRUD, metadata, player presets, channels)
- Video upload and GCP Cloud Transcoder integration
- GCS storage management for source and transcoded assets
- Viewer session tracking and watch analytics
- Badge rules — automatically grant badges at configurable watch thresholds
- Event ingestion from visitor browsers (no auth required)
- Click tracking redirects

## Architecture

The service compiles to **two separate binaries**:

| Binary | Entry Point | Port | Purpose |
|--------|------------|------|---------|
| Intelligence | `cmd/main.go` | `8084` | Creator-facing REST API |
| Processing | `server.go` | `8090` | Transcoding pipeline and GCP callbacks |

## Directory Structure

```
video-service/
├── cmd/
│   └── main.go          # Intelligence service entry point
├── server.go            # Processing service entry point
├── config.go            # Environment-driven configuration
├── handlers/
│   ├── intelligence.go  # Creator-facing handlers
│   ├── handlers.go      # Processing handlers (upload, transcode, playback)
│   └── tracking.go      # Click tracking logic
├── queries/             # MongoDB queries for media and analytics
└── providers/
    ├── storage/         # GCS provider
    ├── transcoder/      # GCP Video Transcoder provider
    └── events/          # Pub/Sub or HTTP callback event publishing
```

## API Endpoints

### Intelligence Service (port 8084)

**Tenant (JWT required):**

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/video/media` | List media |
| `POST` | `/api/video/media` | Create media with metadata |
| `GET` | `/api/video/media/:id` | Get media details |
| `PUT` | `/api/video/media/:id` | Update media |
| `DELETE` | `/api/video/media/:id` | Delete media |
| `GET/POST` | `/api/video/presets` | Player preset management |
| `GET/POST` | `/api/video/channels` | Channel management |
| `GET/POST` | `/api/video/webhooks` | Webhook management |
| `GET` | `/api/video/analytics/overview` | Analytics dashboard |
| `GET` | `/api/video/media/:id/viewers` | List viewers for a video |
| `GET` | `/api/video/media/:id/sessions` | List viewing sessions |

**Public (no auth):**

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/video/events` | Visitor event ingestion — returns `viewer_id`, `session_id`, badges earned |
| `GET` | `/api/track/click/:token` | Click tracking redirect |

### Processing Service (port 8090)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/media` | Create a media record |
| `POST` | `/media/:id/upload` | Handle file upload |
| `POST` | `/media/:id/process` | Trigger GCP transcoding |
| `GET` | `/media/:id/status` | Check transcoding status |
| `GET` | `/media/:id/playback` | Get HLS/DASH/MP4 playback URLs |
| `POST` | `/event` | Ingest a playback event |
| `POST` | `/identify` | Identify a viewer |
| `POST` | `/callbacks/transcoder` | GCP Transcoder completion callback |

## Data Models

All models live in `pkg/models/video.go`.

**`Media`** — The primary asset record.
- `title`, `description`, `kind`, `status`
- `source_asset_id`, `source_url` — Original file reference
- `poster_asset_id`, `poster_url` — Thumbnail
- `duration_sec`, `file_size_bytes`, `mime_type`
- `processing_provider`, `processing_status`, GCS fields — Transcoding state
- `playback_hls_url`, `playback_dash_url`, `playback_mp4_url` — Delivery URLs
- `renditions[]`, `transcript`, `captions[]`, `chapters[]`, `interactions[]`
- `badge_rules[]` — Automatic badge grants at configurable watch percentages
- `player_controls` — UI customization options

**`MediaRendition`** — An encoded variant (HLS, DASH, MP4) with resolution/bitrate.

**`MediaBadgeRule`** — Grants a badge when a viewer reaches a watch percentage threshold.

**`MediaInteraction`** — An interactive overlay element (clickable, form, annotation) at a timestamp.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8084` / `8090` | HTTP listen port (per binary) |
| `MONGO_HOST` | — | MongoDB host |
| `MONGO_DB` | — | Database name |
| `GCP_PROJECT_ID` | — | GCP project |
| `GCS_BUCKET` | — | Source asset bucket |
| `GCS_OUTPUT_BUCKET` | — | Transcoded output bucket |
| `TRANSCODER_REGION` | — | GCP Video Transcoder region |
| `PUBSUB_TOPIC` | — | Pub/Sub topic for transcoding events |
| `SENTANYL_BASE_URL` | — | Base URL for transcoder callbacks |
| `CDN_BASE_URL` | — | Base URL for playback delivery |

## Dependencies

- [`gin-gonic/gin`](https://github.com/gin-gonic/gin) — HTTP framework
- `gopkg.in/mgo.v2` — MongoDB driver
- `cloud.google.com/go/storage` — GCS operations
- `cloud.google.com/go/video` — GCP Cloud Video Transcoder
- `cloud.google.com/go/pubsub` — Event publishing
- `../pkg` — Shared auth, config, db, models, HTTP utilities
