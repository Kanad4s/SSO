package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Usage:
// Flag - --config=./path/to/config/file.yaml
// ENV - CONFIG_PATH=./path/to/config/file.yaml

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("empty config path")
	}
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// Priority: flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
