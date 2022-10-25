package trace

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type Client struct {
	Config
	client *rpc.Client
}

func newClient(conf Config) *Client {
	mac := &qiniumac.Mac{
		AccessKey: conf.AccessKey,
		SecretKey: []byte(conf.SecretKey),
	}
	tr := qiniumac.NewTransport(mac, nil)
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	r := &Client{
		Config: conf,
		client: &rpc.Client{
			Client: httpClient,
		},
	}
	return r
}

type ReportRequest struct {
	IMAppID  string      `json:"im_app"`
	RTCAppId string      `json:"rtc_app"`
	PiliHub  string      `json:"pili_hub"`
	Item     interface{} `json:"item"`
}

func (s *Client) ReportEvent(ctx context.Context, kind string, event interface{}) error {
	log := logger.ReqLogger(ctx)
	r := &ReportRequest{
		IMAppID:  s.IMAppID,
		RTCAppId: s.RTCAppId,
		PiliHub:  s.PiliHub,
		Item:     event,
	}
	url := fmt.Sprintf("%s/report/live/%s", s.ReportHost, kind)
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

func (s *Client) ReportBatchEvent(ctx context.Context, kind string, events []interface{}) error {
	log := logger.ReqLogger(ctx)
	r := &BatchReportRequest{
		IMAppID:  s.IMAppID,
		RTCAppId: s.RTCAppId,
		PiliHub:  s.PiliHub,
		Items:    events,
	}
	url := fmt.Sprintf("%s/report/live/%s/batch", s.ReportHost, kind)
	resp := &api.Response{}
	err := s.client.CallWithJSON(log, resp, url, r)
	if err != nil {
		log.Info("Error ", err)
		return err
	}
	return nil
}
