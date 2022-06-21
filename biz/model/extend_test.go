// @Author: wangsheng
// @Description:
// @File:  extend_test.go
// @Version: 1.0.0
// @Date: 2022/5/23 8:51 下午
// Copyright 2021 QINIU. All rights reserved

package model

import (
	"reflect"
	"testing"
)

func TestCombineExtends(t *testing.T) {

	tests := []struct {
		name string
		dst  Extends
		src  Extends
		want Extends
	}{
		{
			name: "src nil",
			src:  nil,
			dst:  map[string]string{"a": "a1", "b": "b1"},
			want: map[string]string{"a": "a1", "b": "b1"},
		},
		{
			name: "dst nil",
			src:  map[string]string{"a": "a1", "b": "b1"},
			dst:  nil,
			want: map[string]string{"a": "a1", "b": "b1"},
		},
		{
			name: "dst src",
			src:  map[string]string{"a": "a11", "c": "c1"},
			dst:  map[string]string{"a": "a1", "b": "b1"},
			want: map[string]string{"a": "a11", "b": "b1", "c": "c1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CombineExtends(tt.dst, tt.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CombineExtends() = %v, want %v", got, tt.want)
			}
		})
	}
}
