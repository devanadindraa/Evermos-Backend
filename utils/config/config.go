package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host        string      `envconfig:"host"`
	Port        int         `envconfig:"port" validate:"number,required"`
	Environment Environment `envconfig:"environment" validate:"oneof=DEVELOPMENT TEST STAGING PRODUCTION"`
	Version     string      `envconfig:"version" default:"development"`
	Database    Database    `envconfig:"database"`
}

type Database struct {
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Host     string `envconfig:"host"`
	Port     string `envconfig:"port"`
	Name     string `envconfig:"name"`
}

var config *Config

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env: %v\n", err)
	}

	var c Config
	err = envconfig.Process("Backend", &c)
	if err != nil {
		panic("Failed to Process env : " + err.Error())
	}

	config = &c

	return config
}

func GetConfig() *Config {
	if config != nil {
		return config
	}
	return NewConfig()
}
