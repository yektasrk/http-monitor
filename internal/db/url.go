package db

import "time"

type Url struct {
	ID               uint
	Address          string `gorm:"unique"`
	OwnerID          int
	Owner            User `gorm:"foreignKey:ownerID"`
	FailureThreshold int
	Interval         time.Duration
}

type Request struct {
	ID     uint
	url    Url
	time   time.Time
	result string
}

func (client Client) SaveUrl(url Url) error {
	result := client.db.Create(&url)
	return result.Error
}

func (client Client) GetUrlsForOwner(ownerID int) ([]Url, int64, error) {
	urls := []Url{}
	result := client.db.Where("owner_id = ?", ownerID).Find(&urls)
	return urls, result.RowsAffected, result.Error
}

func (client Client) GetRequestsForUrl(urlID int) (map[string]interface{}, error) {
	url := make(map[string]interface{})
	result := client.db.Where("url_id = ?", urlID).Find(&url)
	return url, result.Error
}
