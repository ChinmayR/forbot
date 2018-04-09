package draw

import (
	"math"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm/utils"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func DrawCircleAtTime(ga forbot.GraphAnalysis, atTime time.Time, value float64, color drawing.Color) []chart.Series {
	singleDayDataX := make([]time.Time, 0)
	singleDayDataY := make([]float64, 0)
	singleDayDataY2 := make([]float64, 0)

	closestXvIndex := utils.GetClosestXvIndex(ga.Xv, atTime)

	//length := int(float64(ga.GetGraph().Width) * 0.05)
	length := 2
	for j := -length; j <= length; j++ {
		if closestXvIndex+j >= len(ga.Xv) || closestXvIndex+j < 0 {
			continue
		}
		singleDayDataX = append(singleDayDataX, ga.Xv[closestXvIndex+j])
		val1 := value - math.Abs(float64(j)*0.00005) + (0.00005 * float64(length))
		val2 := value + math.Abs(float64(j)*0.00005) - (0.00005 * float64(length))
		singleDayDataY = append(singleDayDataY, val1)
		singleDayDataY2 = append(singleDayDataY2, val2)
	}

	retVal := make([]chart.Series, 0)

	retVal = append(retVal, chart.TimeSeries{
		Name: "Dot",
		Style: chart.Style{
			Show:        true,
			StrokeColor: color,
			StrokeWidth: 3,
		},
		XValues: singleDayDataX,
		YValues: singleDayDataY,
	})

	retVal = append(retVal, chart.TimeSeries{
		Name: "Dot2",
		Style: chart.Style{
			Show:        true,
			StrokeColor: color,
			StrokeWidth: 3,
		},
		XValues: singleDayDataX,
		YValues: singleDayDataY2,
	})

	return retVal
}
