package dto

import "github.com/qbox/livekit/biz/model"

type CensorConfigDto struct {
	Enable     bool `json:"enable"`
	Pulp       bool `json:"pulp"`
	Terror     bool `json:"terror"`
	Politician bool `json:"politician"`
	Ads        bool `json:"ads"`
	Interval   int  `json:"interval"`
}

func CConfigDtoToEntity(dto *CensorConfigDto) *model.CensorConfig {
	return &model.CensorConfig{
		Enable:     dto.Enable,
		Pulp:       dto.Pulp,
		Terror:     dto.Terror,
		Politician: dto.Politician,
		Ads:        dto.Ads,
		Interval:   dto.Interval,
	}
}

func CConfigEntityToDto(entity *model.CensorConfig) *CensorConfigDto {
	return &CensorConfigDto{
		Enable:     entity.Enable,
		Pulp:       entity.Pulp,
		Terror:     entity.Terror,
		Politician: entity.Politician,
		Ads:        entity.Ads,
		Interval:   entity.Interval,
	}
}

type CensorImageDto struct {
	ID        uint   `json:"id"`
	Url       string `json:"url"`
	JobID     string `json:"job_id"`
	CreatedAt int64  `json:"created_at"`

	Suggestion string `json:"suggestion"`
	Pulp       string `json:"pulp"`
	Terror     string `json:"terror"`
	Politician string `json:"politician"`
	Ads        string `json:"ads"`

	LiveID string `json:"live_id"`

	IsReview     int    `json:"is_review"`
	ReviewAnswer int    `json:"review_answer"`
	ReviewUserId string `json:"review_user_id"`
	ReviewTime   int64  `json:"review_time"`
}

func CensorImageModelToDto(entity *model.CensorImage) *CensorImageDto {
	return &CensorImageDto{
		ID:           entity.ID,
		Url:          entity.Url,
		JobID:        entity.JobID,
		CreatedAt:    int64(entity.CreatedAt * 1000),
		Suggestion:   entity.Suggestion,
		Pulp:         entity.Pulp,
		Terror:       entity.Terror,
		Politician:   entity.Politician,
		Ads:          entity.Ads,
		LiveID:       entity.LiveID,
		IsReview:     entity.IsReview,
		ReviewAnswer: entity.ReviewAnswer,
		ReviewUserId: entity.ReviewUserId,
		ReviewTime:   entity.ReviewTime.UnixMilli() / 1000,
	}
}
