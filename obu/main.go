package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"log"
	"math"
	"math/rand"
	"time"
)

var sendInterval = time.Second

const wsEndpoint = "ws://127.0.0.1:3000/ws"

func genCoord() float64 {
	n := rand.Intn(100) + 1
	f := rand.Float64()
	return float64(n) + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	obuIds := generateOBUIDS(20)
	for {
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIds[i],
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
