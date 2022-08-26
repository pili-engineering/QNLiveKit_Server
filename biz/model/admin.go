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
