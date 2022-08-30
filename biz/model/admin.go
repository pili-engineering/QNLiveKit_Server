package model

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
