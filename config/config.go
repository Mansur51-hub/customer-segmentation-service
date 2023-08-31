package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Port  string `yaml:"port" mapstructure:"http_port"`
	Url   string `yaml:"url" mapstructure:"pg_url"`
	Level string `yaml:"level" mapstructure:"log_level"`
}

func NewConfig() (cfg *Config, err error) {
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")

	cfg = &Config{}

	var result map[string]interface{}

	err = viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	err = viper.Unmarshal(&result)

	err = mapstructure.Decode(result, &cfg)

	if err != nil {
		return nil, fmt.Errorf("error marshalling config: %w", err)
	}

	_ = viper.BindEnv("isLocal", "IS_DOCKER_RUN")

	if !viper.GetBool("isLocal") {
		cfg.Url = viper.GetString("pg_local_url")
	}

	return cfg, nil
}
