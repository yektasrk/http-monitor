package httpserver

import (
	"strconv"

	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/db"
	"github.com/yektasrk/http-monitor/internal/middleware"

	"github.com/labstack/echo/v4"
)

const apiPrefix = "api/v1/"

type httpMonitor struct {
	httpMonitorHandler httpMonitorHandler
}

func New(config *configs.Configuration, urlsToSchedule chan db.Url) (*httpMonitor, error) {
	httpMonitorHandler, err := NewHttpMonitorHandler(config, urlsToSchedule)
	if err != nil {
		return nil, err
	}

	middleware.InitAuth(*httpMonitorHandler.userHandler)

	return &httpMonitor{
		httpMonitorHandler: *httpMonitorHandler,
	}, nil
}

func (httpMonitor httpMonitor) Serve(config configs.HttpServerConfiguration) error {
	e := echo.New()
	e.POST(apiPrefix+"users/", httpMonitor.httpMonitorHandler.createUser)
	e.POST(apiPrefix+"users/login/", httpMonitor.httpMonitorHandler.loginUser)

	urlgroup := e.Group(apiPrefix + "urls")
	urlgroup.Use(middleware.LoginRequired)
	urlgroup.GET("/", httpMonitor.httpMonitorHandler.ListUrls)
	urlgroup.GET("/:id", httpMonitor.httpMonitorHandler.GetUrlStats)
	urlgroup.POST("/", httpMonitor.httpMonitorHandler.CreateUrl)
	urlgroup.GET("/:id/alerts", httpMonitor.httpMonitorHandler.GetUrlAlerts)

	address := config.Host + ":" + strconv.Itoa(config.Port)
	err := e.Start(address)
	if err != nil {
		return err
	}
	return nil
}
