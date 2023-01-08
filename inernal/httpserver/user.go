package httpserver

import (
	"net/http"

	"github.com/yektasrk/http-monitor/pkg/auth"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type UserRequest struct {
	Username string `json:"username" valid:"required, stringlength(6|12)"`
	Password string `json:"password" valid:"required, minstringlength(8)"`
}

func createUser(c echo.Context) error {
	userRequest := UserRequest{}

	if err := c.Bind(&userRequest); err != nil { //TODO: new func: json to interface
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(userRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()) //TODO: error message
	}

	data := struct {
		Username string `json:"username"`
	}{
		Username: userRequest.Username,
	}
	return c.JSON(http.StatusOK, data)
}

func loginUser(c echo.Context) error {
	userRequest := UserRequest{}

	if err := c.Bind(&userRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(userRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := auth.GenerateToken(userRequest.Username)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	return c.JSON(http.StatusBadRequest, data)
}
