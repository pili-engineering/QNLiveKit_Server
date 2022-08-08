package report

import (
	"context"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
	"net/http"
	"strconv"
)

type RClient struct {
	Config
	client *http.Client
}

func NewReportClient(conf Config) *RClient {
	mac := &qiniumac.Mac{
		AccessKey: conf.AccessKey,
		SecretKey: []byte(conf.SecretKey),
	}
	tr := qiniumac.NewTransport(mac, nil)
	r := &RClient{
		Config: conf,
		client: &http.Client{
			Transport: tr,
		},
	}
	return r
}

type RequestReport struct {
	IMAppID  string        `json:"im_app"`
	RTCAppId string        `json:"rtc_app"`
	PiliHub  string        `json:"pili_hub"`
	Item     model.Extends `json:"item"`
}

func (s *RClient) ReportOnlineMessage(ctx context.Context) {
	var log *logger.Logger
	if ctx == nil {
		log = logger.New("ReportOnlineMessage Start")
	} else {
		log = logger.ReqLogger(ctx)
	}
	db := mysql.GetLive(log.ReqID())
	var numLives int
	db.Table("live_entities").Count(&numLives)
	var numUsers int
	db.Table("live_users").Count(&numUsers)
	var value = map[string]string{
		"lives": strconv.Itoa(numLives),
		"users": strconv.Itoa(numUsers),
	}
	s.postReportJson(log, "overview", value)
}

func (s *RClient) postReportJson(log *logger.Logger, kind string, value map[string]string) error {
	r := &RequestReport{
		IMAppID:  s.IMAppID,
		RTCAppId: s.RTCAppId,
		PiliHub:  s.PiliHub,
		Item:     value,
	}
	client := rpc.Client{
		Client: s.client,
	}
	url := s.ReportHost + "/report/live/" + kind
	resp := &api.Response{}
	err := client.CallWithJSON(log, resp, url, r)
	if err != nil {
		log.Info("Error ", err)
		return err
	}
	return nil
}
