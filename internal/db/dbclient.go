package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/yektasrk/http-monitor/configs"
)

type Client struct {
	db *gorm.DB
}

func GetDatabase(config configs.PostgresConfiguration) (*Client, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Database,
		config.Password,
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return &Client{
		db: db,
	}, nil
}
