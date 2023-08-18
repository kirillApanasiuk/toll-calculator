package main

import (
	"context"
	"encoding/json"
	"github.com/kirillApanasiuk/toll-calculator/aggregator/client"
	"github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/IBM/sarama.v1"
	"log"
	"time"
)

func main() {
	var (
		errorChan    = make(chan *sarama.ConsumerError)
		messageChan  = make(chan *sarama.ConsumerMessage)
		calcServicer CalculatorServicer
		aggClient    client.Client
	)
	calcServicer = NewCalculatorService()
	//wrap calc service with log middleware
	calcServicer = NewLogMiddleware(calcServicer)
	aggClient = client.NewHttpClient(config.AGGREGATOR_ENDPOINT)
	grpcClient, err := client.NewGRPCClient(config.GRPC_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}

	_ = aggClient
	consumer, _ := NewKafkaConsumer(config.CONST_HOST, calcServicer)
	consumer.SubscribeLoop(config.CONST_TOPIC, errorChan, messageChan)
	go makeSomethingWithSubscribedData(errorChan, messageChan, calcServicer, grpcClient)

	c := make(chan struct{})
	<-c
}

func makeSomethingWithSubscribedData(
	errorChan chan *sarama.ConsumerError, messsageChan chan *sarama.ConsumerMessage,
	calcService CalculatorServicer, ct client.Client,
) {
	for {
		select {
		case err := <-errorChan:
			logrus.Errorf("err %s", err.Error())
		case msg := <-messsageChan:
			var obuData types.OBUData
			if err := json.Unmarshal(msg.Value, &obuData); err != nil {
				logrus.Error(err)
			}
			distance, err := calcService.CalculateDistance(msg.Value)

			if err != nil {
				logrus.Errorf("calculation error: %s", err)
				continue
			}
			req := types.AggregateRequest{
				Value: distance,
				ObuId: int32(obuData.OBUID),
				Unix:  time.Now().UnixNano(),
			}
			if err := ct.Aggregate(context.Background(), &req); err != nil {
				logrus.Errorf("aggregate error:", err)
			}
		}
	}
}
