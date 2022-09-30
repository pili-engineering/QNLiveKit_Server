package trace

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/qiniumac"
	"github.com/qbox/livekit/utils/rpc"
)

var instance = &Service{}

type Service struct {
	IMAppID  string
	RTCAppId string
	PiliHub  string

	client *rpc.Client
}

func newRpcClient(ak, sk string) *rpc.Client {
	mac := &qiniumac.Mac{
		AccessKey: ak,
		SecretKey: []byte(sk),
	}
	tr := qiniumac.NewTransport(mac, nil)
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	return &rpc.Client{
		Client: httpClient,
	}
}

type ReportRequest struct {
	IMAppID  string      `json:"im_app"`
	RTCAppId string      `json:"rtc_app"`
	PiliHub  string      `json:"pili_hub"`
	Item     interface{} `json:"item"`
}

const reportHost = "https://niucube-api.qiniu.com"

func (s *Service) ReportEvent(ctx context.Context, kind string, event interface{}) error {
	log := logger.ReqLogger(ctx)
	r := &ReportRequest{
		IMAppID:  s.IMAppID,
		RTCAppId: s.RTCAppId,
		PiliHub:  s.PiliHub,
		Item:     event,
	}
	url := fmt.Sprintf("%s/report/live/%s", reportHost, kind)
	resp := &api.Response{}
	err := s.client.CallWithJSON(log, resp, url, r)
	if err != nil {
		log.Info("Error ", err)
		return err
	}
	return nil
}

type BatchReportRequest struct {
	IMAppID  string        `json:"im_app"`
	RTCAppId string        `json:"rtc_app"`
	PiliHub  string        `json:"pili_hub"`
	Items    []interface{} `json:"items"`
}

func (s *Service) ReportBatchEvent(ctx context.Context, kind string, events []interface{}) error {
	log := logger.ReqLogger(ctx)
	r := &BatchReportRequest{
		IMAppID:  s.IMAppID,
		RTCAppId: s.RTCAppId,
		PiliHub:  s.PiliHub,
		Items:    events,
	}
	url := fmt.Sprintf("%s/report/live/%s/batch", reportHost, kind)
	resp := &api.Response{}
	err := s.client.CallWithJSON(log, resp, url, r)
	if err != nil {
		log.Info("Error ", err)
		return err
	}
	return nil
}
