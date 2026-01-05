package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type MinIO struct {
	Endpoint  string `yaml:"endpoint" env-required:"true"`
	Bucket    string `yaml:"bucket" env-required:"true"`
	UseSSL    bool   `yaml:"use_ssl"`
	AccessKey string `env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey string `env:"MINIO_SECRET_KEY" env-required:"true"`
}

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Redis struct {
	Addr     string `yaml:"address" env-required:"true"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env-default:"0"`
}

type RefreshCookie struct {
	Name   string `yaml:"name" env-required:"true"`
	Path   string `yaml:"path" env-required:"true"`
	Expiry string `yaml:"expiry" env-required:"true"`
	Secure bool   `yaml:"secure"`
}

type SessionCookie struct {
	Name      string `yaml:"name" env-required:"true"`
	SecretKey string `env:"SESSION_SECRET" env-required:"true"`
	Path      string `yaml:"path" env-required:"true"`
	Expiry    string `yaml:"expiry" env-required:"true"`
	Secure    bool   `yaml:"secure"`
}

type Cookies struct {
	Refresh RefreshCookie `yaml:"refresh" env-required:"true"`
	Session SessionCookie `yaml:"session" env-required:"true"`
}

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	SqliteDbPath string `yaml:"db_path" env-required:"true"`
	DbSource     string `env:"POSTGRESQL_DB_SOURCE" env-required:"true"`
	MinIO        `yaml:"minio"`
	Redis        `yaml:"redis"`
	Cookies      `yaml:"cookies" env-required:"true"`
	HTTPServer   `yaml:"http_server"`
}

func MustLoad() *Config {
	// DEBUG: Print where the code is actually running
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s", wd)

	// Try loading explicitly to see the specific error
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err) // logs specific error (e.g. permission denied)
	}

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
