package utils

import (
	"math"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/constants"
	"github.com/wcharczuk/go-chart"
)

type Algorithm interface {
	RunAlgo(ga forbot.GraphAnalysis, curTime time.Time) TradeSignal
	RunAlgoFromTo(ga forbot.GraphAnalysis, from, to int) TradeSignal
}

type TradeSignal struct {
	Signal       bool
	LevelCrossed float64
	SeriesToDraw []chart.Series
}

func getCurrentTime(timeSeries []time.Time) time.Time {
	return timeSeries[len(timeSeries)-1]
}

func GetClosestXvIndex(Xv []time.Time, atTime time.Time) int {
	closestXvIndex := 0
	for i, _ := range Xv {
		dur1 := Xv[closestXvIndex].Sub(atTime)
		dur2 := Xv[i].Sub(atTime)
		if math.Abs(dur1.Seconds()) > math.Abs(dur2.Seconds()) {
			closestXvIndex = i
		}
	}
	return closestXvIndex
}

func GetCustomPointForTime(allPoints []constants.StopRunPoint, atTime time.Time) constants.StopRunPoint {
	for _, point := range allPoints {
		if point.Year == atTime.Year() && point.Month == int(atTime.Month()) && point.Day == atTime.Day() {
			return point
		}
	}
	return constants.StopRunPoint{
		Values: []float64{},
	}
}
