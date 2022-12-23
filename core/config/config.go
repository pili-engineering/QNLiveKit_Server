package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

	mac := &qiniumac.Mac{
		AccessKey: "ak",
		SecretKey: []byte("sk"),
	}
	tr := qiniumac.NewTransport(mac, nil)
	httpClient := &http.Client{
		Transport: tr,
	}

	client := &rpc.Client{
		Client: httpClient,
	}
	appId := "appId"
	url := "http://127.0.0.1:8080" + "/v1/app/config/cache/" + appId
	ret := &QiniuCinfig{}

	fileName := "qiniuConfig"
	err = client.GetCall(log, ret, url)
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
	v.SetConfigType("yaml")

	err = v.MergeConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("MergeConfig %s error %v", path, err)
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Errorf("write file fail error %v", err)
	}

	return &Config{
		v,
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
