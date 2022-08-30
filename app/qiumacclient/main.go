// @Author: wangsheng
// @Description:
// @File:  main
// @Version: 1.0.0
// @Date: 2022/5/19 6:06 下午
// Copyright 2021 QINIU. All rights reserved

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qbox/livekit/common/auth/qiniumac"
)

func main() {
	mac := &qiniumac.Mac{
		AccessKey: "ak",
		SecretKey: []byte("sk"),
	}
	transport := qiniumac.NewTransport(mac, nil)
	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Get("http://localhost:8099/server/auth/admin/token?app_id=test_app&user_id=user_1&device_id=testDevice")
	if err != nil {
		fmt.Printf("request error %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response error %v", err)
		return
	}
	fmt.Printf("got response %s", string(data))
}
