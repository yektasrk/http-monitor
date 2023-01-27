package db

import (
	"time"
)

type Url struct {
	ID               uint
	Address          string `gorm:"unique"`
	OwnerID          int
	Owner            User `gorm:"foreignKey:ownerID"`
	FailureThreshold int
	Interval         time.Duration
}

func (client Client) SaveUrl(url Url) (Url, error) {
	result := client.db.Create(&url)
	return url, result.Error
}

func (client Client) GetUrlsForOwner(ownerID int) ([]Url, int64, error) {
	urls := []Url{}
	result := client.db.Where("owner_id = ?", ownerID).Find(&urls)
	return urls, result.RowsAffected, result.Error
}

func (client Client) GetUrls() ([]Url, int64, error) {
	urls := []Url{}
	result := client.db.Find(&urls)
	return urls, result.RowsAffected, result.Error
}
