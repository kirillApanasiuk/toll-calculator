package main

import (
	config2 "github.com/kirillApanasiuk/toll-calculator/config"
	"gopkg.in/IBM/sarama.v1"
	"log"
)

type KafkaSaramaConsumer struct {
	consumer *sarama.Consumer
}

func NewKafkaConsumer() (*KafkaSaramaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{config2.CONST_HOST}, config)
	if err != nil {
		log.Fatal("Error while configuring consumer")
		return nil, err
	}

	log.Println("start consuming")

	return &KafkaSaramaConsumer{
		consumer: &consumer,
	}, nil
}

func (s *KafkaSaramaConsumer) Subscribe(
	topic string, errorChan chan *sarama.ConsumerError,
	messageChan chan *sarama.ConsumerMessage,
) {
	consumer := *s.consumer
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("Error while getting partitions \n %v", err)
	}

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(config2.CONST_TOPIC, partition, sarama.OffsetOldest)

		go func() {
			if err != nil {
				log.Printf("couldn't create partition consumer, err :%s", err.Error())
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
