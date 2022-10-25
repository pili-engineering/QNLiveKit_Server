package auth

const UserCtxKey = "UserCtxKey"
const AdminCtxKey = "AdminCtxKey"

type AdminInfo struct {
	UserId string
	Role   string
}

type UserInfo struct {
	AppId    string
	UserId   string
	DeviceId string
	ImUserId int64
}
