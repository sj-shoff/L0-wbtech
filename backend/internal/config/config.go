package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string       `yaml:"env" env-default:"local"`
	Server     ServerConfig `yaml:"server"`
	Database   Postgres     `yaml:"postgres"`
	Kafka      KafkaConfig  `yaml:"kafka"`
	Migrations string       `yaml:"migrations" env-default:"./migrations"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		slog.Error("Config path is required")
		os.Exit(1)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist", "path", configPath)
		os.Exit(1)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Failed to read config", "error", err)
		os.Exit(1)
	}

	loadSecrets(&cfg)

	return &cfg
}

func fetchConfigPath() string {
	flag.String("config", "", "path to config file")
	flag.Parse()

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return "config.yaml"
}

func loadSecrets(cfg *Config) {
	if cfg.Database.Password == "" {
		if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
			cfg.Database.Password = password
		} else {
			slog.Error("POSTGRES_PASSWORD environment variable not set")
			os.Exit(1)
		}
	}
}
