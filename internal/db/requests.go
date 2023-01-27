package db

import "time"

type Request struct {
	ID         uint
	UrlID      int
	Url        Url
	Time       time.Time
	StatusCode int
}

func (client Client) GetRequestsForUrl(urlID int, from, to time.Time) ([]Request, int64, error) {
	requests := []Request{}
	result := client.db.Where("url_id = ? AND time BETWEEN ? AND ?", urlID, from, to).Find(&requests)
	return requests, result.RowsAffected, result.Error
}

func (client Client) SaveRequest(request Request) error {
	result := client.db.Create(&request)
	return result.Error
}
