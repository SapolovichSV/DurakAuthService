package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	LogLevel int
	Server   `yaml:"server"`
	Database `yaml:"db"`
}
type Server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}
type Database struct {
	DbModel  string `yaml:"dbmodel"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbname"`
}
type logLevel struct {
	Level string `yaml:"log_level"`
}

func Build(pathToYamlConfig string) (Config, error) {
	configData, err := os.ReadFile(pathToYamlConfig)

	if err != nil {
		return Config{}, err
	}

	config := Config{}
	config.LogLevel, err = parseLogLevel(configData)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
func parseLogLevel(data []byte) (int, error) {
	var level logLevel
	if err := yaml.Unmarshal(data, &level); err != nil {
		return -1, err
	}
	if level.Level == "DEV" {
		return -4, nil
	} else if level.Level == "PROD" {
		return 4, nil
	} else {
		return -1, yaml.ErrNotFoundNode
	}
}
