package httpserver

import (
	"fmt"
	"strconv"

	"github.com/yektasrk/http-monitor/configs"

	"github.com/labstack/echo/v4"
)

const apiPrefix = "api/v1/"

type httpMonitor struct {
	httpMonitorHandler httpMonitorHandler
}

func New(config *configs.Configuration) (*httpMonitor, error) {
	httpMonitorHandler, err := NewHttpMonitorHandler(config)
	if err != nil {
		return nil, err
	}

	return &httpMonitor{
		httpMonitorHandler: *httpMonitorHandler,
	}, nil
}

func (httpMonitor httpMonitor) Serve(config configs.HttpServerConfiguration) error {
	e := echo.New()
	e.POST(apiPrefix+"users/", httpMonitor.httpMonitorHandler.createUser)
	e.POST(apiPrefix+"users/login/", httpMonitor.httpMonitorHandler.loginUser)

	address := config.Host + ":" + strconv.Itoa(config.Port)
	err := e.Start(address)
	if err != nil {
		fmt.Errorf("failed to start http listener: ", err)
		return err
	}
	return nil
}
