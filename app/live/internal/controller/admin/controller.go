package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/controller/server"
	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
	"math"
	"net/http"
	"strconv"
)

func RegisterCensorRoutes(group *gin.RouterGroup) {
	censorGroup := group.Group("/censor")
	censorGroup.POST("/config", censorController.UpdateCensorConfig)
	censorGroup.GET("/config", censorController.GetCensorConfig)
	censorGroup.POST("/job/start", censorController.CreateJob)
	censorGroup.POST("/job/close", censorController.CloseJob)
	censorGroup.POST("/job/list", censorController.ListAllJobs)
	censorGroup.POST("/job/query", censorController.QueryJob)

	censorGroup.GET("/record", censorController.SearchRecordImage)
	censorGroup.GET("/audit", censorController.AuditRecordImage)
}

var censorController = &CensorController{}

type CensorController struct {
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (c *CensorController) LoginManager(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &LoginRequest{}
	if err := ctx.BindQuery(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	manService := admin.GetManagerService()
	admin, err := manService.FindAdminByUserName(ctx, req.UserName)
	if err != nil {
		log.Errorf("userName:%s, login error:%v", req.UserName, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	} else if admin.Password != req.Password {
		log.Errorf("userName:%s, login error:%v", req.UserName, api.ErrorLoginWrong)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrorLoginWrong))
		return
	}

	authToken := token.AuthToken{
		UserId: admin.UserId,
		Role:   "admin",
	}

	tokenService := token.GetService()
	if token, err := tokenService.GenAuthToken(&authToken); err != nil {
		log.Errorf("")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		resp := &server.GetAuthTokenResponse{
			Response: api.Response{
				RequestId: log.ReqID(),
				Code:      0,
				Message:   "success",
			},
		}
		resp.Data.AccessToken = token
		resp.Data.ExpiresAt = authToken.ExpiresAt
		ctx.JSON(http.StatusOK, resp)
	}
}

type CensorConfigResponse struct {
	api.Response
	Data *dto.CensorConfigDto `json:"data"`
}

func (c *CensorController) UpdateCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &dto.CensorConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	err := censorService.UpdateCensorConfig(ctx, dto.CConfigDtoToEntity(req))
	if err != nil {
		log.Errorf(" UpdateCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: req,
	})
}

func (c *CensorController) GetCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	censorService := admin.GetCensorService()
	censorConfig, err := censorService.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: dto.CConfigEntityToDto(censorConfig),
	})
}

func (c *CensorController) CallbackCensorJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCallBack{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if req.Image.Code != 200 {
		fmt.Println(req)
		return
	}
	m := &model.CensorImage{
		JobID:      req.Image.Job,
		Url:        req.Image.Url,
		CreatedAt:  req.Image.Timestamp,
		Suggestion: req.Image.Result.Suggestion,
		Pulp:       req.Image.Result.Scenes.Pulp.Suggestion,
		Ads:        req.Image.Result.Scenes.Ads.Suggestion,
		Politician: req.Image.Result.Scenes.Politician.Suggestion,
		Terror:     req.Image.Result.Scenes.Terror.Suggestion,
	}

	censorService := admin.GetCensorService()
	liveCensor, err := censorService.GetLiveCensorJobByJobId(ctx, req.Image.Job)
	if err != nil {
		log.Errorf("GetLiveCensorJobByJobId error %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	m.LiveID = liveCensor.LiveID
	err = censorService.SetCensorImage(ctx, m)
	if err != nil {
		log.Errorf("SetCensorImage error %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "switch mic failed",
			RequestId: log.ReqID(),
		})
		return
	}
	ctx.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (c *CensorController) CreateJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCreateRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	m, err := live.GetService().LiveInfo(ctx, req.LiveId)
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	censorService := admin.GetCensorService()
	err = censorService.CreateCensorJob(ctx, m)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &api.Response{
		RequestId: log.ReqID(),
		Code:      0,
		Message:   "success",
	})

}

func (c *CensorController) CloseJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCloseRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	err := censorService.StopCensorJob(ctx, req.LiveId)
	if err != nil {
		log.Errorf("Stop censor job  error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &api.Response{
		RequestId: log.ReqID(),
		Code:      0,
		Message:   "success",
	})
}

