package main

import (
	"fmt"
	"github.com/kirillApanasiuk/toll-calculator/config"
	"gopkg.in/IBM/sarama.v1"
)

func main() {
	errorChan := make(chan *sarama.ConsumerError)
	messageChan := make(chan *sarama.ConsumerMessage)
	consumer, _ := NewKafkaConsumer()
	consumer.Subscribe(config.CONST_TOPIC, errorChan, messageChan)
	go makeSomethingWithSubscribedData(errorChan, messageChan)
	c := make(chan struct{})
	<-c
}

func makeSomethingWithSubscribedData(errorChan chan *sarama.ConsumerError, messsageChan chan *sarama.ConsumerMessage) {
	for {
		select {
		case err := <-errorChan:
			fmt.Errorf("err %s", err.Error())
		case msg := <-messsageChan:
			fmt.Printf("recieved message %v", string(msg.Value))
		}
	}
}
