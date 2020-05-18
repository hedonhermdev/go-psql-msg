package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Host string
	Port int
}

type Config struct {
	Database DBConfig    `yaml:"database"`
	Kafka    KafkaConfig `yaml:"kafka"`
	Channels []string    `yaml:"channels"`
}

func ParseConfig(r io.Reader) (*Config, error) {
	yaml_bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)

	if c.Database.Host == "" {
		c.Database.Host = "localhost"
	}

	if c.Database.Port == 0 {
		c.Database.Port = 5432
	}

	yaml.Unmarshal(yaml_bytes, c)

	return c, nil
}
