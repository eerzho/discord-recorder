package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Discord Discord
	Logger  Logger
	Minio   Minio
}

type Discord struct {
	Token       string `env:"DISCORD_TOKEN,required"`
	RecordLimit int    `env:"DISCORD_RECORD_LIMIT,default=5"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL,default=info"`
}

type Minio struct {
	Endpoint  string `env:"MINIO_ENDPOINT,required"`
	User      string `env:"MINIO_USER,required"`
	Password  string `env:"MINIO_PASSWORD,required"`
	Bucket    string `env:"MINIO_BUCKET,required"`
	AccessKey string `env:"MINIO_ACCESS_KEY,required"`
	SecretKey string `env:"MINIO_SECRET_KEY,required"`
}

var cfg Config

func Parse() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return err
	}

	return nil
}

func Cfg() *Config {
	return &cfg
}
