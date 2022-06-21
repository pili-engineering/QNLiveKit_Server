// @Author: wangsheng
// @Description:
// @File:  config
// @Version: 1.0.0
// @Date: 2022/5/19 10:28 上午
// Copyright 2021 QINIU. All rights reserved

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig(confPath string, conf interface{}) error {
	viper.SetConfigType("yaml")

	if confPath == "" {
		return fmt.Errorf("need config file")
	}

	f, err := os.Open(confPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = viper.ReadConfig(f); err != nil {
		return err
	}

	if err = viper.Unmarshal(conf); err != nil {
		return err
	}
	return nil
}
