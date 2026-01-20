package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	ServerAddress string `mapstructure:"server_address"`
	Password      string `mapstructure:"password"`
}

type SmtpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type DBConfig struct {
	Host              string `mapstructure:"host"`
	Name              string `mapstructure:"name"`
	Password          string `mapstructure:"password"`
	Port              int    `mapstructure:"port"`
	User              string `mapstructure:"user"`
	MaxOpenConnection int    `mapstructure:"max_open_connection"`
}

type AppConfig struct {
	Redis    RedisConfig `mapstructure:"redis"`
	Smtp     SmtpConfig  `mapstructure:"smtp"`
	DBConfig DBConfig    `mapstructure:"database"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}

func New() *AppConfig {
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error reading config file: %v", err)
	}

	var appConfig AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		logrus.Error("unable to decode into struct", err)
	}
	return &appConfig
}
