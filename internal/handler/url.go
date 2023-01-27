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
	allowedIntervals   []string
	maxUrlPerUser      int
	alertsHistoryCount int
	dbClient           *db.Client
}

func NewUrlHandler(config *configs.Configuration) (*UrlHandler, error) {
	dbClient, err := db.GetDatabase(config.Postgres)
	if err != nil {
		return nil, err
	}

	log.Debug(config.UrlHandler.AllowedIntervals)
	return &UrlHandler{
		allowedIntervals:   config.UrlHandler.AllowedIntervals,
		maxUrlPerUser:      config.UrlHandler.MaxUrlPerUser,
		alertsHistoryCount: config.UrlHandler.AlertsHistoryCount,
		dbClient:           dbClient,
	}, nil
}

func (urlHandler UrlHandler) CreateUrl(owner int, address string, failureThreshold int, intervalStr string) (*db.Url, error) {
	if failureThreshold < 1 {
		log.Error(failureThreshold)
		return nil, InvalidThreshold
	}

	_, urlcount, err := urlHandler.dbClient.GetUrlsForOwner(owner)
	if err != nil {
		return nil, err
	}
	if urlcount >= int64(urlHandler.maxUrlPerUser) {
		return nil, UrlCountPerUseExceeded
	}

	if !utils.Contains(intervalStr, urlHandler.allowedIntervals) {
		return nil, InvalidInterval
	}
	interval, _ := time.ParseDuration(intervalStr)

	url := db.Url{
		OwnerID:          owner,
		Address:          address,
		FailureThreshold: failureThreshold,
		Interval:         interval,
	}
	url, err = urlHandler.dbClient.SaveUrl(url)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, UrlAlreadyExists
	}
	return &url, nil
}

func (urlHandler UrlHandler) UrlStats(urlID int, durationStr string) (int, int, int, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration = time.Duration(24 * time.Hour)
	}
	from := time.Now().Add(-1 * duration)
	to := time.Now()
	requests, reqCount, err := urlHandler.dbClient.GetRequestsForUrl(urlID, from, to)
	if err != nil {
		return 0, 0, 0, err
	}

	successRequests := 0
	failedRequests := 0
	for _, request := range requests {
		if request.StatusCode/100 == 2 {
			successRequests += 1
		} else {
			failedRequests += 1
		}
	}
	return successRequests, failedRequests, int(reqCount), err
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

func (urlHandler UrlHandler) GetAlerts(urlID int) (string, []map[string]interface{}, error) {
	alerts, err := urlHandler.dbClient.GetLatestAlertsForUrl(urlID, urlHandler.alertsHistoryCount)
	if err != nil {
		return "Unknown", nil, err
	}

	var alertsMap []map[string]interface{}
	for _, alert := range alerts {
		urlMap, err := utils.StructToMap(alert, []string{"ID", "TimeFired", "TimeResolved"})
		if err != nil {
			return "Unknown", nil, err
		}
		alertsMap = append(alertsMap, urlMap)
	}

	state := "OK!"
	if len(alerts) > 0 && alerts[0].TimeResolved == nil {
		log.Debug(alerts[0])
		state = "*** FIRING ***"
	}
	return state, alertsMap, err
}
