package forbot

import (
	"fmt"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/ChinmayR/forbot/constants"
	"github.com/ChinmayR/forbot/params"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/util"
)

var oandaCon = NewConnection(constants.ACCOUNT_ID, constants.AUTH_TOKEN, true)

type GraphAnalysis struct {
	Symbol    string
	Xv        []time.Time
	YvClose   []float64
	YvOpen    []float64
	YvLow     []float64
	YvHigh    []float64
	MinYRange float64
	MaxYRange float64
}

func GetGraphAnalysisForSymbol(symbol string, from, to time.Time) GraphAnalysis {
	//location, _ := time.LoadLocation("UTC")
	history := oandaCon.GetCandles(symbol,
		params.InstrumentCandlesParams{
			Granularity: "M15",
			From:        from,
			To:          to,
		})

	graphAnalysis := GraphAnalysis{
		Symbol:    symbol,
		Xv:        GetTimesFromCandles(history.Candles),
		YvClose:   GetCloseFromCandles(history.Candles),
		YvOpen:    GetOpenFromCandles(history.Candles),
		YvLow:     GetLowFromCandles(history.Candles),
		YvHigh:    GetHighFromCandles(history.Candles),
		MinYRange: GetMinFromCandles(history.Candles),
		MaxYRange: GetMaxFromCandles(history.Candles),
	}
	return graphAnalysis
}

func (ga *GraphAnalysis) Handler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "image/png")
	ga.GetGraph().Render(chart.PNG, res)
}

func (ga *GraphAnalysis) SaveGraph(fileNamePrefix string) {
	collector := &chart.ImageWriter{}
	ga.GetGraph().Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	curTime := time.Now()
	fileName := fileNamePrefix + "-image-" + curTime.Format(constants.IMAGE_FORMAT) + ".png"
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, image)
}

func (ga *GraphAnalysis) GetGraph() chart.Chart {
	xv := ga.Xv

	priceSeriesLow := chart.TimeSeries{
		Name: "Low",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xv,
		YValues: ga.YvLow,
	}

	priceSeriesHigh := chart.TimeSeries{
		Name: "High",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xv,
		YValues: ga.YvHigh,
	}

	//supAndResSeries, horizontalGridlines := ga.getSupAndResLevelsForAllDays()
	supAndResSeries, horizontalGridlines := []chart.Series{}, []chart.GridLine{}

	graph := chart.Chart{
		Width:  1280,
		Height: 800,
		XAxis: chart.XAxis{
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			GridLines: ga.getVerticalGridLines(),
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.5f", vf)
				}
				return ""
			},
			Range: &chart.ContinuousRange{
				Max: ga.MaxYRange,
				Min: ga.MinYRange,
			},
			GridMajorStyle: chart.Style{
				Show:        false,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			GridLines: horizontalGridlines,
		},
		Series: append(append([]chart.Series{
			priceSeriesLow,
			priceSeriesHigh,
		}, ga.getCustomLevelsForAllDays()...), supAndResSeries...),
	}

	return graph
}

func (ga *GraphAnalysis) getVerticalGridLines() []chart.GridLine {
	gridLines := make([]chart.GridLine, 0)
	for i, timeVline := range ga.Xv {
		if i < len(ga.Xv)-1 && ga.Xv[i].YearDay() != ga.Xv[i+1].YearDay() {
			gridLines = append(gridLines, chart.GridLine{Value: util.Time.ToFloat64(timeVline)})
		}
	}
	return gridLines
}

func (ga *GraphAnalysis) getCustomLevelsForAllDays() []chart.Series {
	var retVal = make([]chart.Series, 0)
	for _, eachDay := range constants.GetStopRunPointsForSymbol(ga.Symbol) {
		for _, eachLevel := range eachDay.Values {
			series := ga.GetLevelsForDay(eachDay.Year, eachDay.Month, eachDay.Day, eachLevel)
			retVal = append(retVal, series)
		}
	}
	return retVal
}

