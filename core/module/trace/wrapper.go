package trace

import "context"

func SetRtcAppId(rtcAppId string) {
	instance.RTCAppId = rtcAppId
}

func SetPiliHub(piliHub string) {
	instance.PiliHub = piliHub
}

func SetImAppID(imAppId string) {
	instance.IMAppID = imAppId
}

func ReportEvent(ctx context.Context, kind string, event interface{}) error {
	return instance.ReportEvent(ctx, kind, event)
}

func ReportBatchEvent(ctx context.Context, kind string, events []interface{}) error {
	return instance.ReportBatchEvent(ctx, kind, events)
}
