package model

import (
	"github.com/sirupsen/logrus"
	"j2pay-server/pkg/logger"
	"log"
)

type GormLogger log.Logger

func (l *GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logger.Logger.WithFields(
			logrus.Fields{
				//"module":  "gorm",
				//"type":    "sql",
				"rows": v[5],
				//"src_ref": v[1],
				"values": v[4],
			},
		).Info(v[3])
	case "log":
		logger.Logger.WithFields(
			logrus.Fields{
				//"module":  "gorm",
				//"type":    "sql",
			},
		).Error(v[2])
	default:
		logger.Logger.Error(v...)
	}
}
