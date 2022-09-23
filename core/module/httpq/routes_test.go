package httpq

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServer_selectGroup(t *testing.T) {
	s := newServer(&Config{Addr: ""})
	type args struct {
	}
	tests := []struct {
		name         string
		relativePath string
		want         gin.IRoutes
		want1        string
	}{
		{
			name:         "client",
			relativePath: "/client/hello",
			want:         s.clientGroup,
			want1:        "/hello",
		},
		{
			name:         "server",
			relativePath: "/server/world",
			want:         s.serverGroup,
			want1:        "/world",
		},
		{
			name:         "admin",
			relativePath: "/admin/wang",
			want:         s.adminGroup,
			want1:        "/wang",
		},
		{
			name:         "callback",
			relativePath: "/callback/li",
			want:         s.callbackGroup,
			want1:        "/li",
		},
		{
			name:         "other",
			relativePath: "/xx/hello",
			want:         s.engine,
			want1:        "/xx/hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := s.selectGroup(tt.relativePath)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectGroup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("selectGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
