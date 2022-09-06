package admin

import (
	"context"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
	"net/http"
)

type CCService interface {
	UpdateCensorConfig(ctx context.Context, mod *model.CensorConfig) error
	GetCensorConfig(ctx context.Context) (*model.CensorConfig, error)
	CreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity) error
	StopCensorJob(ctx context.Context, liveId string) error
	GetLiveCensorJobByLiveId(ctx context.Context, liveId string) (*model.LiveCensor, error)
	GetLiveCensorJobByJobId(ctx context.Context, jobId string) (*model.LiveCensor, error)

	GetCensorImageById(ctx context.Context, imageId uint) (*model.CensorImage, error)

	SetCensorImage(ctx context.Context, image *model.CensorImage) error
	SearchCensorImage(ctx context.Context, isReview, pageNum, pageSize int, liveId string) (image []model.CensorImage, totalCount int, err error)

	JobList(ctx context.Context, req *JobListRequest, resp *JobListResponse) error
	JobQuery(ctx context.Context, req *JobQueryRequest, resp *JobQueryResponse) error
}

type Config struct {
	AccessKey      string
	SecretKey      string
	CensorCallback string
}

type CensorService struct {
	Config
	CClient rpc.Client
}

func InitCensorService(config Config) {
	mac := &qiniumac.Mac{
		AccessKey: config.AccessKey,
		SecretKey: []byte(config.SecretKey),
	}
	c := &http.Client{
		Transport: qiniumac.NewTransport(mac, nil),
	}
	cService = &CensorService{
		Config: config,
		CClient: rpc.Client{
			Client: c,
		},
	}
}

var cService CCService

func GetCensorService() CCService {
	return cService
}

func (c *CensorService) UpdateCensorConfig(ctx context.Context, mod *model.CensorConfig) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod.ID = 1
	err := db.Model(model.CensorConfig{}).Save(mod).Error
	if err != nil {
		return api.ErrDatabase
	}
	return nil
}

func (c *CensorService) GetCensorConfig(ctx context.Context) (*model.CensorConfig, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod := &model.CensorConfig{}
	result := db.Model(model.CensorConfig{}).First(mod)
	if result.Error != nil {
		if result.RecordNotFound() {
			mod.Enable = true
			mod.ID = 1
			mod.Pulp = true
			mod.Interval = 20
			err := db.Model(model.CensorConfig{}).Save(mod).Error
			if err != nil {
				return nil, api.ErrDatabase
			}
		} else {
			return nil, api.ErrDatabase
		}
	}
	return mod, nil
}

type JobCreateRequest struct {
	Data   JobLiveData     `json:"data"`
	Params JobCreateParams `json:"params"`
}

type JobLiveData struct {
	ID   string `json:"ID"`
	Url  string `json:"uri"`
	Info string `json:"info"`
}

type JobCreateParams struct {
	Image    JobImage `json:"image"`
	HookUrl  string   `json:"hook_url"`
	HookAuth bool     `json:"hook_auth"`
}

type JobImage struct {
	IsOn          bool          `json:"is_on"`
	Scenes        []string      `json:"scenes"`
	IntervalMsecs int           `json:"interval_msecs"`
	Saver         JobImageSaver `json:"saver"`
	HookRule      int           `json:"hook_rule"`
}

type JobImageSaver struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

type JobCreateResponse struct {
	RequestId string                `json:"request_id"` //请求ID
	Code      int                   `json:"code"`       //错误码，0 成功，其他失败
	Message   string                `json:"message"`    //错误信息
	Data      JobCreateResponseData `json:"data"`
}

type JobCreateResponseData struct {
	JobID string `json:"job"`
}

func (c *CensorService) CreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity) error {
	log := logger.ReqLogger(ctx)
	config, err := c.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig Error %v", err)
		return err
	}
	if config.Enable == false {
		return nil
	}
	resp, err := c.postCreateCensorJob(ctx, liveEntity, config)
	if err != nil {
		log.Errorf("postCreateCensorJob Error %v", err)
		return err
	}
	err = c.SetLiveCensorJob(ctx, liveEntity.LiveId, resp.Data.JobID, config)
	if err != nil {
		log.Errorf("SetLiveCensorJob Error %v", err)
		return nil
	}
	return nil
}

func (c *CensorService) postCreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity, config *model.CensorConfig) (*JobCreateResponse, error) {
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
	req.Params.Image.Saver.Bucket = "niu-cube"
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor"
	resp := &JobCreateResponse{}
	err := c.CClient.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return resp, nil
}

func (c *CensorService) SetLiveCensorJob(ctx context.Context, liveId string, JobId string, config *model.CensorConfig) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{
		LiveID:     liveId,
		JobID:      JobId,
		Interval:   config.Interval,
		Politician: config.Politician,
		Pulp:       config.Pulp,
		Ads:        config.Ads,
		Terror:     config.Terror,
	}
	err := db.Model(model.LiveCensor{}).Save(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CensorService) GetLiveCensorJobByLiveId(ctx context.Context, liveId string) (*model.LiveCensor, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{}
	err := db.Model(model.LiveCensor{}).First(m, "live_id = ?", liveId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) StopCensorJob(ctx context.Context, liveId string) error {
	log := logger.ReqLogger(ctx)
	liveCensorJob, err := c.GetLiveCensorJobByLiveId(ctx, liveId)
	if err != nil {
		log.Errorf("GetLiveCensorJobByLiveId Error %v", err)
		return err
	}
	req := &JobCreateResponseData{
		JobID: liveCensorJob.JobID,
	}
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/close"
	resp := &api.Response{}
	err = c.CClient.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		log.Errorf("post StopCensorJob Error %v", err)
		return err
	}
	return nil
}

