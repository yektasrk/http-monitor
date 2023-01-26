package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Hash(text string) string {
	algorithm := sha256.New()
	algorithm.Write([]byte(text))
	sha := algorithm.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	return shaStr
}

func ParsRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(i); err != nil {
		log.Error("Error in validating struct", err)
		return err
	}
	return nil
}
