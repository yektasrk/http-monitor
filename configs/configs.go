package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	HttpServer HttpServerConfiguration
	Postgres   PostgresConfiguration
	Logger     LoggerConfiguration
	JwtAuth    JwtAuth
	UrlHandler UrlHandlerConfiguration
}

type HttpServerConfiguration struct {
	Host string
	Port int
}

type PostgresConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type JwtAuth struct {
	SecretKey  string
	ExpireTime string
}

type LoggerConfiguration struct {
	Level string
}

type UrlHandlerConfiguration struct {
	MaxUrlPerUser      int
	AlertsHistoryCount int
	AllowedIntervals   []string
}

func Load(filename string) (*Configuration, error) {
	var config Configuration
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
