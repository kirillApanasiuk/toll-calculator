package main

import (
	"fmt"
	"math/rand"
	"time"
)

const sendInterval = 10

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func genCoord() float64 {
	n := rand.Intn(100) + 1
	f := rand.Float64()
	return float64(n) + f
}

func main() {
	for {
		fmt.Println(genCoord())
		time.Sleep(sendInterval * time.Millisecond)
	}
}
