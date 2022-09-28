package admin

import (
	"context"
	"net/http"
	"strings"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type JobService interface {
	JobQuery(ctx context.Context, req *JobQueryRequest, resp *JobQueryResponse) error
	JobCreate(ctx context.Context, liveEntity *model.LiveEntity, config *model.CensorConfig) (*JobCreateResponse, error)
	JobClose(ctx context.Context, req *JobCreateResponseData) error
	JobList(ctx context.Context, req *JobListRequest, resp *JobListResponse) error
	ImageBucketToUrl(url string) string
}

var service JobService

func InitJobService(config Config) {
	service = NewCensorClient(config)
}

func GetJobService() JobService {
	return service
}

type Config struct {
	AccessKey      string
	SecretKey      string
	CensorCallback string
	CensorBucket   string
	CensorAddr     string
}

type CensorClient struct {
	AccessKey      string
	SecretKey      string
	CensorCallback string
	CensorBucket   string
	CensorAddr     string
	client         *rpc.Client
}

func NewCensorClient(config Config) *CensorClient {
	mac := &qiniumac.Mac{
		AccessKey: config.AccessKey,
		SecretKey: []byte(config.SecretKey),
	}
	c := &http.Client{
		Transport: qiniumac.NewTransport(mac, nil),
	}
	return &CensorClient{
		AccessKey:      config.AccessKey,
		SecretKey:      config.AccessKey,
		CensorBucket:   config.CensorBucket,
		CensorCallback: config.CensorCallback,
		CensorAddr:     config.CensorAddr,
		client: &rpc.Client{
			Client: c,
		},
	}
}

func (c *CensorClient) ImageBucketToUrl(url string) string {
	split := strings.Split(url, c.CensorBucket)
	return c.CensorAddr + split[1]
}

func (c *CensorClient) JobCreate(ctx context.Context, liveEntity *model.LiveEntity, config *model.CensorConfig) (*JobCreateResponse, error) {
	log := logger.ReqLogger(ctx)
	req := &JobCreateRequest{}
	req.Data.Url = liveEntity.PushUrl
	req.Params.Image.IsOn = config.Enable
	req.Params.Image.IntervalMsecs = config.Interval * 1000

	req.Params.HookAuth = false
	req.Params.HookUrl = c.CensorCallback + "/manager/censor/callback"

	s := make([]string, 0)
	if config.Pulp {
		s = append(s, "pulp")
	}
	if config.Terror {
		s = append(s, "terror")
	}
	if config.Politician {
		s = append(s, "politician")
	}
	if config.Ads {
		s = append(s, "ads")
	}

	req.Params.Image.Scenes = s
	req.Params.Image.HookRule = 0 //图片审核结果回调规则，0/1。默认为 0，返回判定结果违规的审核结果；设为 1 时，返回所有审核结果。
	req.Params.Image.Saver.Bucket = c.CensorBucket
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor"
	resp := &JobCreateResponse{}
	err := c.client.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return resp, nil
}

func (c *CensorClient) JobClose(ctx context.Context, req *JobCreateResponseData) error {
	log := logger.ReqLogger(ctx)
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/close"

	resp := &api.Response{}
	err := c.client.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *CensorClient) JobQuery(ctx context.Context, req *JobQueryRequest, resp *JobQueryResponse) error {
	log := logger.ReqLogger(ctx)
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/query"
	err := c.client.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *CensorClient) JobList(ctx context.Context, req *JobListRequest, resp *JobListResponse) error {
	log := logger.ReqLogger(ctx)
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/list"
	err := c.client.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		return err
	}
	return nil
}
