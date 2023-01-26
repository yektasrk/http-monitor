package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yektasrk/http-monitor/internal/handler"
	"github.com/yektasrk/http-monitor/pkg/auth"
)

var userHandler handler.UserHandler

func InitAuth(handler handler.UserHandler) {
	userHandler = handler
}

func LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) < 2 || authHeaderParts[0] != "Bearer" {
			return echo.ErrBadRequest
		}

		jwtToken := authHeaderParts[1]
		userID, err := auth.ValidateToken(userHandler.JwtSecretKey, jwtToken)
		if err != nil {
			return echo.ErrBadRequest
		}

		c.Set("userID", userID)
		return next(c)
	}
}
