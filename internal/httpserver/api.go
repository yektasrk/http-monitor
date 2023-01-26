package httpserver

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/handler"

	"github.com/yektasrk/http-monitor/pkg/utils"

	"github.com/labstack/echo/v4"
)

type httpMonitorHandler struct {
	userHandler *handler.UserHandler
}

type UserRequest struct {
	Username string `json:"username" valid:"required, stringlength(6|12)"`
	Password string `json:"password" valid:"required, minstringlength(8)"`
}

func NewHttpMonitorHandler(config *configs.Configuration) (*httpMonitorHandler, error) {
	userHandler, err := handler.NewUserHandler(config)
	if err != nil {
		return nil, err
	}

	return &httpMonitorHandler{
		userHandler: userHandler,
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
		log.Error("Error in creating user", err)
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
		log.Error("Error in authenticating user", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	return c.JSON(http.StatusOK, data)
}

func (httpMonitor httpMonitorHandler) ListUrls(c echo.Context) error {
	userID := c.Get("userID")
	fmt.Println(userID)
	urls := []string{"1", "2"}

	data := struct {
		Urls []string `json:"urls"`
	}{
		Urls: urls,
	}
	return c.JSON(http.StatusOK, data)
}