func (ga *GraphAnalysis) GetLevelsForDay(year, month, day int, value float64) chart.Series {
	singleDayDataX := make([]time.Time, 0)
	singleDayDataY := make([]float64, 0)
	for i, _ := range ga.Xv {
		yearCur, monthCur, dayCur := ga.Xv[i].Date()
		if yearCur == year && int(monthCur) == month && dayCur == day {
			singleDayDataX = append(singleDayDataX, ga.Xv[i])
			singleDayDataY = append(singleDayDataY, value)
		}
	}

	return chart.TimeSeries{
		Name: "DayLine",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(2),
		},
		XValues: singleDayDataX,
		YValues: singleDayDataY,
	}
}

type Point struct {
	XVal     time.Time
	YVal     float64
	Strength int
}

func (ga *GraphAnalysis) getSupAndResLevelsForAllDays() ([]chart.Series, []chart.GridLine) {
	var retVal = make([]chart.Series, 0)
	var gridLines = make([]chart.GridLine, 0)

	for _, eachPoint := range ga.FindSupAndRes() {
		fmt.Printf("Date: %s, Level: %f, Strength: %d\n", eachPoint.XVal.String(), eachPoint.YVal, eachPoint.Strength)

		xValues := make([]time.Time, 0)
		yValues := make([]float64, 0)
		for _, timeVal := range ga.Xv {
			if isInSameDay(eachPoint.XVal, timeVal) {
				xValues = append(xValues, timeVal)
				yValues = append(yValues, eachPoint.YVal)
			}
		}

		retVal = append(retVal, chart.TimeSeries{
			Name: "DayLine",
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetDefaultColor(4),
			},
			XValues: xValues,
			YValues: yValues,
		})

		gridLines = append(gridLines, chart.GridLine{Value: eachPoint.YVal})
	}
	return retVal, gridLines
}

func isInSameDay(i, j time.Time) bool {
	return i.Year() == j.Year() && i.Month() == j.Month() && i.Day() == j.Day()
}

func (ga *GraphAnalysis) FindSupAndRes() []Point {
	// number of points to consider on each side of the current point
	numPointsToConsider := 40
	strengthThreshold := 1
	var localMaxOrMin = make([]Point, 0)

	for i, point := range ga.YvClose {
		if i < numPointsToConsider || i > (len(ga.YvClose)-1-numPointsToConsider) {
			continue
		}
		if isMaxima(ga.YvClose, i, numPointsToConsider) || isMinima(ga.YvClose, i, numPointsToConsider) {
			localMaxOrMin = append(localMaxOrMin, Point{XVal: ga.Xv[i], YVal: point, Strength: 1})
		}
	}

	var retVals = make([]Point, 0)
	threshold := 0.001 // 10 pips threshold
	for i, maxOrMin := range localMaxOrMin {
		for j, otherMaxOrMin := range localMaxOrMin {
			if i == j {
				continue
			}
			if math.Abs(maxOrMin.YVal-otherMaxOrMin.YVal) < threshold {
				localMaxOrMin[i].Strength += 1
			}
		}
		if localMaxOrMin[i].Strength >= strengthThreshold {
			// loop over all the existing found support or resistance and
			// only add this one if none of the already found are not within
			// the given range
			isNewSupOrRes := true
			for _, alreadyFound := range retVals {
				if math.Abs(alreadyFound.YVal-localMaxOrMin[i].YVal) < threshold {
					isNewSupOrRes = false
				}
			}
			if isNewSupOrRes {
				retVals = append(retVals, localMaxOrMin[i])
			}
		}
	}
	return retVals
}

func isMaxima(data []float64, i int, numPointsToConsider int) bool {
	for j := 1; j <= numPointsToConsider; j++ {
		if data[i-j] > data[i] || data[i+j] > data[i] {
			return false
		}
	}
	return true
}

func isMinima(data []float64, i int, numPointsToConsider int) bool {
	for j := 1; j <= numPointsToConsider; j++ {
		if data[i-j] < data[i] || data[i+j] < data[i] {
			return false
		}
	}
	return true
}
