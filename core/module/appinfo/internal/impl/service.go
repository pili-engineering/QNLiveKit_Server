package impl

var appInfo = AppInfo{}

type AppInfo struct {
	IMAppID string `json:"im_app_id"`
}

func GetAppInfo() AppInfo {
	return appInfo
}

func SetImAppId(imAppId string) {
	appInfo.IMAppID = imAppId
}
