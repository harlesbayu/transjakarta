package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MQTTConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RabbitMQ struct {
	URL string `mapstructure:"url"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	MQTT     MQTTConfig     `mapstructure:"mqtt"`
	RabbitMQ RabbitMQ       `mapstructure:"rabbitmq"`
}

func NewFromJSONFile(path string, filename string) Config {
	var sb strings.Builder
	sb.WriteString(path)
	sb.WriteString("/")
	sb.WriteString(filename)
	log.Println("loading config from file ... ", sb.String())

	viper.SetConfigFile(sb.String())
	viper.SetConfigType("json")
	viper.WatchConfig()

	// default value
	viper.SetDefault("Logger.Stdout", false)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return load()
}

func NewFromJSONByte(b []byte) Config {
	log.Println("loading config from byte ... ", b)
	viper.SetConfigType("json")
	if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
		panic(err)
	}

	return load()
}

func load() Config {
	cfg := Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	log.Println("config loaded!")
	return cfg
}
