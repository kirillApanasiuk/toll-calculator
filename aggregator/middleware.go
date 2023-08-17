package main

import (
	"github.com/kirillApanasiuk/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next Aggregator
}

func (m *LogMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	currentTime := time.Now()
	defer func(start time.Time) {
		var (
			distance      float64
			amount        float64
			obuIdentifier int
		)

		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
			obuIdentifier = inv.OBUID
		}

		logrus.WithFields(
			logrus.Fields{
				"took":        time.Since(start).Milliseconds(),
				"err":         err,
				"obuId":       obuIdentifier,
				"totalDist":   distance,
				"totalAmount": amount,
			},
		).Info("CalculateInvoice")
	}(currentTime)
	inv, err = m.next.CalculateInvoice(obuId)
	return
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took": time.Since(start),
				"err":  err,
			},
		).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
