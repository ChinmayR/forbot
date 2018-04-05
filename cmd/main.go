package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/backtester"
	"github.com/ChinmayR/forbot/constants"
	"github.com/aws/aws-lambda-go/lambda"
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
	fmt.Println(runMain())
	fmt.Println("Ran main successfully")
	return Response{
		Message: fmt.Sprintf("Success for request Id %f", request.ID),
		Ok:      true,
	}, nil
}

type GraphToSym struct {
	Symbol        string
	GraphAnalysis forbot.GraphAnalysis
}

func runMain() []GraphToSym {
	var retVal []GraphToSym
	for _, sym := range constants.Symbols {
		graphAnalysis := forbot.GetGraphAnalysisForSymbol(
			sym,
			time.Now().AddDate(0, 0, -10),
			time.Now().Add(time.Minute*-1))
		retVal = append(retVal, GraphToSym{Symbol: sym, GraphAnalysis: graphAnalysis})
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

		symbol := constants.GBP_USD
		handler := backtester.RunBacktest(
			symbol,
			time.Now().AddDate(0, 0, -10),
			time.Now().Add(time.Minute*-1))

		http.HandleFunc("/"+symbol, handler)
		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":8080", nil))

	} else if executionType == NORMAL {

		graphToSym := runMain()
		for _, graphToSymIter := range graphToSym {
			sym := graphToSymIter.Symbol
			graphAnalysis := graphToSymIter.GraphAnalysis

			//graphAnalysis.SaveGraph(sym)
			http.HandleFunc("/"+sym, graphAnalysis.Handler)
		}

		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
