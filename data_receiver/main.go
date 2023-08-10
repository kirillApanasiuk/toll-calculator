package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	kafkaConfig "github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"gopkg.in/IBM/sarama.v1"
	"log"
	"net/http"
)

const KafkaTopic = "obudata"

func main() {

	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", receiver.handleWS)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("data receiver works fine")
	}
}

type DataReceiver struct {
	msgChannel chan types.OBUData
	conn       *websocket.Conn
	prod       *sarama.SyncProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer([]string{kafkaConfig.CONST_HOST}, config)
	if err != nil {
		fmt.Println("error when building producer")
		return nil, err
	}

	fmt.Println("Producer successfully builded")

	return &DataReceiver{
		msgChannel: make(chan types.OBUData, 128),
		prod:       &producer,
	}, nil
}

func (dr *DataReceiver) produceData(topic string, data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(b),
	}
	producer := *dr.prod
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("SendMessage err: ", err)
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}
func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("NewOBU connected client connected!")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Println("received message", data)
		if err := dr.produceData(kafkaConfig.CONST_TOPIC, data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
