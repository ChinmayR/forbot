package algorithm

import (
	"math"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/constants"
)

func getCurrentTime(timeSeries []time.Time) time.Time {
	return timeSeries[len(timeSeries)-1]
}

func GetClosestXvIndex(ga forbot.GraphAnalysis, atTime time.Time) int {
	closestXvIndex := 0
	for i, _ := range ga.Xv {
		dur1 := ga.Xv[closestXvIndex].Sub(atTime)
		dur2 := ga.Xv[i].Sub(atTime)
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
