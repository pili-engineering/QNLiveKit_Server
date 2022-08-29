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
