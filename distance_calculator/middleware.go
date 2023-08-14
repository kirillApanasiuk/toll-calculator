package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func (l *LogMiddleware) CalculateDistance(dataBytes []byte) (dist float64, err error) {
	defer func(start time.Time) {
		var log = logrus.New()
		log.Formatter = &logrus.TextFormatter{DisableColors: true, FullTimestamp: true}
		log.WithFields(
			logrus.Fields{
				"took": time.Since(start),
				"err":  err,
				"dist": dist,
			},
		).Info("calculate distance")
	}(time.Now())

	dist, err = l.next.CalculateDistance(dataBytes)
	return
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{next: next}
}
