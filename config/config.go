package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type PostgresConfig struct {
	Name string

	User string
	Pass string

	Host string
	Port string
}

func (c PostgresConfig) ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Pass, c.Host, c.Port, c.Name,
	)
}

func (c PostgresConfig) DataSourceString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Name, c.Pass, "disable")
}

type AppConfig struct {
	PostgresConfig `yaml:"-"`

	HttpConfig struct {
		Port string `yaml:"port"`
	} `yaml:"http"`
}

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	var config AppConfig

	bconfig, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		logrus.Fatalf("cant read file ./config/config.yaml cause: %s", err.Error())
	}

	err = yaml.Unmarshal(bconfig, &config)
	if err != nil {
		logrus.Fatalf("cant unmarshal file, cause: %s", err.Error())
	}

	config.PostgresConfig = PostgresConfig{
		Name: os.Getenv("POSTGRES_NAME"),
		User: os.Getenv("POSTGRES_USER"),
		Pass: os.Getenv("POSTGRES_PASS"),
		Host: os.Getenv("POSTGRES_HOST"),
		Port: os.Getenv("POSTGRES_PORT"),
	}

	return config
}
