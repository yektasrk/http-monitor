package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/db"
	"github.com/yektasrk/http-monitor/internal/handler"

	"github.com/yektasrk/http-monitor/pkg/utils"

	"github.com/labstack/echo/v4"
)

type httpMonitorHandler struct {
	urlsToSchedule chan db.Url
	userHandler    *handler.UserHandler
	urlHandler     *handler.UrlHandler
}

type UserRequest struct {
	Username string `json:"username" valid:"required, stringlength(6|12)"`
	Password string `json:"password" valid:"required, minstringlength(8)"`
}

type UrlRequest struct {
	Address          string `json:"address" valid:"required, url"`
	FailureThreshold string `json:"failureThreshold" valid:"required, numeric"`
	Interval         string `json:"interval" valid:"required"`
}

func NewHttpMonitorHandler(config *configs.Configuration, urlsToSchedule chan db.Url) (*httpMonitorHandler, error) {
	userHandler, err := handler.NewUserHandler(config)
	if err != nil {
		return nil, err
	}

	urlHandler, err := handler.NewUrlHandler(config)
	if err != nil {
		return nil, err
	}

	return &httpMonitorHandler{
		urlsToSchedule: urlsToSchedule,
		userHandler:    userHandler,
		urlHandler:     urlHandler,
	}, nil
}

func (httpMonitor httpMonitorHandler) createUser(c echo.Context) error {
	userRequest := UserRequest{}
	err := utils.ParsRequest(c, &userRequest)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = httpMonitor.userHandler.CreateUser(userRequest.Username, userRequest.Password)
	if errors.Is(err, handler.UserAlreadyExists) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		log.Error("Error in creating user: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		Username string `json:"username"`
	}{
		Username: userRequest.Username,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) loginUser(c echo.Context) error {
	userRequest := UserRequest{}
	err := utils.ParsRequest(c, &userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := httpMonitor.userHandler.AuthUser(userRequest.Username, userRequest.Password)
	if errors.Is(err, handler.UserNotFoundError) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	} else if errors.Is(err, handler.PasswordNotCorrect) {
		return echo.NewHTTPError(http.StatusOK, err.Error())
	} else if err != nil {
		log.Error("Error in authenticating user: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) CreateUrl(c echo.Context) error {
	userID := c.Get("userID").(int)
	urlRequest := UrlRequest{}
	err := utils.ParsRequest(c, &urlRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	failureThresholdInt, _ := strconv.Atoi(urlRequest.FailureThreshold)
	url, err := httpMonitor.urlHandler.CreateUrl(userID, urlRequest.Address, failureThresholdInt, urlRequest.Interval)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		//TODO
	}
	httpMonitor.urlsToSchedule <- *url

	data := struct {
		ID  uint   `json:"id"`
		Url string `json:"url"` //TODO
	}{
		ID:  url.ID,
		Url: url.Address,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) GetUrlStats(c echo.Context) error {
	urlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
	}

	duration := c.QueryParam("duration")
	successRequests, failedRequests, allRequests, err := httpMonitor.urlHandler.UrlStats(urlID, duration)
	if err != nil {
		log.Error("Error in getting url stats: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		SuccessRequests int `json:"successRequests"`
		FailedRequests  int `json:"failedRequests"`
		AllRequests     int `json:"allRequests"`
	}{
		SuccessRequests: successRequests,
		FailedRequests:  failedRequests,
		AllRequests:     allRequests,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) ListUrls(c echo.Context) error {
	userID := c.Get("userID").(int)
	urls, count, err := httpMonitor.urlHandler.ListUserUrls(userID)
	if err != nil {
		log.Error("Error in list urls: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		Urls  []map[string]interface{} `json:"urls"`
		Count int                      `json:"count"`
	}{
		Urls:  urls,
		Count: count,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) GetUrlAlerts(c echo.Context) error {
	urlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
	}

	state, latestAlerts, err := httpMonitor.urlHandler.GetAlerts(urlID)
	if err != nil {
		log.Error("Error in get alerts: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		CurrentState string                   `json:"currentState"`
		LatestAlerts []map[string]interface{} `json:"latestAlerts"`
	}{
		CurrentState: state,
		LatestAlerts: latestAlerts,
	}
	return c.JSON(http.StatusOK, data)
}
