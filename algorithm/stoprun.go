package algorithm

import (
	"fmt"

	"github.com/ChinmayR/forbot"
	"github.com/ChinmayR/forbot/constants"
)

func RunStopRunAlgo(ga forbot.GraphAnalysis) bool {
	customLevels := constants.GetStopRunPointsForSymbol(ga.Symbol)
	fmt.Println(customLevels)
	startValidHour, endValidHour := constants.GetValidHourForSymbol(ga.Symbol)

	currentHour := getCurrentTime(ga.Xv).Hour()
	if currentHour > startValidHour && currentHour < endValidHour {
		// search the last 10 candles for the first candle pattern

		// if first candle pattern found then search again for second candle pattern

		// if third candle pattern found then send text
	}

	return false
}
