package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/qiniumac"
	"github.com/qbox/livekit/utils/rpc"
)

type Config struct {
	v *viper.Viper
}

type QiniuCinfig struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      []byte `json:"data"`
}

type QiniuConfigReq struct {
	Rtc  Rtc  `json:"rtc,omitempty"`
	Pili Pili `json:"pili,omitempty"`
	Kodo Kodo `json:"kodo,omitempty"`
	Im   Im   `json:"im,omitempty"`
}
type Rtc struct {
	AppId string `json:"app_id"`
}

type Pili struct {
	Hub string `json:"hub"`
}

type Kodo struct {
	Bucket string `json:"bucket"`
}

type Im struct {
	AppId string `json:"app_id"`
}

func LoadConfig(log *logger.Logger, path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s error %v", path, err)
	}
	defer f.Close()

	v := viper.New()
	v.SetConfigType("yaml")
	err = v.ReadConfig(f)
	if err != nil {
		return nil, fmt.Errorf("read config file %s error %v", path, err)
	}

	ak := v.GetString("account.access_key")
	sk := v.GetString("account.secret_key")

	hub := v.GetString("pili.hub")
	bucket := v.GetString("kodo.bucket")
	rtc := v.GetString("rtc.app_id")
	im := v.GetString("im.app_id")
	req := &QiniuConfigReq{
		Rtc:  Rtc{AppId: rtc},
		Pili: Pili{Hub: hub},
		Kodo: Kodo{Bucket: bucket},
		Im:   Im{AppId: im},
	}

	mac := &qiniumac.Mac{
		AccessKey: ak,
		SecretKey: []byte(sk),
	}
	tr := qiniumac.NewTransport(mac, nil)
	httpClient := &http.Client{
		Transport: tr,
	}

	client := &rpc.Client{
		Client: httpClient,
	}
	url := "https://live-admin.qiniu.com" + "/v1/app/config/cache"
	ret := &QiniuCinfig{}

	fileName := "qiniuConfig"
	err = client.CallWithJSON(log, ret, url, req)
	var data []byte
	if err != nil || ret.Code != 0 {
		log.Errorf("read qiniu config file %s error %v", path, err)
		data, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("MergeConfig %s error %v", path, err)
		}
	} else {
		data = ret.Data
	}

	vip := viper.New()
	vip.SetConfigType("yaml")
	err = vip.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Errorf("write file fail error %v", err)
	}

	vip.SetConfigType("yaml")
	split := strings.Split(path, "/")
	leng := len(split[len(split)-1])
	vip.SetConfigName(path[len(path)-leng : len(path)-5])
	if len(split) == 1 {
		vip.AddConfigPath(".")
	} else {
		vip.AddConfigPath(path[:len(path)-leng])
	}
	vip.MergeInConfig()
	if err != nil {
		return nil, fmt.Errorf("MergeConfig %s error %v", path, err)
	}

	return &Config{
		vip,
	}, nil
}

func (c *Config) Sub(key string) *Config {
	if v := c.v.Sub(key); v == nil {
		return nil
	} else {
		return &Config{
			v,
		}
	}
}

func (c *Config) Unmarshal(val interface{}) error {
	return c.v.Unmarshal(val)
}
