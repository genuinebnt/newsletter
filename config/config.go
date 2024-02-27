package config

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Application struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseSettings struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

type Config struct {
	Application Application      `yaml:"application"`
	Database    DatabaseSettings `yaml:"database"`
}

func GetConfiguration() (*Config, error) {
	appEnv := os.Getenv("APP_ENVIRONMENT")
	if appEnv == "" {
		appEnv = "local"
	} else if appEnv != "production" && appEnv != "local" {
		return nil, fmt.Errorf("wrong App environment. Expected 'production' or 'local'")
	}

	configFile := appEnv + ".yaml"

	currDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	configDir := filepath.Join(currDir, "configuration")

	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetConfigFile(filepath.Join(configDir, "base.yaml"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config files: %v", err)
	}

	v.SetConfigFile(filepath.Join(configDir, configFile))

	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config file: %v", err)
	}

	fmt.Println("Merged configuration:", v.AllSettings())

	config := &Config{}

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	log.Println(config.Application.Host, config.Application.Port, config.Database.Name, config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password)

	return config, nil
}

func (settings DatabaseSettings) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", settings.Username, settings.Password, settings.Host, settings.Port, settings.Name,
	)
}

func (settings DatabaseSettings) ConnnectionStringWithoutDB() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d", settings.Username, settings.Password, settings.Host, settings.Port,
	)
}
