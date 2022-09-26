package cron

import (
	"gopkg.in/robfig/cron.v2"
)

func AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return instance.AddSingleTaskFunc(spec, cmd)
}

func AddSingleTaskFunc(spec string, cmd func()) (cron.EntryID, error) {
	return instance.AddSingleTaskFunc(spec, cmd)
}
