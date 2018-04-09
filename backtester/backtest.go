package backtester

import (
	"net/http"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm/utils"
	"github.com/wcharczuk/go-chart"
)

type HandlerFunc func(res http.ResponseWriter, req *http.Request)

func RunBacktest(symbol string, from time.Time, to time.Time) {

}

func GetAnalyzedGraphAndHandler(symbol string, from time.Time, to time.Time, algoToRun utils.Algorithm) (HandlerFunc, chart.Chart, utils.TradeSignal) {
	ga := forbot.GetGraphAnalysisForSymbol(symbol, from, to)
	graph := ga.GetGraph()

	// the last xV in the range below will set the trade signal so the signal
	// is sent for the "to" time (the last time in the "from" to "to" range)
	//
	//var tradeSignal utils.TradeSignal
	////for i, xV := range ga.Xv {
	//for i := 0; i < len(ga.Xv); i += 1 {
	//	xV := ga.Xv[i]
	//	//if i != len(ga.Xv)-1 {
	//	//	continue
	//	//}
	//	tradeSignal = algoToRun.RunAlgo(ga, xV)
	//	if tradeSignal.Signal {
	//		//graph.Series = append(graph.Series, draw.DrawCircleAtTime(ga, xV, ga.YvClose[i], drawing.ColorRed)...)
	//		graph.Series = append(graph.Series, tradeSignal.SeriesToDraw...)
	//	}
	//}

	tradeSignal := algoToRun.RunAlgoFromTo(ga, 0, len(ga.Xv))
	graph.Series = append(graph.Series, tradeSignal.SeriesToDraw...)

	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "image/png")
		graph.Render(chart.PNG, res)
	}, graph, tradeSignal
}
