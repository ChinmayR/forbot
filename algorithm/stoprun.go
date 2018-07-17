package algorithm

import (
	"log"
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

func (ba StopRunAlgo) RunAlgo(ga forbot.GraphAnalysis, curTime time.Time) utils.TradeSignal {
	panic("UNSUPPORTED, NEED A RANGE FROM, TO")
}

func (ba StopRunAlgo) RunAlgoFromTo(ga forbot.GraphAnalysis, from, to int) utils.TradeSignal {
	customLevels := constants.GetStopRunPointsForSymbol(ga.Symbol)
	startValidHour, endValidHour := constants.GetValidHourForSymbol(ga.Symbol)

	firstCandleInfo, foundFirstCandle := &CandleInfo{CandleNumber: NONE}, false
	secondCandleInfo, foundSecondCandle := &CandleInfo{CandleNumber: NONE}, false
	thirdCandleInfo, foundThirdCandle := &CandleInfo{CandleNumber: NONE}, false
	var seriesToDraw = make([]chart.Series, 0)
	retTradeSignal := utils.TradeSignal{Signal: false, SeriesToDraw: seriesToDraw}
	sleepForCandles := 0

	for i := from; i < to; i += 1 {
		// set this to false so that this is only reflective of the last candle
		retTradeSignal.Signal = false

		if sleepForCandles > 0 {
			sleepForCandles--
			continue
		}
		xV := ga.Xv[i]
		currentHour := xV.Hour()
		customPointsForCurTime := utils.GetCustomPointForTime(customLevels, xV)
		if currentHour >= startValidHour && currentHour <= endValidHour {
			// 1) search the last 8 candles for the first candle pattern
			// 2) if first candle pattern found then search again for second candle pattern
			// 3) if third candle pattern found then send text

			if foundFirstCandle && foundSecondCandle && !foundThirdCandle {
				// if more than 8 candles have passed, then reset, else try to find second candle
				if i > (firstCandleInfo.CandleIndex + 8) {
					firstCandleInfo, foundFirstCandle = &CandleInfo{CandleNumber: NONE}, false
					secondCandleInfo, foundSecondCandle = &CandleInfo{CandleNumber: NONE}, false
					thirdCandleInfo, foundThirdCandle = &CandleInfo{CandleNumber: NONE}, false
				} else {
					thirdCandleInfo, foundThirdCandle = getThirdCandleInfo(ga, xV, *firstCandleInfo)
					if foundThirdCandle {
						log.Printf("Found third candle for %v with %+v\n", ga.Symbol, thirdCandleInfo)
						// Draw third candle
						seriesToDraw = append(seriesToDraw, draw.DrawCircleAtTime(ga, ga.Xv[thirdCandleInfo.CandleIndex], ga.YvClose[thirdCandleInfo.CandleIndex], drawing.ColorRed)...)

						retTradeSignal.Signal = true
						retTradeSignal.LevelCrossed = firstCandleInfo.LevelPierced

						// reset both first and second and third candle state
						firstCandleInfo, foundFirstCandle = &CandleInfo{CandleNumber: NONE}, false
						secondCandleInfo, foundSecondCandle = &CandleInfo{CandleNumber: NONE}, false
						thirdCandleInfo, foundThirdCandle = &CandleInfo{CandleNumber: NONE}, false
						sleepForCandles = 6 // sleep for 2 hours (if granularity is 15 mins candles)

						continue
					}
				}
			} else if foundFirstCandle && !foundSecondCandle {
				// if more than 8 candles have passed, then reset, else try to find second candle
				if i > (firstCandleInfo.CandleIndex + 8) {
					firstCandleInfo, foundFirstCandle = &CandleInfo{CandleNumber: NONE}, false
					secondCandleInfo, foundSecondCandle = &CandleInfo{CandleNumber: NONE}, false
					thirdCandleInfo, foundThirdCandle = &CandleInfo{CandleNumber: NONE}, false
				} else {
					secondCandleInfo, foundSecondCandle = getSecondCandleInfo(ga, xV, *firstCandleInfo)
					if foundSecondCandle {
						log.Printf("Found second candle for %v with %+v\n", ga.Symbol, secondCandleInfo)
						// Draw second candle
						seriesToDraw = append(seriesToDraw, draw.DrawCircleAtTime(ga, ga.Xv[secondCandleInfo.CandleIndex], ga.YvClose[secondCandleInfo.CandleIndex], drawing.ColorBlue)...)
						continue
					}
				}
			}
			if !foundFirstCandle {
				firstCandleInfo, foundFirstCandle = getFirstCandleInfo(ga, xV, customPointsForCurTime)
				if foundFirstCandle {
					log.Printf("Found first candle for %v with %+v\n", ga.Symbol, firstCandleInfo)
					// Draw first candle
					seriesToDraw = append(seriesToDraw, draw.DrawCircleAtTime(ga, ga.Xv[firstCandleInfo.CandleIndex], ga.YvClose[firstCandleInfo.CandleIndex], drawing.ColorGreen)...)
					continue
				}
			}
		}
	}

	retTradeSignal.SeriesToDraw = seriesToDraw
	return retTradeSignal
}

const (
	NONE = iota
	CANDLE_ONE
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
	closestXvIndex := utils.GetClosestXvIndex(ga.Xv, curTime)
	for _, custLvl := range customPointsForCurTime.Values {
		if custLvl < ga.YvHigh[closestXvIndex] && custLvl > ga.YvLow[closestXvIndex] {
			didPierceUp := false
			// check if pierced up or down
			if ga.YvOpen[closestXvIndex] < custLvl {
				didPierceUp = true
			}
			// check if the candle pierced by too much
			pipAmount := constants.ValidHours[ga.Symbol].Pip
			piercedByPips := math.Abs(ga.YvHigh[closestXvIndex] - custLvl)
			piercedByPips = piercedByPips / pipAmount
			if didPierceUp && piercedByPips > upperPipThresholdOnPierce {
				return nil, false
			}
			piercedByPips = math.Abs(ga.YvLow[closestXvIndex] - custLvl)
			piercedByPips = piercedByPips / pipAmount
			if !didPierceUp && piercedByPips > upperPipThresholdOnPierce {
				return nil, false
			}

			return &CandleInfo{
				CandleNumber: CANDLE_ONE,
				CandleIndex:  closestXvIndex,
				LevelPierced: custLvl,
				DidPierceUp:  didPierceUp,
			}, true
		}
	}

	return nil, false
}

func getSecondCandleInfo(ga forbot.GraphAnalysis, curTime time.Time, firstCandleInfo CandleInfo) (*CandleInfo, bool) {
	closestXvIndex := utils.GetClosestXvIndex(ga.Xv, curTime)
	if firstCandleInfo.DidPierceUp && ga.YvClose[closestXvIndex] > firstCandleInfo.LevelPierced {
		return nil, false
	}
	if !firstCandleInfo.DidPierceUp && ga.YvClose[closestXvIndex] < firstCandleInfo.LevelPierced {
		return nil, false
	}
	return &CandleInfo{
		CandleNumber: CANDLE_TWO,
		CandleIndex:  closestXvIndex,
		LevelPierced: firstCandleInfo.LevelPierced,
		DidPierceUp:  firstCandleInfo.DidPierceUp,
	}, true

	return nil, false
}

func getThirdCandleInfo(ga forbot.GraphAnalysis, curTime time.Time, firstCandleInfo CandleInfo) (*CandleInfo, bool) {
	closestXvIndex := utils.GetClosestXvIndex(ga.Xv, curTime)
	midPointOfFirstCandle := (ga.YvHigh[firstCandleInfo.CandleIndex] + ga.YvLow[firstCandleInfo.CandleIndex]) / float64(2)
	returnTrue := false
	if firstCandleInfo.DidPierceUp && ga.YvHigh[closestXvIndex] > midPointOfFirstCandle {
		returnTrue = true
	}
	if !firstCandleInfo.DidPierceUp && ga.YvLow[closestXvIndex] < midPointOfFirstCandle {
		returnTrue = true
	}
	if returnTrue {
		return &CandleInfo{
			CandleNumber: CANDLE_THREE,
			CandleIndex:  closestXvIndex,
			LevelPierced: firstCandleInfo.LevelPierced,
			DidPierceUp:  firstCandleInfo.DidPierceUp,
		}, true
	}

	return nil, false
}
