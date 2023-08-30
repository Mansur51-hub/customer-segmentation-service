package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Net struct {
	Port string `yaml:"port"`
}

type PG struct {
	Url string `yaml:"url"`
}

type Log struct {
	Level string `yaml:"level"`
}

type Config struct {
	Net `yaml:"NET"`
	PG  `yaml:"PG"`
	Log `yaml:"LOG"`
}

func NewConfig(path string) (cfg *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("LOG.level", "debug")

	cfg = &Config{}

	err = viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, fmt.Errorf("error marshalling config: %w", err)
	}

	return cfg, nil
}
