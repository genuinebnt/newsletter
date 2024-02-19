package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Application struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseSettings struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	DatabaseName string `yaml:"database_name"`
}

type Config struct {
	Application Application      `yaml:"application"`
	Database    DatabaseSettings `yaml:"database"`
}

func GetConfiguration() (*Config, error) {
	appEnv := os.Getenv("APP_ENVIRONMENT")
	if appEnv != "production" && appEnv != "" {
		return nil, fmt.Errorf("wrong app environment. use either local or production")
	} else if appEnv == "" {
		appEnv = "local"
	}

	configFile := appEnv + ".yaml"

	currDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	configDir := filepath.Join(currDir, "configuration")

	v := viper.New()
	v.SetConfigFile(filepath.Join(configDir, "base.yaml"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config files: %v", err)
	}

	v.SetConfigFile(filepath.Join(configDir, configFile))
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config file: %v", err)
	}

	config := &Config{}

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return config, nil
}

func (settings DatabaseSettings) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", settings.Username, settings.Password, settings.Host, settings.Port, settings.DatabaseName,
	)
}

func (settings DatabaseSettings) ConnnectionStringWithoutDB() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d", settings.Username, settings.Password, settings.Host, settings.Port,
	)
}
