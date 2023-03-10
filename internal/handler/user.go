package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/db"
	"github.com/yektasrk/http-monitor/pkg/auth"
	"github.com/yektasrk/http-monitor/pkg/utils"
)

var (
	UserNotFoundError  = errors.New("User Not Found")
	UserAlreadyExists  = errors.New("User Already Exists")
	PasswordNotCorrect = errors.New("Password Not Correct")
)

type UserHandler struct {
	JwtTokenExpireTime time.Duration
	JwtSecretKey       []byte
	dbClient           *db.Client
}

func NewUserHandler(config *configs.Configuration) (*UserHandler, error) {
	jwtTokenExpireTime, err := time.ParseDuration(config.JwtAuth.ExpireTime)
	jwtKeybytes := []byte(config.JwtAuth.SecretKey)

	dbClient, err := db.GetDatabase(config.Postgres)
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		JwtTokenExpireTime: jwtTokenExpireTime,
		JwtSecretKey:       jwtKeybytes,
		dbClient:           dbClient,
	}, nil
}

func (userHandler UserHandler) CreateUser(username, password string) error {
	protectedPassword := utils.Hash(password)
	user := db.User{
		Username: username,
		Password: protectedPassword,
	}
	err := userHandler.dbClient.SaveUser(user)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return UserAlreadyExists
	}
	return err
}

func (userHandler UserHandler) AuthUser(username, password string) (string, error) {
	user, err := userHandler.dbClient.GetUser(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", UserNotFoundError
	} else if err != nil {
		return "", err
	}

	protectedPassword := utils.Hash(password)
	if user.Password != protectedPassword {
		return "", PasswordNotCorrect
	}

	return auth.GenerateToken(userHandler.JwtSecretKey, userHandler.JwtTokenExpireTime, user.ID)
}
