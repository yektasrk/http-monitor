package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	HttpServer HttpServerConfiguration
}

type HttpServerConfiguration struct {
	Host string
	Port int
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
