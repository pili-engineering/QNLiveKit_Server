package mysql

import (
	olog "log"
	"os"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/qbox/livekit/utils/logger"
)

const (
	maxOpenConns    = 16384
	connMaxLifeTime = 5
)

// Mysql mysql 实例
type Mysql struct {
	*gorm.DB
}

var (
	clientMap = map[string]*Mysql{}
)

// Init init mysql
func Init(configs ...*ConfigStructure) {

	for _, config := range configs {
		var c = &Mysql{}

		dbURI := config.GetURL()

		db, err := gorm.Open("mysql", dbURI)
		if err != nil {
			panic(err)
		}

		if err := db.DB().Ping(); err != nil {
			panic(err)
		}

		// set pool size
		if config.MaxOpenConns > 0 {
			if config.MaxOpenConns > maxOpenConns {
				config.MaxOpenConns = maxOpenConns
			}
			db.DB().SetMaxOpenConns(config.MaxOpenConns)
		}

		if config.MaxIdleConns > 0 {
			if config.MaxIdleConns > config.MaxOpenConns {
				config.MaxIdleConns = config.MaxOpenConns
			}
			db.DB().SetMaxIdleConns(config.MaxIdleConns)
		}

		// set op response timeout
		if config.ConnMaxLifeTime == 0 {
			config.ConnMaxLifeTime = connMaxLifeTime
		}

		db.DB().SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)

		//set to log
		//db.LogMode(true)

		// panic as early as possible
		if err := db.DB().Ping(); err != nil {
			panic(err)
		}

		c.DB = db
		name := config.Default
		if name == "" {
			name = "default"
		}

		if config.ReadOnly {
			name += "_readonly"
		}

		clientMap[name] = c
	}

	return
}

// Get 获取一个mysql实例 当只有一个mysql配置的是就使用这个func来获取一个mysql实例
func Get(xReqID ...string) *gorm.DB {
	return get("", xReqID...)
}

func get(name string, xReqID ...string) *gorm.DB {
	var reqID = ""
	if len(xReqID) > 0 {
		reqID = xReqID[0]
	} else {
		reqID = log.GenReqID()
	}

	if name == "" {
		name = "default"
	}
	client := clientMap[name]

	db := client.New()
	db.LogMode(true)
	db.SetLogger(Logger{olog.New(os.Stdout, "", 0), reqID})
	return db

}

func GetLive(xReqID ...string) *gorm.DB {
	return get("live", xReqID...)
}

func GetLiveReadOnly(xReqID ...string) *gorm.DB {
	return get("live_readonly", xReqID...)
}

// GetStructGormFields 返回结构体的mysql对应的column
func GetStructGormFields(st interface{}) []string {
	fields := make([]string, 0)
	scope := gorm.Scope{Value: st}
	structFields := scope.Fields()

	for _, v := range structFields {
		fields = append(fields, v.DBName)
	}

	return fields
}
