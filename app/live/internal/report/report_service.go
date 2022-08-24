package report

import (
	"context"
	"github.com/qbox/livekit/biz/model"
)

type RService interface {
	ReportOnlineMessage(ctx context.Context)
	PostStatsSingleLive(context.Context, []*model.StatsSingleLiveEntity) error
	GetStatsSingleLive(ctx context.Context, liveId string) (*CommonStats, error)
}

type Config struct {
	IMAppID    string
	RTCAppId   string
	PiliHub    string
	AccessKey  string
	SecretKey  string
	ReportHost string
}

var service RService

func GetService() RService {
	return service
}

func InitService(config Config) {
	service = NewReportClient(config)
}
