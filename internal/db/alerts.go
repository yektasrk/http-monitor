package db

import "time"

type Alert struct {
	ID           uint
	UrlID        int
	Url          Url
	TimeFired    time.Time
	TimeResolved *time.Time
}

func (client Client) SaveAlert(alert Alert) error {
	result := client.db.Create(&alert)
	return result.Error
}

func (client Client) GetLastAlert() (Alert, error) { //TODO
	alert := Alert{}
	result := client.db.Last(&alert)
	return alert, result.Error
}

func (client Client) UpdateAlert(alert Alert, fields Alert) error {
	result := client.db.Model(&alert).Update(&fields)
	return result.Error
}

func (client Client) GetLatestAlertsForUrl(urlID int, limit int) ([]Alert, error) {
	alerts := []Alert{}
	result := client.db.Order("id desc").Where("url_id = ?", urlID).Find(&alerts).Limit(limit)
	return alerts, result.Error
}
