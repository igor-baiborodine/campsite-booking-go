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
	filesOpt := dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))
	pathsOpt := dotenv.Paths(".")

	path := os.Getenv("ENVIRONMENT_CONFIG_PATH")
	if path != "" {
		pathsOpt = dotenv.Paths(path)
	}
	if err := dotenv.Load(filesOpt, pathsOpt); err != nil {
		return cfg, err
	}
	err := envconfig.Process("", &cfg)

	return cfg, err
}
