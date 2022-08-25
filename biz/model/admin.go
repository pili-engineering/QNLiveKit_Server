package model

type ManagerEntity struct {
	ID          uint   `gorm:"primary_key"`
	UserId      string `json:"user_id"`
	PassWord    string `json:"pass_word"`
	Description string `json:"description"`
}

func (e ManagerEntity) TableName() string {
	return "admin_user"
}
