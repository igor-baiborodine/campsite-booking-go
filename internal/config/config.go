package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	PGConfig struct {
		Conn string `envconfig:"PG_CONN" default:"host=postgres dbname=${CAMPGROUNDS_DB} user=${CAMPGROUNDS_USER} password=${CAMPGROUNDS_PASSWORD}"`
	}

	RPCConfig struct {
		Host string `default:"0.0.0.0"`
		Port string `default:":8085"`
	}

	AppConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL"        default:"DEBUG"`
		PG              PGConfig
		RPC             RPCConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func (c RPCConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func InitConfig() (AppConfig, error) {
	cfg := AppConfig{}
	if err := dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return cfg, err
	}
	err := envconfig.Process("", &cfg)

	return cfg, err
}

// ReplaceEnvPlaceholders replaces placeholders of the format ${VAR_NAME}
// with their environment variable values only if the variable is set.
func ReplaceEnvPlaceholders(configValue string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	result := re.ReplaceAllStringFunc(configValue, func(placeholder string) string {
		envVar := placeholder[2 : len(placeholder)-1]
		val, ok := os.LookupEnv(envVar)
		if ok {
			return val
		}
		return placeholder
	})
	return result
}
