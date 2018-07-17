package constants

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/ChinmayR/forbot/crawler"
)

const (
	IS_LAMBDA = iota
	IS_BACKTEST
	NORMAL
)

const ExecutionType = IS_BACKTEST
const FILENAME = "allPoints"

const (
	ACCOUNT_ID = "509983"
	AUTH_TOKEN = "cb9b68b71444f1c24aa5c9e9e40ae409-3603616a28fd4e42580d23d90856ecd5"

	RFC3339      = "2006-01-02T15:04:05Z07:00"
	IMAGE_FORMAT = "01-02T15-04"
)

const (
	EUR_USD = "EUR_USD"
	GBP_USD = "GBP_USD"
	USD_JPY = "USD_JPY"
	EUR_JPY = "EUR_JPY"
	AUD_USD = "AUD_USD"
)

var Symbols = []string{EUR_USD, GBP_USD, USD_JPY, EUR_JPY, AUD_USD}

type HourRange struct {
	StartHour int
	EndHour   int
	Pip       float64
}

var ValidHours = map[string]HourRange{
	"EUR_USD": HourRange{
		StartHour: 0,
		EndHour:   24,
		Pip:       1 / 10000.0,
	},
	"GBP_USD": HourRange{
		StartHour: 0,
		EndHour:   24,
		Pip:       1 / 10000.0,
	},
	"USD_JPY": HourRange{
		StartHour: 0,
		EndHour:   24,
		Pip:       1 / 100.0,
	},
	"EUR_JPY": HourRange{
		StartHour: 0,
		EndHour:   24,
		Pip:       1 / 100.0,
	},
}

func GetValidHourForSymbol(symbol string) (int, int) {
	validHoursForSymbol := ValidHours[symbol]
	return validHoursForSymbol.StartHour, validHoursForSymbol.EndHour
}

type StopRunPoint struct {
	Symbol string
	Year   int
	Month  int
	Day    int
	Values []float64
}

var StopRunPoints_EUR_USD = []StopRunPoint{
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 15, Values: []float64{1.2521, 1.2388, 1.2303}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 16, Values: []float64{1.2521, 1.2458, 1.2388}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 19, Values: []float64{1.2456}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 20, Values: []float64{1.2435}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 21, Values: []float64{1.2372}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 22, Values: []float64{1.2372, 1.2435, 1.2212}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 23, Values: []float64{1.2319, 1.2297}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 26, Values: []float64{}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 27, Values: []float64{1.2354, 1.2435, 1.2260, 1.2212}},
	//{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 28, Values: []float64{1.2260, 1.2282}},
	//{Symbol: EUR_USD, Year: 2018, Month: 3, Day: 29, Values: []float64{1.2355, 1.2375}},
	//{Symbol: EUR_USD, Year: 2018, Month: 4, Day: 9, Values: []float64{1.2342, 1.2375, 1.2261, 1.2219, 1.2168}},
	//{Symbol: EUR_USD, Year: 2018, Month: 4, Day: 10, Values: []float64{1.2289}},
	//
	//{Symbol: EUR_USD, Year: 2018, Month: 6, Day: 6, Values: []float64{1.1738, 1.1824, 1.1645}},
	//{Symbol: EUR_USD, Year: 2018, Month: 6, Day: 7, Values: []float64{1.1789, 1.1824, 1.1742, 1.1645}},
	//{Symbol: EUR_USD, Year: 2018, Month: 6, Day: 8, Values: []float64{1.1838, 1.1742, 1.1645}},
}

var StopRunPoints_GBP_USD = []StopRunPoint{
	//{Symbol: GBP_USD, Year: 2018, Month: 4, Day: 5, Values: []float64{1.4094, 1.4200, 1.4015}},
	//{Symbol: GBP_USD, Year: 2018, Month: 4, Day: 6, Values: []float64{1.39898}},
	//{Symbol: GBP_USD, Year: 2018, Month: 6, Day: 5, Values: []float64{1.3350, 1.3330}},
	//{Symbol: GBP_USD, Year: 2018, Month: 6, Day: 6, Values: []float64{1.3412, 13477, 1.3335, 1.3303, 1.3259}},
	//{Symbol: GBP_USD, Year: 2018, Month: 6, Day: 7, Values: []float64{1.3436, 1.3477, 1.3303, 1.3259}},
	//{Symbol: GBP_USD, Year: 2018, Month: 6, Day: 8, Values: []float64{1.3469, 1.3379, 1.3303, 1.3259}},
}

var StopRunPoints_USD_JPY = []StopRunPoint{
	//{Symbol: EUR_USD, Year: 2018, Month: 6, Day: 1, Values: []float64{109.07, 109.81, 108.12}},
}

var StopRunPoints_EUR_JPY = []StopRunPoint{
	//{Symbol: EUR_JPY, Year: 2018, Month: 4, Day: 5, Values: []float64{131.77, 132.32, 130.87, 130.01}},
	//{Symbol: EUR_JPY, Year: 2018, Month: 4, Day: 6, Values: []float64{131.77, 132.32, 130.87, 130.01}},
	//{Symbol: EUR_JPY, Year: 2018, Month: 4, Day: 9, Values: []float64{131.77, 132.32, 130.87, 130.01}},
}

var StopRunPoints_AUD_USD = []StopRunPoint{
	//{Symbol: AUD_USD, Year: 2018, Month: 6, Day: 1, Values: []float64{0.7555, 0.7531}},
}

func HasStopRunPointsForToday(symbol string) bool {
	crawledPoints, err := crawler.GetTodayManipulationPoints()
	if err != nil {
		panic(err)
	}
	for _, manPoint := range crawledPoints[0].ManPoints {
		if manPoint.Symbol == symbol {
			return true
		}
	}
	return false
}

func GetStopRunPointsForSymbol(symbol string) []StopRunPoint {
	retVal := make([]StopRunPoint, 0)

	crawledPoints := make([]*crawler.CrawledPoints, 0)

	if ExecutionType == IS_BACKTEST {
		pointsFile, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			return retVal
		}
		b, err := ioutil.ReadAll(pointsFile)
		if err != nil {
			log.Println(err)
			return retVal
		}
		json.Unmarshal(b, &crawledPoints)
	} else {
		points, err := crawler.GetTodayManipulationPoints()
		if err != nil {
			panic(err)
		}
		crawledPoints = points
	}

	for _, crawledPoint := range crawledPoints {
		for _, manPoint := range crawledPoint.ManPoints {
			if manPoint.Symbol == symbol {
				forDay := crawledPoint.Date
				log.Printf("Adding crawled data for %v day:%v month:%v year:%v, points: %v\n", manPoint.Symbol, forDay.Day(), forDay.Month(), forDay.Year(), manPoint.Points)
				retVal = append(retVal,
					StopRunPoint{
						Symbol: manPoint.Symbol,
						Year:   forDay.Year(),
						Month:  int(forDay.Month()),
						Day:    forDay.Day(),
						Values: manPoint.Points,
					})
			}
		}
	}

	switch symbol {
	case EUR_USD:
		retVal = append(retVal, StopRunPoints_EUR_USD...)
	case GBP_USD:
		retVal = append(retVal, StopRunPoints_GBP_USD...)
	case USD_JPY:
		retVal = append(retVal, StopRunPoints_USD_JPY...)
	case EUR_JPY:
		retVal = append(retVal, StopRunPoints_EUR_JPY...)
	case AUD_USD:
		retVal = append(retVal, StopRunPoints_AUD_USD...)
	default:
		retVal = append(retVal, StopRunPoints_EUR_USD...)
	}
	return retVal
}
