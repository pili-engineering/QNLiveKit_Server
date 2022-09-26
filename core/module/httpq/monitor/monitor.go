package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/qbox/livekit/common/trace"
)

type MonitorItem struct {
	Method   string `json:"method"`
	Host     string `json:"host"`
	Path     string `json:"path"`
	Handler  string `json:"handler"`
	Status   int    `json:"status"`
	Duration int    `json:"duration"`
	LogTime  int64  `json:"logTime"`
}

var (
	monitorItems     []*MonitorItem
	monitorItemsLock sync.Mutex
)

func init() {
	monitorItems = make([]*MonitorItem, 0, 1024)
}

func monitor(method, host, path, handler string, status int, duration int) {
	item := MonitorItem{
		Method:   method,
		Host:     host,
		Path:     path,
		Handler:  handler,
		Status:   status,
		Duration: duration,
		LogTime:  time.Now().UnixNano() / int64(time.Millisecond),
	}

	addMonitorItem(&item)
}

func addMonitorItem(item *MonitorItem) {
	monitorItemsLock.Lock()
	defer monitorItemsLock.Unlock()

	monitorItems = append(monitorItems, item)
}

// AllAndResetMonitorItems 获取当前所有的监控记录，并清空
func allAndResetMonitorItems() []*MonitorItem {
	monitorItemsLock.Lock()
	defer monitorItemsLock.Unlock()

	ret := monitorItems
	monitorItems = make([]*MonitorItem, 0, len(monitorItems))
	return ret
}

func ReportMonitorItems(ctx context.Context) {
	items := allAndResetMonitorItems()
	if len(items) == 0 {
		return
	}

	events := make([]interface{}, 0, len(items))
	for _, item := range items {
		events = append(events, item)
	}

	trace.ReportBatchEvent(ctx, "api", events)
}
