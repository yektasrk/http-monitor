package monitor

import (
	"errors"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yektasrk/http-monitor/configs"
	"github.com/yektasrk/http-monitor/internal/db"

	"github.com/go-co-op/gocron"
)

const NoResponseCode = -1

var (
	UrlAlreadyScheduled = errors.New("Url Already Scheduled")
)

type MonitorJob struct {
	syncInterval      time.Duration
	UrlsToSchedule    chan db.Url
	dbClient          *db.Client
	scheduledUrls     map[uint]*gocron.Job
	scheduler         *gocron.Scheduler
	failesCountPerUrl map[uint]int
}

func New(config *configs.Configuration) (*MonitorJob, error) {
	syncInterval, err := time.ParseDuration(config.Scheduler.SyncInterval)
	if err != nil {
		return nil, err
	}

	dbClient, err := db.GetDatabase(config.Postgres)
	if err != nil {
		return nil, err
	}

	scheduler := gocron.NewScheduler(time.Local)
	scheduler.StartAsync()

	monitorJob := MonitorJob{
		syncInterval:      syncInterval,
		UrlsToSchedule:    make(chan db.Url, 100),
		scheduledUrls:     make(map[uint]*gocron.Job),
		failesCountPerUrl: make(map[uint]int),
		dbClient:          dbClient,
		scheduler:         scheduler,
	}

	if err = monitorJob.syncUrls(); err != nil {
		log.Error("Error in syncing urls: ", err)
	}
	return &monitorJob, nil
}

func (monitorJob *MonitorJob) RunForever() {
	ticker := time.NewTicker(monitorJob.syncInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := monitorJob.syncUrls(); err != nil {
					log.Error("Error in syncing urls: ", err)
				}
			case url := <-monitorJob.UrlsToSchedule:
				log.Debug("new url! ", url.ID, url.Address)
				monitorJob.scheduleUrl(url)
			}
		}
	}()
}

func (monitorJob *MonitorJob) scheduleUrl(url db.Url) error {
	if _, ok := monitorJob.scheduledUrls[url.ID]; ok {
		return UrlAlreadyScheduled
	}

	job, err := monitorJob.scheduler.Every(url.Interval).Do(monitorJob.monitorUrl, url)
	if err != nil {
		return err
	}

	monitorJob.scheduledUrls[url.ID] = job
	monitorJob.failesCountPerUrl[url.ID] = 0
	return nil
}

func (monitorJob *MonitorJob) syncUrls() error {
	allUrls, _, err := monitorJob.dbClient.GetUrls()
	if err != nil {
		return err
	}

	for _, url := range allUrls {
		if _, ok := monitorJob.scheduledUrls[url.ID]; !ok {
			monitorJob.UrlsToSchedule <- url
		}
	}
	return nil
}

func (monitorJob *MonitorJob) monitorUrl(url db.Url) error {
	request := monitorJob.sendRequest(url)
	if request.StatusCode/100 != 2 {
		log.Debug("REQUEST FAILED: ", url.Address, " ", request.StatusCode)
		if monitorJob.failesCountPerUrl[url.ID]+1 == url.FailureThreshold {
			monitorJob.fireAlert(url)
		}
	} else if monitorJob.failesCountPerUrl[url.ID] > 0 {
		monitorJob.resolveAlert(url)
	}

	log.Debug(url.Address, " health: ", request.StatusCode)
	return monitorJob.dbClient.SaveRequest(request)
}

func (monitorJob *MonitorJob) sendRequest(url db.Url) db.Request {
	request := db.Request{
		UrlID: int(url.ID),
		Time:  time.Now(),
	}

	address := url.Address
	if !strings.Contains("http", address) {
		address = "http://" + url.Address
	}
	resp, err := http.Get(address)
	if err != nil {
		request.StatusCode = NoResponseCode
	} else {
		request.StatusCode = resp.StatusCode
	}
	return request
}

func (monitorJob *MonitorJob) fireAlert(url db.Url) error {
	log.Debug(" **** FIRING ALARM ***** ", url.Address)
	alert := db.Alert{
		UrlID:     int(url.ID),
		TimeFired: time.Now(),
	}
	if err := monitorJob.dbClient.SaveAlert(alert); err != nil {
		log.Errorln("Error in firing alert ", err)
		return err
	}
	monitorJob.failesCountPerUrl[url.ID] += 1
	return nil
}

func (monitorJob *MonitorJob) resolveAlert(url db.Url) error {
	alert, err := monitorJob.dbClient.GetLastAlert()
	if err != nil {
		return err
	}

	now := time.Now()
	resolvedAlert := db.Alert{
		TimeResolved: &now,
	}
	if err := monitorJob.dbClient.UpdateAlert(alert, resolvedAlert); err != nil {
		log.Errorln("Error in resolving alert ", err)
		return err
	}
	monitorJob.failesCountPerUrl[url.ID] = 0
	return nil
}
