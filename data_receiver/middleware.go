package main

import (
	"github.com/kirillApanasiuk/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(topicName string, data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"obuId": data.OBUID,
				"lat":   data.Lat,
				"long":  data.Long,
				"took":  time.Since(start),
			},
		).Info("producing")
	}(time.Now())
	return l.next.ProduceData(topicName, data)
}
