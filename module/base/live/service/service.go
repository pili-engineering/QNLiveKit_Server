package service

import (
	"context"
	"time"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/utils/timestamp"
)

var Instance IService

type IService interface {
	CreateLive(context context.Context, req *CreateLiveRequest) (live *model.LiveEntity, err error)

	GetLiveAuthor(ctx context.Context, liveId string) (*model.LiveUserEntity, error)

	DeleteLive(context context.Context, liveId string, anchorId string) (err error)

	StartLive(context context.Context, liveId string, anchorId string) (roomToken string, err error)

	StopLive(context context.Context, liveId string, anchorId string) (err error)

	AdminStopLive(ctx context.Context, liveId string, reason string, adminId string) error

	LiveInfo(context context.Context, liveId string) (live *model.LiveEntity, err error)

	LiveListAnchor(context context.Context, pageNum, pageSize int, anchorId string) (lives []model.LiveEntity, totalCount int, err error)

	LiveList(context context.Context, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error)

	LiveUserList(context context.Context, liveId string, pageNum, pageSize int) (users []model.LiveRoomUserEntity, totalCount int, err error)

	UpdateExtends(context context.Context, liveId string, extends model.Extends) (err error)

	JoinLiveRoom(context context.Context, liveId string, userId string) (err error)

	LeaveLiveRoom(context context.Context, liveId string, userId string) (err error)

	SearchLive(context context.Context, keyword string, flag, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error)

	// CurrentLiveRoom 查找主播当前所在的直播间
	CurrentLiveRoom(ctx context.Context, userId string) (liveEntity *model.LiveEntity, err error)

	Heartbeat(context context.Context, liveId string, userId string) (liveEntity *model.LiveEntity, err error)

	// StartRelay 绑定直播间与跨房PK 会话
	StartRelay(ctx context.Context, roomId, userId string, sid string) (err error)

	// StopRelay 解绑指定的跨房PK 会话
	StopRelay(ctx context.Context, roomId, userId string, sid string) (err error)

	TimeoutLiveUser(ctx context.Context, now time.Time)

	TimeoutLiveRoom(ctx context.Context, now time.Time)

	FindLiveRoomUser(context context.Context, liveId string, userId string) (liveRoomUser *model.LiveRoomUserEntity, err error)

	CheckLiveAnchor(ctx context.Context, liveId string, userId string) error

	UpdateLiveRelatedReview(context context.Context, liveId string, latest *int) (err error)

	AddLike(ctx context.Context, liveId string, userId string, count int64) (my, total int64, err error)

	FlushCacheLikes(ctx context.Context)
}

type CreateLiveRequest struct {
	AnchorId        string               `json:"anchor_id"`
	Title           string               `json:"title"`
	Notice          string               `json:"notice"`
	CoverUrl        string               `json:"cover_url"`
	StartAt         timestamp.Timestamp  `json:"start_at"`
	EndAt           timestamp.Timestamp  `json:"end_at"`
	PublishExpireAt *timestamp.Timestamp `json:"publish_expire_at"`
	Extends         model.Extends        `json:"extends" gorm:"type:varchar(512)"`
}
