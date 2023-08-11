package main

import (
	"encoding/json"
	"fmt"
	kafkaConfig "github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"gopkg.in/IBM/sarama.v1"
	"log"
)

type DataProducer interface {
	ProduceData(topicName string, data types.OBUData) error
}

type KafkaSaramaProducer struct {
	producer *sarama.SyncProducer
}

func NewKafkaSaramaProducer() (DataProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{kafkaConfig.CONST_HOST}, config)
	if err != nil {
		fmt.Println("error when building producer")
		return nil, err
	}
	return &KafkaSaramaProducer{
		producer: &producer,
	}, nil
}

func (saramaProducer *KafkaSaramaProducer) ProduceData(topicName string, data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     topicName,
		Partition: -1,
		Value:     sarama.StringEncoder(b),
	}
	producer := *saramaProducer.producer
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("SendMessage err: ", err)
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topicName, partition, offset)
	return nil
}
