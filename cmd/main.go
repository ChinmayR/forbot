package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/constants"
	"github.com/ChinmayR/forbot/params"
)

var oandaCon = forbot.NewConnection(constants.ACCOUNT_ID, constants.AUTH_TOKEN, true)

func main() {
	//location, _ := time.LoadLocation("UTC")
	history := oandaCon.GetCandles("EUR_USD",
		params.InstrumentCandlesParams{
			Granularity: "M15",
			From:        time.Now().AddDate(0, 0, -50),
			To:          time.Now(),
			//From: time.Date(2018, 2, 15, 0, 0, 0, 0, location),
			//To:   time.Date(2018, 2, 28, 0, 0, 0, 0, location),
		})

	graphAnalysis := &forbot.GraphAnalysis{
		Xv:        forbot.GetTimesFromCandles(history.Candles),
		YvClose:   forbot.GetCloseFromCandles(history.Candles),
		YvOpen:    forbot.GetOpenFromCandles(history.Candles),
		YvLow:     forbot.GetLowFromCandles(history.Candles),
		YvHigh:    forbot.GetHighFromCandles(history.Candles),
		MinYRange: forbot.GetMinFromCandles(history.Candles),
		MaxYRange: forbot.GetMaxFromCandles(history.Candles),
	}

	graphAnalysis.SaveGraph()
	log.Println("Server started...")

	http.HandleFunc("/", graphAnalysis.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
