package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/yektasrk/http-monitor/configs"
)

func ConfigureLogger(config configs.LoggerConfiguration) error {
	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2000-01-02 15:04:05",
	})
	return nil
}
