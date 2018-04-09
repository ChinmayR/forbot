package algorithm

import (
	"fmt"
	"math"
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/algorithm/utils"
	"github.com/ChinmayR/forbot/backtester/draw"
	"github.com/ChinmayR/forbot/constants"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type StopRunAlgo struct{}

func (ba StopRunAlgo) RunAlgoFromTo(ga forbot.GraphAnalysis, from, to int) utils.TradeSignal {
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

func (sra StopRunAlgo) RunAlgo(ga forbot.GraphAnalysis, curTime time.Time) utils.TradeSignal {
	customLevels := constants.GetStopRunPointsForSymbol(ga.Symbol)
	startValidHour, endValidHour := constants.GetValidHourForSymbol(ga.Symbol)

	currentHour := curTime.Hour()
	customPointsForCurTime := utils.GetCustomPointForTime(customLevels, curTime)
	var seriesToDraw = make([]chart.Series, 0)
	if currentHour >= startValidHour && currentHour <= endValidHour {
		// 1) search the last 8 candles for the first candle pattern
		// 2) if first candle pattern found then search again for second candle pattern
		// 3) if third candle pattern found then send text

		firstCandleInfo, foundFirstCandle := getFirstCandleInfo(ga, curTime, customPointsForCurTime)
		if foundFirstCandle {
			fmt.Printf("Found first candle for %v with %+v\n", ga.Symbol, firstCandleInfo)
			// Draw first candle
			seriesToDraw = append(seriesToDraw, draw.DrawCircleAtTime(ga, ga.Xv[firstCandleInfo.CandleIndex], ga.YvClose[firstCandleInfo.CandleIndex], drawing.ColorRed)...)

			secondCandleInfo, foundSecondCandle := getSecondCandleInfo(ga, curTime, *firstCandleInfo)
			if foundSecondCandle {
				fmt.Printf("Found second candle for %v with %+v\n", ga.Symbol, secondCandleInfo)
				// Draw second candle
				seriesToDraw = append(seriesToDraw, draw.DrawCircleAtTime(ga, ga.Xv[secondCandleInfo.CandleIndex], ga.YvClose[secondCandleInfo.CandleIndex], drawing.ColorBlue)...)
				// ignore the third candle for now and just send the signal
				return utils.TradeSignal{
					Signal:       true,
					LevelCrossed: firstCandleInfo.LevelPierced,
					SeriesToDraw: seriesToDraw,
				}
			}
		}

	}

	return utils.TradeSignal{
		Signal:       false,
		SeriesToDraw: seriesToDraw,
	}
}

const (
	CANDLE_ONE = iota
	CANDLE_TWO
	CANDLE_THREE
)

type CandleType int

type CandleInfo struct {
	CandleNumber CandleType
	CandleIndex  int
	LevelPierced float64
	DidPierceUp  bool
}

const upperPipThresholdOnPierce = 15 //pips

func getFirstCandleInfo(ga forbot.GraphAnalysis, curTime time.Time, customPointsForCurTime constants.StopRunPoint) (*CandleInfo, bool) {
	candlesToLookBack := 8
	closestXvIndex := utils.GetClosestXvIndex(ga.Xv, curTime)
	for i := candlesToLookBack; i >= 0; i-- {
		indexToLookAt := closestXvIndex - i
		for _, custLvl := range customPointsForCurTime.Values {
			if custLvl < ga.YvHigh[indexToLookAt] && custLvl > ga.YvLow[indexToLookAt] {
				didPierceUp := false
				// check if pierced up or down
				if ga.YvOpen[indexToLookAt] < custLvl {
					didPierceUp = true
				}
				// check if the candle pierced by too much
				pipAmount := constants.ValidHours[ga.Symbol].Pip
				piercedByPips := math.Abs(ga.YvHigh[indexToLookAt] - custLvl)
				piercedByPips = piercedByPips / pipAmount
				if didPierceUp && piercedByPips > upperPipThresholdOnPierce {
					return nil, false
				}
				piercedByPips = math.Abs(ga.YvLow[indexToLookAt] - custLvl)
				piercedByPips = piercedByPips / pipAmount
				if !didPierceUp && piercedByPips > upperPipThresholdOnPierce {
					return nil, false
				}

				return &CandleInfo{
					CandleNumber: CANDLE_ONE,
					CandleIndex:  indexToLookAt,
					LevelPierced: custLvl,
					DidPierceUp:  didPierceUp,
				}, true
			}
		}
	}

	return nil, false
}

func getSecondCandleInfo(ga forbot.GraphAnalysis, curTime time.Time, firstCandleInfo CandleInfo) (*CandleInfo, bool) {
	indexToStartLooking := firstCandleInfo.CandleIndex + 1
	for i := indexToStartLooking; i >= 0; i-- {
		if firstCandleInfo.DidPierceUp && ga.YvClose[i] > firstCandleInfo.LevelPierced {
			return nil, false
		}
		if !firstCandleInfo.DidPierceUp && ga.YvClose[i] < firstCandleInfo.LevelPierced {
			return nil, false
		}
		return &CandleInfo{
			CandleNumber: CANDLE_TWO,
			CandleIndex:  i,
			LevelPierced: firstCandleInfo.LevelPierced,
			DidPierceUp:  firstCandleInfo.DidPierceUp,
		}, true
	}

	return nil, false
}
