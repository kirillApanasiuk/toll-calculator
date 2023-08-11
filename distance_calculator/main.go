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
		calcServicer = NewCalculatorService()
		consumer, _  = NewKafkaConsumer(config.CONST_HOST, calcServicer)
	)

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

			logrus.Infof("recieved message %v", string(msg.Value))

			res, err := calcService.CalculateDistance(msg.Value)
			if err != nil {
				logrus.Errorf("Custom message %v \n", err)
			}

			logrus.Printf("Youhooo my data %f.2", res)
		}
	}
}
