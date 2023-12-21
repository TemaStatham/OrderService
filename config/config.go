package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NatsConfig : конфигурация для nuts-streaming
type NatsConfig struct {
	URL       string `mapstructure:"url"`
	ClusterID string `mapstructure:"cluster_id"`
	ClientID  string `mapstructure:"client_id"`
	Subject   string `mapstructure:"subject"`
}

// DBConfig :
type DBConfig struct {
	DBHost     string `mapstructure:"dbhost"`
	DBPort     string `mapstructure:"dbport"`
	DBUser     string `mapstructure:"dbuser"`
	DBPassword string `mapstructure:"dbpassword"`
	DBName     string `mapstructure:"dbname"`
	DBSSLMode  string `mapstructure:"dbsslmode"`
}

// Config : структура файла конфигурации
type Config struct {
	NATS NatsConfig `mapstructure:"nats"`
	DB   DBConfig   `mapstructure:"db"`
}

// Load : загружает файл конфигурации
func Load(configType, configName, configPath string) (*Config, error) {
	var config *Config

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора конфигурации: %w", err)
	}

	return config, nil
}
