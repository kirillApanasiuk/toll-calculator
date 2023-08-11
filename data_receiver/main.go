package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	kafkaConfig "github.com/kirillApanasiuk/toll-calculator/config"
	"github.com/kirillApanasiuk/toll-calculator/types"
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
	prod       DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p   DataProducer
		err error
	)
	p, err = NewKafkaSaramaProducer()
	if err != nil {
		log.Fatalf("The kafka sarama producer can't be setup -> %v \n", err)
	}

	p = NewLogMiddleware(p)
	fmt.Println("Producer successfully builded")

	return &DataReceiver{
		msgChannel: make(chan types.OBUData, 128),
		prod:       p,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("something happens here")
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
		if err := dr.prod.ProduceData(kafkaConfig.CONST_TOPIC, data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
