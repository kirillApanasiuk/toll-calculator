package main

import (
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
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer([]string{kafkaConfig.CONST_HOST}, config)
	if err != nil {
		log.Fatal("failed to initialize NewSyncProducer, err:", err)
		return
	}
	defer producer.Close()

	for i := 0; i < 5; i++ {
		str := fmt.Sprint("hello kirill apanasiuk")
		msg := &sarama.ProducerMessage{
			Topic:     kafkaConfig.CONST_TOPIC,
			Partition: -1,
			Value:     sarama.StringEncoder(str),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Println("SendMessage err: ", err)
			return
		}
		log.Printf("[producer] partition id: %d; offset:%d, value: %s\n", partition, offset, str)
	}

	return
	receiver := NewDataReceiver()
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
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgChannel: make(chan types.OBUData, 128),
	}
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
		fmt.Printf("recieved OBU data from [%d] :: <lat %2.f, long %2f>\n", data.OBUID, data.Lat, data.Long)
		//dr.msgChannel <- data
	}
}
