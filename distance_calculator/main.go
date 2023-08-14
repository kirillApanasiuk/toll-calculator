package main

import (
	"github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/IBM/sarama.v1"
)

func main() {
	var (
		errorChan    = make(chan *sarama.ConsumerError)
		messageChan  = make(chan *sarama.ConsumerMessage)
		calcServicer CalculatorServicer
	)
	calcServicer = NewCalculatorService()
	//wrap calc service with log middleware
	calcServicer = NewLogMiddleware(calcServicer)

	consumer, _ := NewKafkaConsumer(config.CONST_HOST, calcServicer)
	consumer.SubscribeLoop(config.CONST_TOPIC, errorChan, messageChan)
	go makeSomethingWithSubscribedData(errorChan, messageChan, calcServicer)

	c := make(chan struct{})
	<-c
}

func makeSomethingWithSubscribedData(
	errorChan chan *sarama.ConsumerError, messsageChan chan *sarama.ConsumerMessage,
	calcService CalculatorServicer,
) {
	for {
		select {
		case err := <-errorChan:
			logrus.Errorf("err %s", err.Error())
		case msg := <-messsageChan:
			res, _ := calcService.CalculateDistance(msg.Value)
			_ = res
		}
	}
}
