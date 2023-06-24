package logger

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func LoggerInfo(status string) {
	log.WithFields(log.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
	}).Info(status)
}

func LoggerFatal(status string)  {
	log.WithFields(log.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
	}).Fatal(status)
}

func LoggerWarn(status string)  {
	log.WithFields(log.Fields{
		"at":time.Now().Format("2006-01-02 15:04:05"),
	}).Warn(status)
}