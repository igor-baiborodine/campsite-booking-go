package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	PGConfig struct {
		Conn string `required:"true"`
	}

	RpcConfig struct {
		Host string `default:"0.0.0.0"`
		Port string `default:":8085"`
	}

	AppConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG              PGConfig
		Rpc             RpcConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func (c RpcConfig) Address() string {
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
