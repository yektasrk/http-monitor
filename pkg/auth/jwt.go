package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var (
	iss              = "http-monitor"
	tokenNotValidErr = errors.New("token is not valid")
)

func GenerateToken(secretKey []byte, expireTime time.Duration, ID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  ID,
		"iss": iss,
		"exp": time.Now().Unix() + int64(expireTime.Seconds()),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Error("error in signing token: ", err)
	}
	return tokenString, err
}

func ValidateToken(secretKey []byte, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		log.Error("error in validating token: ", err)
		return -1, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return -1, tokenNotValidErr
	}
	if claims["iss"].(string) != iss || int64(claims["exp"].(float64)) < time.Now().Unix() {
		return -1, tokenNotValidErr
	}

	id := int(claims["id"].(float64))
	return id, nil
}
