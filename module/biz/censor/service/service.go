package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
)

type Service interface {
	UpdateCensorConfig(ctx context.Context, mod *model.CensorConfig) error
	GetCensorConfig(ctx context.Context) (*model.CensorConfig, error)
	CreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity) error
	StopCensorJob(ctx context.Context, liveId string) error
	GetLiveCensorJobByLiveId(ctx context.Context, liveId string) (*model.LiveCensor, error)
	GetLiveCensorJobByJobId(ctx context.Context, jobId string) (*model.LiveCensor, error)

	GetCensorImageById(ctx context.Context, imageId uint) (*model.CensorImage, error)
	SaveLiveCensorJob(ctx context.Context, liveId string, jobId string, config *model.CensorConfig) error
	SaveCensorImage(ctx context.Context, image *model.CensorImage) error
	BatchUpdateCensorImage(ctx context.Context, images []uint, updates map[string]interface{}) error
	SearchCensorImage(ctx context.Context, isReview *int, pageNum, pageSize int, liveId *string) (image []model.CensorImage, totalCount int, err error)
	SearchCensorLive(ctx context.Context, isReview *int, pageNum, pageSize int) (censorLive []CensorLive, totalCount int, err error)
	GetUnreviewCount(ctx context.Context, liveId string) (len int, err error)

	ImageBucketToUrl(url string) string
}

type CensorLive struct {
	LiveId         string `json:"live_id"`
	Title          string `json:"title"`
	AnchorId       string `json:"anchor_id"`
	Nick           string `json:"nick"`
	Status         int    `json:"live_status"`
	AnchorStatus   int    `json:"anchor_status"`
	StopReason     string `json:"stop_reason"`
	StopAt         int64  `json:"stop_at"`
	StartAt        int64  `json:"start_at"`
	Count          int    `json:"count"`           //待审核次数
	ViolationCount int    `json:"violation_count"` //违规次数
	AiCount        int    `json:"ai_count"`        ///ai预警次数
	Time           int64  `json:"time"`
	PushUrl        string `json:"push_url"`
	RtmpPlayUrl    string `json:"rtmp_play_url"`
	FlvPlayUrl     string `json:"flv_play_url"`
	HlsPlayUrl     string `json:"hls_play_url"`
}
