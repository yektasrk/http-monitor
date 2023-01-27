package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"

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
		log.Error("Error in validating struct: ", err)
		return err
	}
	return nil
}

func StructToMap(s interface{}, exportedFields []string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Unexpexted type")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if Contains(fi.Name, exportedFields) {
			out[fi.Name] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func Contains(obj string, list []string) bool {
	for _, e := range list {
		if obj == e {
			return true
		}
	}
	return false
}
