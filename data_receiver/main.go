package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"log"
	"net/http"
)

func main() {
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWS)
	err := http.ListenAndServe(":3000", nil)
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
