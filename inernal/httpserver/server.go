package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Serve() error {
	e := echo.New()

	err := e.Start(":8000")
	if err != nil {
		fmt.Errorf("failed to start http listener: ", err)
		return err
	}
	return nil
}
