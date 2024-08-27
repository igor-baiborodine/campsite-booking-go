package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	// given
	os.Setenv("LOG_LEVEL", "INFO")
	os.Setenv("SHUTDOWN_TIMEOUT", "15s")
	// when
	cfg, err := InitConfig()
	// then
	assert.NoError(t, err)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, 15*time.Second, cfg.ShutdownTimeout)
}

func TestReplaceEnvPlaceholders(t *testing.T) {
	os.Setenv("DB_NAME", "db_name")
	os.Setenv("DB_USER", "db_user")
	os.Setenv("DB_PASSWORD", "db_password")

	tests := map[string]struct {
		configValue string
		want        string
	}{
		"ConfigValueWithPlaceholders_EnvVarsPresent_ReplacedWithEnvVars": {
			configValue: "host=postgres dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD}",
			want:        "host=postgres dbname=db_name user=db_user password=db_password",
		},
		"ConfigValueWithoutPlaceholders_NotChanged": {
			configValue: "host=postgres dbname=db_name user=db_user password=db_password",
			want:        "host=postgres dbname=db_name user=db_user password=db_password",
		},
		"ConfigValueWithPlaceholders_EnvVarsNotPresent_NotChanged": {
			configValue: "host=postgres dbname=${DB_NAME_2} user=${DB_USER_2} password=${DB_PASSWORD_2}",
			want:        "host=postgres dbname=${DB_NAME_2} user=${DB_USER_2} password=${DB_PASSWORD_2}",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			got := ReplaceEnvPlaceholders(tc.configValue)
			// then
			assert.Equal(t, tc.want, got,
				"ReplaceEnvPlaceholders() got = %v, want %v", got, tc.want)
		})
	}
}
