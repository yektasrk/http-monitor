package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const apiPrefix = "api/v1/"

func Serve() error {
	e := echo.New()

	e.POST(apiPrefix+"users/", createUser)
	e.POST(apiPrefix+"users/login/", loginUser)

	err := e.Start(":8000") // TODO: config
	if err != nil {
		fmt.Errorf("failed to start http listener: ", err)
		return err
	}
	return nil
}
