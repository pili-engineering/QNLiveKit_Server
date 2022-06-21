// @Author: wangsheng
// @Description:
// @File:  password_test.go
// @Version: 1.0.0
// @Date: 2022/5/23 4:22 下午
// Copyright 2021 QINIU. All rights reserved

package password

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomPassword(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{
			name:   "6",
			length: 6,
			want:   6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomPassword(tt.length)
			assert.Equal(t, tt.want, len(got))
		})
	}
}
