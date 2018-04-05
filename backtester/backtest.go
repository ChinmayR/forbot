package backtester

import (
	"math"
	"net/http"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm"
	"github.com/wcharczuk/go-chart"
)

type HandlerFunc func(res http.ResponseWriter, req *http.Request)

func RunBacktest(symbol string, from time.Time, to time.Time) HandlerFunc {
	ga := forbot.GetGraphAnalysisForSymbol(symbol, from, to)

	graph := ga.GetGraph()
	graph.Title = "Backtesting"

	//location, _ := time.LoadLocation("UTC")
	for i, xV := range ga.Xv {
		if algorithm.RunBasicAlgo(ga, xV) {
			graph.Series = append(graph.Series,
				DrawCircleAtTime(ga, xV, ga.YvClose[i])...)
		}
	}

	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "image/png")
		graph.Render(chart.PNG, res)
	}
}

func DrawCircleAtTime(ga forbot.GraphAnalysis, atTime time.Time, value float64) []chart.Series {
	singleDayDataX := make([]time.Time, 0)
	singleDayDataY := make([]float64, 0)
	singleDayDataY2 := make([]float64, 0)

	closestXvIndex := algorithm.GetClosestXvIndex(ga, atTime)

	length := 6
	for j := -length; j <= length; j++ {
		if closestXvIndex+j >= len(ga.Xv) || closestXvIndex+j < 0 {
			continue
		}
		singleDayDataX = append(singleDayDataX, ga.Xv[closestXvIndex+j])
		val1 := value - math.Abs(float64(j)*0.00005) + 0.0003
		val2 := value + math.Abs(float64(j)*0.00005) - 0.0003
		singleDayDataY = append(singleDayDataY, val1)
		singleDayDataY2 = append(singleDayDataY2, val2)
	}

	retVal := make([]chart.Series, 0)

	retVal = append(retVal, chart.TimeSeries{
		Name: "Dot",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(2),
			StrokeWidth: 3,
		},
		XValues: singleDayDataX,
		YValues: singleDayDataY,
	})

	retVal = append(retVal, chart.TimeSeries{
		Name: "Dot2",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(2),
			StrokeWidth: 3,
		},
		XValues: singleDayDataX,
		YValues: singleDayDataY2,
	})

	return retVal
}
