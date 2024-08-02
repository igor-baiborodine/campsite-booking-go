package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_InitConfig(t *testing.T) {
	tests := map[string]struct {
		env        string
		configPath string
		want       *AppConfig
		wantErr    error
	}{
		"Success": {
			env:        "test",
			configPath: "./../../",
			want: &AppConfig{
				Environment: "test",
				LogLevel:    "WARN",
				PG: PGConfig{
					Conn: "host=localhost dbname=test_campgrounds user=test_campgrounds_user password=test_campgrounds_pass",
				},
				RPC:             RPCConfig{Host: "0.0.0.0", Port: ":8085"},
				ShutdownTimeout: time.Duration(30) * time.Second,
			},
			wantErr: nil,
		},
	}
	defer os.Setenv("ENVIRONMENT", "")
	defer os.Setenv("ENVIRONMENT_CONFIG_PATH", "")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			os.Setenv("ENVIRONMENT", tc.env)
			os.Setenv("ENVIRONMENT_CONFIG_PATH", tc.configPath)
			// when
			got, err := InitConfig()
			// then
			if tc.wantErr != nil {
				assert.Equal(
					t,
					tc.wantErr,
					err,
					"InitConfig() error = %v, wantErr %v",
					err,
					tc.wantErr,
				)
				return
			}
			assert.Equal(t, tc.want, &got, "InitConfig() got = %v, want %v", err, tc.wantErr)
		})
	}
}
