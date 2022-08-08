package report

import (
	"context"
)

type RService interface {
	ReportOnlineMessage(ctx context.Context)
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
