package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm"
	"github.com/ChinmayR/forbot/algorithm/utils"
	"github.com/ChinmayR/forbot/backtester"
	"github.com/ChinmayR/forbot/constants"
	"github.com/ChinmayR/forbot/image_upload"
	"github.com/ChinmayR/forbot/twilio"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wcharczuk/go-chart"
)

type Request struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func LambdaHandler(request Request) (Response, error) {

	retval := Response{
		Message: fmt.Sprintf("Success for request Id %f", request.ID),
		Ok:      true,
	}

	PST, _ := time.LoadLocation("America/Los_Angeles")
	curTimeInPST := time.Now().In(PST)
	if curTimeInPST.Hour() >= 10 && curTimeInPST.Hour() < 21 {
		log.Printf("Skipping run for %v, not meant to run from 10AM till 9PM\n", curTimeInPST.String())
		return retval, nil
	}

	graphToSym := runMain(
		time.Now().AddDate(0, 0, -5),
		time.Now().Add(time.Minute*-1),
		algorithm.StopRunAlgo{})

	for _, graphToSymIter := range graphToSym {
		sym := graphToSymIter.Symbol

		http.HandleFunc("/"+sym, graphToSymIter.Handler)

		tradeSignal := graphToSymIter.TradeSignal
		log.Printf("Trade signal: %v", tradeSignal)
		if tradeSignal.Signal {
			fileName := SaveAnalyzedGraph(graphToSymIter.AnalyzedGraph, sym)
			image_upload.UploadImage(fileName)
			twilio.SendMsgFromData(fmt.Sprintf("%s - Crossed level %s", sym,
				strconv.FormatFloat(tradeSignal.LevelCrossed, 'f', -1, 64)), fileName)
		}

		PST, _ := time.LoadLocation("America/Los_Angeles")
		curTimeInPST := time.Now().In(PST)
		// send a text with image if it is between 10:30PM to 10:45PM so the stop points are reflected
		if constants.HasStopRunPointsForToday(graphToSymIter.Symbol) && curTimeInPST.Hour() == 22 && curTimeInPST.Minute() >= 30 && curTimeInPST.Minute() < 45 {
			fileName := SaveAnalyzedGraph(graphToSymIter.AnalyzedGraph, sym)
			image_upload.UploadImage(fileName)
			os.RemoveAll(fileName)
			twilio.SendMsgFromData("New levels for pair "+graphToSymIter.Symbol, fileName)
		}
	}
	return retval, nil
}

type GraphToSym struct {
	Symbol        string
	GraphAnalysis forbot.GraphAnalysis
	Handler       backtester.HandlerFunc
	AnalyzedGraph chart.Chart
	TradeSignal   utils.TradeSignal
}

func runMain(from, to time.Time, algoToRun utils.Algorithm) []GraphToSym {
	log.Printf("Running main now from: %v to: %v algoToRun: %v\n", from, to, algoToRun)
	var retVal []GraphToSym
	for _, sym := range constants.Symbols {
		log.Printf("Get graph analysis for symbol: %v\n", sym)
		graphAnalysis := forbot.GetGraphAnalysisForSymbol(
			sym,
			from,
			to)

		log.Printf("Get analyzed graph for symbol: %v\n", sym)
		handler, graph, tradeSignal := backtester.GetAnalyzedGraphAndHandler(
			&graphAnalysis,
			algoToRun)

		log.Printf("Done for %v, with tradeSignal: %v \n", sym, tradeSignal.Signal)
		retVal = append(retVal, GraphToSym{
			Symbol:        sym,
			GraphAnalysis: graphAnalysis,
			Handler:       handler,
			AnalyzedGraph: graph,
			TradeSignal:   tradeSignal,
		})
	}
	return retVal
}

const (
	IS_LAMBDA = iota
	IS_BACKTEST
	NORMAL
)

func main() {
	// run "GOOS=linux go build -o main && zip -r main.zip main"

	// run "GOOS=linux go build -o main"
	// then run "zip -r main.zip main"
	// then upload main.zip
	executionType := NORMAL
	log.Printf("Execution type: %v\n", executionType)
	if executionType == IS_LAMBDA {

		lambda.Start(LambdaHandler)

	} else if executionType == IS_BACKTEST {

		graphToSym := runMain(
			time.Now().AddDate(0, 0, -5),
			time.Now().Add(time.Minute*-1),
			algorithm.StopRunAlgo{})

		for _, graphToSymIter := range graphToSym {
			http.HandleFunc("/"+graphToSymIter.Symbol, graphToSymIter.Handler)
		}

		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":8080", nil))

	} else if executionType == NORMAL {

		graphToSym := runMain(
			time.Now().AddDate(0, 0, -5),
			time.Now().Add(time.Minute*-1),
			algorithm.StopRunAlgo{})

		for _, graphToSymIter := range graphToSym {
			sym := graphToSymIter.Symbol

			http.HandleFunc("/"+sym, graphToSymIter.Handler)

			tradeSignal := graphToSymIter.TradeSignal

			if tradeSignal.Signal {
				fileName := SaveAnalyzedGraph(graphToSymIter.AnalyzedGraph, sym)
				image_upload.UploadImage(fileName)
				os.RemoveAll(fileName)
				twilio.SendMsgFromData(fmt.Sprintf("%s - Crossed level %s", sym,
					strconv.FormatFloat(tradeSignal.LevelCrossed, 'f', -1, 64)), fileName)
			}

			PST, _ := time.LoadLocation("America/Los_Angeles")
			curTimeInPST := time.Now().In(PST)
			// send a text with image if it is between 10:30PM to 10:45PM so the stop points are reflected
			if constants.HasStopRunPointsForToday(graphToSymIter.Symbol) && curTimeInPST.Hour() == 22 && curTimeInPST.Minute() >= 30 && curTimeInPST.Minute() < 45 {
				fileName := SaveAnalyzedGraph(graphToSymIter.AnalyzedGraph, sym)
				image_upload.UploadImage(fileName)
				os.RemoveAll(fileName)
				twilio.SendMsgFromData("New levels for pair "+graphToSymIter.Symbol, fileName)
			}
		}

		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

// this method is the same as chart.go SaveGraph(...)
func SaveAnalyzedGraph(graph chart.Chart, fileNamePrefix string) string {
	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	curTime := time.Now()
	fileName := "/tmp/" + fileNamePrefix + "-image-" + curTime.Format(constants.IMAGE_FORMAT) + ".png"
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, image)
	return fileName
}
