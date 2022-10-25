package trace

import "context"

type Config struct {
	IMAppID    string
	RTCAppId   string
	PiliHub    string
	AccessKey  string
	SecretKey  string
	ReportHost string
}

var client *Client

func InitService(conf Config) error {
	client = newClient(conf)
	return nil
}

func ReportEvent(ctx context.Context, kind string, event interface{}) error {
	return client.ReportEvent(ctx, kind, event)
}

func ReportBatchEvent(ctx context.Context, kind string, events []interface{}) error {
	return client.ReportBatchEvent(ctx, kind, events)
}
