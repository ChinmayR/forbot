package forbot

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ChinmayR/forbot/params"
)

type Candle struct {
	Open  string `json:"o"`
	Close string `json:"c"`
	Low   string `json:"l"`
	High  string `json:"h"`
}

type Candles struct {
	Complete bool      `json:"complete"`
	Volume   int       `json:"volume"`
	Time     time.Time `json:"time"`
	Mid      Candle    `json:"mid"`
}

type InstrumentHistory struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}

type Bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}

type BrokerBook struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []Bucket  `json:"buckets"`
}

func (c *OandaConnection) GetCandles(instrument string, params params.InstrumentCandlesParams) InstrumentHistory {
	endpoint := "/instruments/" + instrument + "/candles"
	candles := c.Request(endpoint, params.ToValues())
	data := InstrumentHistory{}
	unmarshalJson(candles, &data)

	return data
}

func (c *OandaConnection) OrderBook(instrument string) BrokerBook {
	endpoint := "/instruments/" + instrument + "/orderBook"
	orderbook := c.Request(endpoint, url.Values{})
	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data
}

func (c *OandaConnection) PositionBook(instrument string) BrokerBook {
	endpoint := "/instruments/" + instrument + "/positionBook"
	orderbook := c.Request(endpoint, url.Values{})
	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data
}

func GetTimesFromCandles(candles []Candles) []time.Time {
	var retVal = make([]time.Time, 0)
	for _, point := range candles {
		retVal = append(retVal, point.Time)
	}
	return retVal
}

func GetCloseFromCandles(candles []Candles) []float64 {
	var retVal = make([]float64, 0)
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.Close, 64)
		retVal = append(retVal, conv)
	}
	return retVal
}

func GetOpenFromCandles(candles []Candles) []float64 {
	var retVal = make([]float64, 0)
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.Open, 64)
		retVal = append(retVal, conv)
	}
	return retVal
}

func GetLowFromCandles(candles []Candles) []float64 {
	var retVal = make([]float64, 0)
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.Low, 64)
		retVal = append(retVal, conv)
	}
	return retVal
}

func GetHighFromCandles(candles []Candles) []float64 {
	var retVal = make([]float64, 0)
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.High, 64)
		retVal = append(retVal, conv)
	}
	return retVal
}

func GetMinFromCandles(candles []Candles) float64 {
	conv, _ := strconv.ParseFloat(candles[0].Mid.Low, 64)
	var retVal = conv
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.Low, 64)
		if conv < retVal {
			retVal = conv
		}
	}
	return retVal
}

func GetMaxFromCandles(candles []Candles) float64 {
	conv, _ := strconv.ParseFloat(candles[0].Mid.High, 64)
	var retVal = conv
	for _, point := range candles {
		conv, _ := strconv.ParseFloat(point.Mid.High, 64)
		if conv > retVal {
			retVal = conv
		}
	}
	return retVal
}
