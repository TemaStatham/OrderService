package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServConfig :
type ServConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// NatsConfig : конфигурация для nuts-streaming
type NatsConfig struct {
	URL         string `mapstructure:"url"`
	ClusterID   string `mapstructure:"cluster_id"`
	ClientID    string `mapstructure:"client_id"`
	Canal       string `mapstructure:"canal"`
	Subject     string `mapstructure:"subject"`
	QueueGroup  string `mapstructure:"queue_group"`
	DurableName string `mapstructure:"dur_name"`
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
	NatsConfig `mapstructure:"nats"`
	DBConfig   `mapstructure:"db"`
	ServConfig `mapstructure:"serv"`
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
