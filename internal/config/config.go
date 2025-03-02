package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	App struct {
		Env string `yaml:"env"`
	}

	JWT struct {
		Secret string        `yaml:"secret" env:"JWT_SECRET_KEY"`
		Exp    time.Duration `yaml:"exp"`
	}

	Http struct {
		Address string `yaml:"address"`
	}

	Storage struct {
		Sqlite struct {
			PathToDB string `yaml:"path"`
		}
	}
}

func New() (*Config, error) {

	config := &Config{}
	if err := cleanenv.ReadConfig("./config/local.yaml", config); err != nil {
		return nil, err
	}

	return config, nil
}
