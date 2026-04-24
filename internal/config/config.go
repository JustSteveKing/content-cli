package config

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type DefaultsConfig struct {
	Author string `mapstructure:"author" yaml:"author"`
	Draft  bool   `mapstructure:"draft" yaml:"draft"`
}

type CollectionConfig struct {
	Dir            string         `mapstructure:"dir" yaml:"dir"`
	Format         string         `mapstructure:"format" yaml:"format"`
	Template       string         `mapstructure:"template" yaml:"template,omitempty"`
	RequiredFields []string       `mapstructure:"required_fields" yaml:"required_fields"`
	OptionalFields []string       `mapstructure:"optional_fields" yaml:"optional_fields,omitempty"`
	Slug           string         `mapstructure:"slug" yaml:"slug"`
	Defaults       DefaultsConfig `mapstructure:"defaults" yaml:"defaults"`
}

type Config struct {
	DefaultCollection string                      `mapstructure:"default_collection" yaml:"default_collection,omitempty"`
	Collections       map[string]CollectionConfig `mapstructure:"collections" yaml:"collections"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".content.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Exists() bool {
	_, err := os.Stat(".content.yaml")
	return err == nil
}

func Write(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(".content.yaml", data, 0644)
}
