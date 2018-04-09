package algorithm

import (
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm/utils"
	"github.com/ChinmayR/forbot/backtester/draw"
	"github.com/ChinmayR/forbot/constants"
	"github.com/wcharczuk/go-chart/drawing"
)

type BasicAlgo struct{}

func (ba BasicAlgo) RunAlgoFromTo(ga forbot.GraphAnalysis, from, to int) utils.TradeSignal {
	retTradeSignal := utils.TradeSignal{}
	for i := from; i < to; i += 1 {
		xV := ga.Xv[i]
		tradeSignal := ba.RunAlgo(ga, xV)
		if tradeSignal.Signal {
			retTradeSignal.Signal = tradeSignal.Signal
			retTradeSignal.LevelCrossed = tradeSignal.LevelCrossed
			retTradeSignal.SeriesToDraw = append(retTradeSignal.SeriesToDraw, tradeSignal.SeriesToDraw...)
		}
	}
	return retTradeSignal
}

func (ba BasicAlgo) RunAlgo(ga forbot.GraphAnalysis, curTime time.Time) utils.TradeSignal {
	customLevels := constants.GetStopRunPointsForSymbol(ga.Symbol)
	startValidHour, endValidHour := constants.GetValidHourForSymbol(ga.Symbol)

	// For the curTime passed in, if its closest candle penetrates any of the custom
	// levels for that day, then return true
	currentHour := curTime.Hour()
	customPointsForCurTime := utils.GetCustomPointForTime(customLevels, curTime)
	if currentHour >= startValidHour && currentHour <= endValidHour {
		for _, custLvl := range customPointsForCurTime.Values {
			closestXvIndex := utils.GetClosestXvIndex(ga.Xv, curTime)
			if custLvl < ga.YvHigh[closestXvIndex] && custLvl > ga.YvLow[closestXvIndex] {
				return utils.TradeSignal{
					Signal:       true,
					LevelCrossed: custLvl,
					SeriesToDraw: draw.DrawCircleAtTime(ga, curTime, ga.YvClose[closestXvIndex], drawing.ColorRed),
				}
			}
		}
	}
	return utils.TradeSignal{
		Signal: false,
	}
}
