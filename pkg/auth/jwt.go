package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
		return -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"].(string) == iss || claims["exp"].(int64) <= time.Now().Unix() {
			id := int(claims["id"].(float64))
			return id, nil
		}
		return -1, tokenNotValidErr
	}
	return -1, err
}
