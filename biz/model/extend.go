// @Author: wangsheng
// @Description:
// @File:  extend
// @Version: 1.0.0
// @Date: 2022/5/23 11:12 上午
// Copyright 2021 QINIU. All rights reserved

package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Extends map[string]string

//写入mysql 转换
func (c Extends) Value() (driver.Value, error) {
	data, _ := json.Marshal(c)
	return string(data), nil
}

//从mysql 读取转换
func (c *Extends) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch s := value.(type) {
	case string:
		err := json.Unmarshal([]byte(s), c)
		return err

	case []byte:
		err := json.Unmarshal(s, c)
		return err
	}
	return fmt.Errorf("can not convert %v to Extend", value)
}

//以Src 中的内容，覆盖dst 中的内容
func CombineExtends(dst, src Extends) Extends {
	e := make(map[string]string)
	if len(dst) > 0 {
		for k, v := range dst {
			e[k] = v
		}
	}

	if len(src) > 0 {
		for k, v := range src {
			e[k] = v
		}
	}

	return e
}
