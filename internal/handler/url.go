package handler

import (
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/db"
	"github.com/yektasrk/http-monitor/pkg/utils"
)

var (
	InvalidThreshold       = errors.New("Invalid Failure threshold")
	UrlAlreadyExists       = errors.New("URL Already Exists")
	UrlCountPerUseExceeded = errors.New("Max Number of URLs Reached")
	InvalidInterval        = errors.New("Invalid Interval")
)

type UrlHandler struct {
	allowedIntervals []string
	maxUrlPerUser    int
	dbClient         *db.Client
}

func NewUrlHandler(config *configs.Configuration) (*UrlHandler, error) {
	dbClient, err := db.GetDatabase(config.Postgres)
	if err != nil {
		return nil, err
	}

	log.Debug(config.UrlHandler.AllowedIntervals)
	return &UrlHandler{
		allowedIntervals: config.UrlHandler.AllowedIntervals,
		maxUrlPerUser:    config.UrlHandler.MaxUrlPerUser,
		dbClient:         dbClient,
	}, nil
}

func (urlHandler UrlHandler) CreateUrl(owner int, address string, failureThreshold int, intervalStr string) error {
	if failureThreshold < 1 {
		log.Error(failureThreshold)
		return InvalidThreshold
	}

	_, urlcount, err := urlHandler.dbClient.GetUrlsForOwner(owner)
	if err != nil {
		return err
	}
	if urlcount >= int64(urlHandler.maxUrlPerUser) {
		return UrlCountPerUseExceeded
	}

	if !utils.Contains(intervalStr, urlHandler.allowedIntervals) {
		return InvalidInterval
	}
	interval, _ := time.ParseDuration(intervalStr)

	url := db.Url{
		OwnerID:          owner,
		Address:          address,
		FailureThreshold: failureThreshold,
		Interval:         interval,
	}
	err = urlHandler.dbClient.SaveUrl(url)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return UrlAlreadyExists
	}
	return err
}

func (urlHandler UrlHandler) UrlStats(urlID int) (map[string]interface{}, error) {
	url, err := urlHandler.dbClient.GetRequestsForUrl(urlID)
	if err != nil {
		return nil, err
	}
	log.Println(url)
	return url, err
}

func (urlHandler UrlHandler) ListUserUrls(ownerID int) ([]map[string]interface{}, int, error) {
	urls, urlcount, err := urlHandler.dbClient.GetUrlsForOwner(ownerID)
	if err != nil {
		return nil, 0, err
	}

	var urlsMap []map[string]interface{}
	for _, url := range urls {
		urlMap, err := utils.StructToMap(url, []string{"ID", "Address", "FailureThreshold", "Interval"})
		urlMap["Interval"] = urlMap["Interval"].(time.Duration).String()
		if err != nil {
			return nil, 0, err
		}
		urlsMap = append(urlsMap, urlMap)
	}
	return urlsMap, int(urlcount), err
}
