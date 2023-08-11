package main

import (
	config2 "github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/IBM/sarama.v1"
)

// This can also be called KafkaTransport
type KafkaSaramaConsumer struct {
	consumer    *sarama.Consumer
	calcService CalculatorServicer
}

func NewKafkaConsumer(host string, calcService CalculatorServicer) (*KafkaSaramaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{host}, config)
	if err != nil {
		logrus.Errorf("Error while configuring consumer")
		return nil, err
	}

	logrus.Info("start consuming ...")

	return &KafkaSaramaConsumer{
		consumer:    &consumer,
		calcService: calcService,
	}, nil
}

func (s *KafkaSaramaConsumer) SubscribeLoop(
	topic string, errorChan chan *sarama.ConsumerError,
	messageChan chan *sarama.ConsumerMessage,
) {
	consumer := *s.consumer
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logrus.Printf("Error while getting partitions \n %v", err)
	}

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(config2.CONST_TOPIC, partition, sarama.OffsetOldest)

		go func() {
			if err != nil {
				logrus.Printf("couldn't create partition consumer, err :%s", err.Error())
			}

			for {
				select {
				case myError := <-pc.Errors():
					errorChan <- myError

				case msg := <-pc.Messages():
					messageChan <- msg
				}
			}
		}()
	}
}
