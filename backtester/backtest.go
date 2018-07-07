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

func GetAnalyzedGraphAndHandler(ga *forbot.GraphAnalysis, algoToRun utils.Algorithm) (HandlerFunc, chart.Chart, utils.TradeSignal) {
	//ga := forbot.GetGraphAnalysisForSymbol(symbol, from, to)
	graph := ga.GetGraph()

	// the last xV in the range below will set the trade signal so the signal
	// is sent for the "to" time (the last time in the "from" to "to" range)
	tradeSignal := algoToRun.RunAlgoFromTo(*ga, 0, len(ga.Xv))
	graph.Series = append(graph.Series, tradeSignal.SeriesToDraw...)

	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "image/png")
		graph.Render(chart.PNG, res)
	}, graph, tradeSignal
}
