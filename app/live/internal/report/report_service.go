package report

import (
	"context"
	"github.com/qbox/livekit/biz/model"
)

type RService interface {
	ReportOnlineMessage(ctx context.Context)
	PostStatsSingleLive(context.Context, []*model.StatsSingleLiveEntity) error
	GetStatsSingleLive(ctx context.Context, liveId string) (*CommonStats, error)
	UpdateSingleLive(ctx context.Context, entity *model.StatsSingleLiveEntity) error
}

var service RService

func GetService() RService {
	return service
}

func InitService() {
	service = NewReportClient()
}
