package main

import (
	"flag"
	"fmt"

	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/httpserver"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./configs.yaml", "Path to configuration file.")
	flag.Parse()

	config, err := configs.Load(configFile)
	if err != nil {
		panic("failed to load configs")
	}

	httpMonitor, err := httpserver.New(config)
	if err != nil {
		fmt.Print(err)
		panic("failed to initialize http monitor service")
	}

	if err := httpMonitor.Serve(config.HttpServer); err != nil {
		panic("failed to start http server")
	}
}