func (c *CensorController) SearchRecordImage(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	IsReview := ctx.DefaultQuery("is_review", "2")
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	liveId := ctx.Query("live_id")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	IsReviewInt, err := strconv.Atoi(IsReview)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "is_review is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	//0： 没审核 1：审核 2：都需要list出来/*
	images, count, err := admin.GetCensorService().SearchCensorImage(ctx, IsReviewInt, pageNumInt, pageSizeInt, liveId)
	if err != nil {
		log.Errorf("search censor image  failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search  censor image failed",
			RequestId: log.ReqID(),
		})
		return
	}

	endPage := false
	if len(images) < pageSizeInt {
		endPage = true
	}
	response := &CensorImageListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = images
	ctx.JSON(http.StatusOK, response)
}

func (c *CensorController) ListAllJobs(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	requset := &admin.JobListRequest{
		Start:  req.Start.Unix(),
		End:    req.End.Unix(),
		Status: req.Status,
		Limit:  req.Limit,
		Marker: req.Marker,
	}
	resp := &admin.JobListResponse{}
	err := admin.GetCensorService().JobList(ctx, requset, resp)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *CensorController) QueryJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorQueryRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	requset := &admin.JobQueryRequest{
		Start:       req.Start.Unix(),
		End:         req.End.Unix(),
		Job:         req.Job,
		Suggestions: req.Suggestions,
	}
	resp := &admin.JobQueryResponse{}
	err := admin.GetCensorService().JobQuery(ctx, requset, resp)
	if err != nil {
		log.Errorf("Job Query error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *CensorController) AuditRecordImage(ctx *gin.Context) {
	userInfo := ctx.MustGet(liveauth.AdminCtxKey).(*liveauth.AdminInfo)
	log := logger.ReqLogger(ctx)
	req := &CensorAuditRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if len(req.Images) == 0 {
		log.Errorf("invalid request %+v", req)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	for _, idx := range req.Images {
		image, err := censorService.GetCensorImageById(ctx, idx)
		if err != nil {
			log.Errorf("AuditRecordImage fail %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}
		if image.IsReview == 1 {
			continue
		}
		image.IsReview = 1
		image.ReviewAnswer = req.ReviewAnswer
		image.ReviewTime = timestamp.Now()
		image.ReviewUserId = userInfo.UserId
		err = censorService.SetCensorImage(ctx, image)
		if err != nil {
			log.Errorf("AuditRecordImage fail %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}
	}
	ctx.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

type CensorCreateRequest struct {
	LiveId string `json:"live_id"`
}

type CensorAuditRequest struct {
	Images       []uint `json:"image_list"`
	ReviewAnswer int    `json:"review_answer"`
}

type CensorListRequest struct {
	Start  timestamp.Timestamp `json:"start"`
	End    timestamp.Timestamp `json:"end"`
	Status string              `json:"status"`
	Limit  int                 `json:"limit"`
	Marker string              `json:"marker"`
}

type CensorCloseRequest struct {
	LiveId string `json:"live_id"`
}

type CensorQueryRequest struct {
	Start       timestamp.Timestamp `json:"start"`
	End         timestamp.Timestamp `json:"end"`
	Job         string              `json:"job"`
	Suggestions []string            `json:"suggestions"`
}

type CensorCallBack struct {
	Job  string `json:"job"`
	Live struct {
		Id   string `json:"id"`
		Uri  string `json:"uri"`
		Info string `json:"info"`
	} `json:"live"`
	Error struct {
		Timestamp int    `json:"timestamp"`
		Message   string `json:"message"`
	} `json:"error"`
	Image struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Job       string `json:"job"`
		Url       string `json:"url"`
		Timestamp int    `json:"timestamp"`
		Result    struct {
			Suggestion string `json:"suggestion"`
			Scenes     struct {
				Ads struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"ads"`
				Politician struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"politician"`
				Pulp struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"pulp"`
				Terror struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"terror"`
			} `json:"scenes"`
		} `json:"result"`
	} `json:"image"`
	Audio struct {
		Result struct {
			Suggestion string `json:"suggestion"`
		} `json:"result"`
	} `json:"audio"`
}

type CensorImageListResponse struct {
	api.Response
	Data struct {
		TotalCount int                 `json:"total_count"`
		PageTotal  int                 `json:"page_total"`
		EndPage    bool                `json:"end_page"`
		List       []model.CensorImage `json:"list"`
	} `json:"data"`
}
