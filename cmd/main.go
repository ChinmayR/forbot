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
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func LambdaHandler(request Request) (Response, error) {
	fmt.Println("Running main now")
	fmt.Println(runMain(
		time.Now().AddDate(0, 0, -2),
		time.Now().Add(time.Minute*-1),
		algorithm.BasicAlgo{}))

	fmt.Println("Ran main successfully")
	return Response{
		Message: fmt.Sprintf("Success for request Id %f", request.ID),
		Ok:      true,
	}, nil
}

type GraphToSym struct {
	Symbol        string
	GraphAnalysis forbot.GraphAnalysis
	Handler       backtester.HandlerFunc
	AnalyzedGraph chart.Chart
	TradeSignal   utils.TradeSignal
}

func runMain(from, to time.Time, algoToRun utils.Algorithm) []GraphToSym {
	var retVal []GraphToSym
	for _, sym := range constants.Symbols {
		graphAnalysis := forbot.GetGraphAnalysisForSymbol(
			sym,
			from,
			to)

		handler, graph, tradeSignal := backtester.GetAnalyzedGraphAndHandler(
			sym,
			from,
			to,
			algoToRun)

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
	NORMAL = iota
	IS_LAMBDA
	IS_BACKTEST
)

func main() {
	// run "GOOS=linux go build -o main"
	executionType := IS_BACKTEST
	if executionType == IS_LAMBDA {

		lambda.Start(LambdaHandler)

	} else if executionType == IS_BACKTEST {

		graphToSym := runMain(
			time.Now().AddDate(0, 0, -5),
			time.Now().Add(time.Minute*-1),
			algorithm.BasicAlgo{})

		for _, graphToSymIter := range graphToSym {
			http.HandleFunc("/"+graphToSymIter.Symbol, graphToSymIter.Handler)
		}

		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":8080", nil))

	} else if executionType == NORMAL {

		graphToSym := runMain(
			time.Now().AddDate(0, 0, -5),
			time.Now().Add(time.Minute*-1),
			algorithm.BasicAlgo{})

		for _, graphToSymIter := range graphToSym {
			sym := graphToSymIter.Symbol

			http.HandleFunc("/"+sym, graphToSymIter.Handler)

			tradeSignal := graphToSymIter.TradeSignal
			//if sym == "GBP_USD" {
			if tradeSignal.Signal {
				fileName := SaveAnalyzedGraph(graphToSymIter.AnalyzedGraph, sym)
				image_upload.UploadImage(fileName)
				twilio.SendMsgFromData(fmt.Sprintf("%s - Crossed level %s", sym,
					strconv.FormatFloat(tradeSignal.LevelCrossed, 'f', -1, 64)), fileName)
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
	fileName := fileNamePrefix + "-image-" + curTime.Format(constants.IMAGE_FORMAT) + ".png"
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, image)
	return fileName
}
