package model

import "github.com/qbox/livekit/utils/timestamp"

type ManagerEntity struct {
	ID          uint   `gorm:"primary_key"`
	UserName    string `json:"user_name"`
	UserId      string `json:"user_id"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func (e ManagerEntity) TableName() string {
	return "admin_user"
}

type CensorConfig struct {
	ID         uint `gorm:"primary_key"`
	Enable     bool `json:"enable"`
	Pulp       bool `json:"pulp"`
	Terror     bool `json:"terror"`
	Politician bool `json:"politician"`
	Ads        bool `json:"ads"`
	Interval   int  `json:"interval"`
}

func (e CensorConfig) TableName() string {
	return "censor_config"
}

type LiveCensor struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	LiveID     string `json:"live_id"`
	JobID      string `json:"job_id"`
	Pulp       bool   `json:"pulp"`
	Terror     bool   `json:"terror"`
	Politician bool   `json:"politician"`
	Ads        bool   `json:"ads"`
	Interval   int    `json:"interval"`
}

func (e LiveCensor) TableName() string {
	return "live_censor"
}

const (
	AuditResultPass  = 1 //审核结果为通过
	AuditResultBlock = 2 //审核结果违规
)

type CensorImage struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Url       string `json:"url"`
	JobID     string `json:"job_id"`
	CreatedAt int    `json:"created_at"`
	//Scenes       string              `json:"scenes"`

	Suggestion string `json:"suggestion"`
	Pulp       string `json:"pulp"`
	Terror     string `json:"terror"`
	Politician string `json:"politician"`
	Ads        string `json:"ads"`

	LiveID string `json:"live_id"`

	IsReview     int                  `json:"is_review"`
	ReviewAnswer int                  `json:"review_answer"`
	ReviewUserId string               `json:"review_user_id"`
	ReviewTime   *timestamp.Timestamp `json:"review_time"`
}

func (e CensorImage) TableName() string {
	return "censor_image"
}
