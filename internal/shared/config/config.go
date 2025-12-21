package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type MinIO struct {
	Endpoint  string `yaml:"endpoint" env-required:"true"`
	Bucket    string `yaml:"bucket" env-required:"true"`
	UseSSL    bool   `yaml:"use_ssl" env-required:"true"`
	AccessKey string `env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey string `env:"MINIO_SECRET_KEY" env-required:"true"`
}

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	MinIO      `yaml:"minio"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Config file does not exist: %s", configPath)
		}
		log.Fatalf("Failed to stat config file: %v", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}
