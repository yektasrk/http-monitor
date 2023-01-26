package main

import (
	"flag"

	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/httpserver"
	"github.com/yektasrk/http-monitor/pkg/logger"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./configs.yaml", "Path to configuration file.")
	flag.Parse()

	config, err := configs.Load(configFile)
	if err != nil {
		panic("failed to load configs")
	}

	err = logger.ConfigureLogger(config.Logger)
	if err != nil {
		panic("failed to configure logger")
	}

	httpMonitor, err := httpserver.New(config)
	if err != nil {
		log.Fatal("failed to initialize http monitor service", err)
	}

	if err := httpMonitor.Serve(config.HttpServer); err != nil {
		log.Fatal("failed to start http server", err)
	}
}
