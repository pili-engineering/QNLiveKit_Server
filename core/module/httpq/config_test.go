package httpq

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			name:    "success",
			addr:    "127.0.0.1:12345",
			wantErr: false,
		},
		{
			name:    "zero addr",
			addr:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Addr: tt.addr,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