func (c *CensorService) GetCensorImageById(ctx context.Context, imageId uint) (*model.CensorImage, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.CensorImage{}
	err := db.Model(model.CensorImage{}).First(m, "id = ?", imageId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) GetLiveCensorJobByJobId(ctx context.Context, jobId string) (*model.LiveCensor, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{}
	err := db.Model(model.LiveCensor{}).First(m, "job_id = ?", jobId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) SetCensorImage(ctx context.Context, image *model.CensorImage) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(model.CensorImage{}).Save(image).Error
	if err != nil {
		return err
	}
	return nil
}

// SearchCensorImage 0： 没审核 1：审核 2：都需要list出来/*
func (c *CensorService) SearchCensorImage(ctx context.Context, isReview, pageNum, pageSize int, liveId string) (image []model.CensorImage, totalCount int, err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLiveReadOnly(log.ReqID())
	image = make([]model.CensorImage, 0)

	if liveId == "" {
		if isReview == 0 {
			err = db.Where(" is_review = ? ", isReview).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
			err = db.Model(&model.CensorImage{}).Where(" is_review = ? ", isReview).Count(&totalCount).Error
		} else if isReview == 1 {
			err = db.Model(&model.CensorImage{}).Where(" is_review = ? ", isReview).Count(&totalCount).Error
			err = db.Where(" is_review = ?", isReview).Order("review_answer desc").Order("created_at desc ").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
		} else {
			err = db.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
			err = db.Model(&model.CensorImage{}).Count(&totalCount).Error
		}
	} else {
		if isReview == 0 {
			err = db.Where(" is_review = ? and live_id = ?", isReview, liveId).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
			err = db.Model(&model.CensorImage{}).Where(" is_review = ?  and live_id = ?", isReview, liveId).Count(&totalCount).Error
		} else if isReview == 1 {
			err = db.Model(&model.CensorImage{}).Where(" is_review = ? and live_id = ?", isReview, liveId).Count(&totalCount).Error
			err = db.Where(" is_review = ? and live_id = ?", isReview, liveId).Order("review_answer desc").Order("created_at desc ").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
		} else {
			err = db.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
			err = db.Model(&model.CensorImage{}).Count(&totalCount).Error
		}
	}

	if err != nil {
		log.Errorf("search CensorImage error: %v", err)
		return
	}
	return
}

type JobListRequest struct {
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Marker string `json:"marker"`
}

type JobListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Marker string `json:"marker"`
		Items  []struct {
			Id   string `json:"id"`
			Data struct {
				Id   string `json:"id"`
				Uri  string `json:"uri"`
				Info string `json:"info"`
			} `json:"data"`
			Params struct {
				HookUrl  string `json:"hook_url"`
				HookAuth bool   `json:"hook_auth"`
				Image    struct {
					IsOn          bool     `json:"is_on"`
					Scenes        []string `json:"scenes"`
					IntervalMsecs int      `json:"interval_msecs"`
					Saver         struct {
						Uid    int    `json:"uid"`
						Bucket string `json:"bucket"`
						Prefix string `json:"prefix"`
					} `json:"saver"`
					HookRule int `json:"hook_rule"`
				} `json:"image"`
			} `json:"params"`
			Message   string `json:"message"`
			Status    string `json:"status"`
			CreatedAt int    `json:"created_at"`
			UpdatedAt int    `json:"updated_at"`
		} `json:"items"`
	} `json:"data"`
}

func (c *CensorService) JobList(ctx context.Context, req *JobListRequest, resp *JobListResponse) error {
	log := logger.ReqLogger(ctx)
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/list"
	err := c.CClient.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		return err
	}
	return nil
}

type JobQueryRequest struct {
	Job         string   `json:"job"`
	Suggestions []string `json:"suggestions"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
}

type JobQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Marker string `json:"marker"`
		Items  struct {
			Image []struct {
				Code      int    `json:"code,omitempty"`
				Message   string `json:"message,omitempty"`
				Job       string `json:"job,omitempty"`
				Timestamp int    `json:"timestamp,omitempty"`
				Url       string `json:"url,omitempty"`
				Result    struct {
					Suggestion string `json:"suggestion"`
					Scenes     struct {
						Pulp struct {
							Suggestion string `json:"suggestion"`
							Details    []struct {
								Suggestion string  `json:"suggestion,omitempty"`
								Label      string  `json:"label,omitempty"`
								Score      float64 `json:"score,omitempty"`
							} `json:"details"`
						} `json:"pulp"`
					} `json:"scenes"`
				} `json:"result,omitempty"`
			} `json:"image"`
			Audio []struct {
				Code      int    `json:"code,omitempty"`
				Message   string `json:"message,omitempty"`
				Job       string `json:"job,omitempty"`
				Start     int    `json:"start,omitempty"`
				End       int    `json:"end,omitempty"`
				Url       string `json:"url,omitempty"`
				AudioText string `json:"audio_text,omitempty"`
				Result    struct {
					Suggestion string `json:"suggestion"`
					Scenes     struct {
						Antispam struct {
							Suggestion string `json:"suggestion"`
							Details    []struct {
								Suggestion string  `json:"suggestion"`
								Label      string  `json:"label"`
								Text       string  `json:"text"`
								Score      float64 `json:"score"`
							} `json:"details"`
						} `json:"antispam"`
					} `json:"scenes"`
				} `json:"result,omitempty"`
			} `json:"audio"`
		} `json:"items"`
	} `json:"data"`
}

func (c *CensorService) JobQuery(ctx context.Context, req *JobQueryRequest, resp *JobQueryResponse) error {
	log := logger.ReqLogger(ctx)
	reqUrl := "http://ai.qiniuapi.com/v3/live/censor/query"
	err := c.CClient.CallWithJSON(log, resp, reqUrl, req)
	if err != nil {
		return err
	}
	return nil
}
