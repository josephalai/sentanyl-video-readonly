package transcoder

import (
	"context"
	"fmt"
	"log"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	transcoderpb "cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// TranscoderProvider abstracts video transcoding operations.
type TranscoderProvider interface {
	// CreateJob submits a transcoding job.
	CreateJob(input TranscodeInput) (*TranscodeJob, error)

	// GetJobStatus retrieves the current status of a transcoding job.
	GetJobStatus(jobName string) (*TranscodeJob, error)

	// CancelJob cancels a running transcoding job.
	CancelJob(jobName string) error
}

// TranscodeInput describes the input for a transcoding job.
type TranscodeInput struct {
	InputURI     string // gs://bucket/path/to/source.mp4
	OutputPrefix string // gs://bucket/outputs/media-id/
	MediaID      string
	TenantID     string
	// HLS configuration
	EnableHLS   bool
	EnableDASH  bool
	EnableMP4   bool
	Resolutions []string // ["1080p", "720p", "480p"]
}

// TranscodeJob represents the state of a transcoding job.
type TranscodeJob struct {
	JobName      string
	State        string // "PENDING", "RUNNING", "SUCCEEDED", "FAILED"
	InputURI     string
	OutputPrefix string
	DurationSec  int64
	Error        string
	// Output paths
	HLSURL    string
	DASHURL   string
	MP4URL    string
	PosterURL string
}

// GCPTranscoderProvider implements TranscoderProvider using Google Cloud Transcoder API.
type GCPTranscoderProvider struct {
	ProjectID string
	Region    string
	client    *transcoder.Client
}

// NewGCPTranscoderProvider creates a new GCP Transcoder provider.
func NewGCPTranscoderProvider(projectID, region string) (*GCPTranscoderProvider, error) {
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Transcoder client: %w", err)
	}
	return &GCPTranscoderProvider{
		ProjectID: projectID,
		Region:    region,
		client:    client,
	}, nil
}

// Close releases the Transcoder client resources.
func (t *GCPTranscoderProvider) Close() error {
	if t.client != nil {
		return t.client.Close()
	}
	return nil
}

func (t *GCPTranscoderProvider) CreateJob(input TranscodeInput) (*TranscodeJob, error) {
	ctx := context.Background()

	// Build elementary streams for video renditions
	elementaryStreams := []*transcoderpb.ElementaryStream{
		{
			Key: "video-stream-720p",
			ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
				VideoStream: &transcoderpb.VideoStream{
					CodecSettings: &transcoderpb.VideoStream_H264{
						H264: &transcoderpb.VideoStream_H264CodecSettings{
							HeightPixels: 720,
							WidthPixels:  1280,
							BitrateBps:   2500000,
							FrameRate:    30,
						},
					},
				},
			},
		},
		{
			Key: "video-stream-480p",
			ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
				VideoStream: &transcoderpb.VideoStream{
					CodecSettings: &transcoderpb.VideoStream_H264{
						H264: &transcoderpb.VideoStream_H264CodecSettings{
							HeightPixels: 480,
							WidthPixels:  854,
							BitrateBps:   1000000,
							FrameRate:    30,
						},
					},
				},
			},
		},
		{
			Key: "audio-stream",
			ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
				AudioStream: &transcoderpb.AudioStream{
					Codec:      "aac",
					BitrateBps: 128000,
				},
			},
		},
	}

	// Build mux streams
	muxStreams := []*transcoderpb.MuxStream{}

	if input.EnableHLS {
		muxStreams = append(muxStreams,
			&transcoderpb.MuxStream{
				Key:               "hls-720p",
				Container:         "ts",
				ElementaryStreams: []string{"video-stream-720p", "audio-stream"},
				SegmentSettings: &transcoderpb.SegmentSettings{
					SegmentDuration: nil,
				},
			},
			&transcoderpb.MuxStream{
				Key:               "hls-480p",
				Container:         "ts",
				ElementaryStreams: []string{"video-stream-480p", "audio-stream"},
			},
		)
	}

	if input.EnableMP4 {
		muxStreams = append(muxStreams,
			&transcoderpb.MuxStream{
				Key:               "mp4-720p",
				Container:         "mp4",
				ElementaryStreams: []string{"video-stream-720p", "audio-stream"},
			},
		)
	}

	// Build manifests
	manifests := []*transcoderpb.Manifest{}
	if input.EnableHLS {
		manifests = append(manifests, &transcoderpb.Manifest{
			FileName: "manifest.m3u8",
			Type:     transcoderpb.Manifest_HLS,
			MuxStreams: []string{"hls-720p", "hls-480p"},
		})
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", t.ProjectID, t.Region)

	req := &transcoderpb.CreateJobRequest{
		Parent: parent,
		Job: &transcoderpb.Job{
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					Inputs: []*transcoderpb.Input{
						{
							Key: "input0",
							Uri: input.InputURI,
						},
					},
					ElementaryStreams: elementaryStreams,
					MuxStreams:        muxStreams,
					Manifests:         manifests,
					Output: &transcoderpb.Output{
						Uri: input.OutputPrefix,
					},
				},
			},
		},
	}

	job, err := t.client.CreateJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcoder job: %w", err)
	}

	log.Printf("[GCP Transcoder] Created job %s (input=%s, output=%s)", job.Name, input.InputURI, input.OutputPrefix)

	return &TranscodeJob{
		JobName:      job.Name,
		State:        mapJobState(job.State),
		InputURI:     input.InputURI,
		OutputPrefix: input.OutputPrefix,
	}, nil
}

func (t *GCPTranscoderProvider) GetJobStatus(jobName string) (*TranscodeJob, error) {
	ctx := context.Background()

	req := &transcoderpb.GetJobRequest{
		Name: jobName,
	}

	job, err := t.client.GetJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcoder job status: %w", err)
	}

	result := &TranscodeJob{
		JobName: job.Name,
		State:   mapJobState(job.State),
	}

	if job.Error != nil {
		result.Error = job.Error.Message
	}

	return result, nil
}

func (t *GCPTranscoderProvider) CancelJob(jobName string) error {
	ctx := context.Background()

	req := &transcoderpb.DeleteJobRequest{
		Name: jobName,
	}

	if err := t.client.DeleteJob(ctx, req); err != nil {
		return fmt.Errorf("failed to cancel transcoder job: %w", err)
	}

	log.Printf("[GCP Transcoder] Cancelled job %s", jobName)
	return nil
}

func mapJobState(state transcoderpb.Job_ProcessingState) string {
	switch state {
	case transcoderpb.Job_PENDING:
		return "PENDING"
	case transcoderpb.Job_RUNNING:
		return "RUNNING"
	case transcoderpb.Job_SUCCEEDED:
		return "SUCCEEDED"
	case transcoderpb.Job_FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}
