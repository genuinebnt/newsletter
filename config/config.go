package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database        DatabaseSettings `yaml:"database"`
	ApplicationPort int              `yaml:"application_port"`
}

type DatabaseSettings struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	DatabaseName string `yaml:"database_name"`
}

//go:embed config.yaml
var configFile []byte

func GetConfiguration() (*Config, error) {
	config := &Config{}

	err := yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
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
