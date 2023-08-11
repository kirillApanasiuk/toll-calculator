package main

import (
	"encoding/json"
	"fmt"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"math"
)

type CalculatorServicer interface {
	CalculateDistance(dataBytes []byte) (float64, error)
}

type CalculatorService struct {
	points [][]float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		points: make([][]float64, 0),
	}
}

func (s *CalculatorService) CalculateDistance(dataBytes []byte) (float64, error) {
	var obuData types.OBUData
	if err := json.Unmarshal(dataBytes, &obuData); err != nil {
		logrus.Error(err)
	}

	distance := 0.0

	if len(s.points) > 0 {
		prePoint := s.points[len(s.points)-1]
		distance = calculateDistance(prePoint[0], prePoint[1], obuData.Lat, obuData.Long)
	}
	s.points = append(s.points, []float64{obuData.Lat, obuData.Long})

	fmt.Printf("calculating the distance %v \n", obuData)

	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
