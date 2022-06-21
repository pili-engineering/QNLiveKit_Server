// @Author: wangsheng
// @Description:
// @File:  password
// @Version: 1.0.0
// @Date: 2022/5/23 4:09 下午
// Copyright 2021 QINIU. All rights reserved

package password

import "math/rand"

const passwordCharSet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()[]{}.,"

func RandomPassword(length int) string {
	charSetLen := len(passwordCharSet)
	p := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(charSetLen)
		p += passwordCharSet[idx : idx+1]
	}
	return p
}
