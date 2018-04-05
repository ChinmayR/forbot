package algorithm

import (
	"time"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/constants"
)

func RunBasicAlgo(ga forbot.GraphAnalysis, curTime time.Time) bool {
	customLevels := constants.GetStopRunPointsForSymbol(ga.Symbol)
	startValidHour, endValidHour := constants.GetValidHourForSymbol(ga.Symbol)

	// For the curTime passed in, if its closest candle penetrates any of the custom
	// levels for that day, then return true
	currentHour := curTime.Hour()
	customPointsForCurTime := GetCustomPointForTime(customLevels, curTime)
	if currentHour >= startValidHour && currentHour <= endValidHour {
		for _, custLvl := range customPointsForCurTime.Values {
			closestXvIndex := GetClosestXvIndex(ga, curTime)
			if custLvl < ga.YvHigh[closestXvIndex] && custLvl > ga.YvLow[closestXvIndex] {
				return true
			}
		}
	}
	return false
}
