package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/base/stats/internal/impl"
)

type IService interface {
	ReportOnlineMessage(ctx context.Context)
	PostStatsSingleLive(context.Context, []*model.StatsSingleLiveEntity) error
	GetStatsSingleLive(ctx context.Context, liveId string) (*impl.CommonStats, error)
	UpdateSingleLive(ctx context.Context, entity *model.StatsSingleLiveEntity) error

	SaveStatsSingleLive(ctx context.Context, entities []*model.StatsSingleLiveEntity) error
}

var Instance IService
