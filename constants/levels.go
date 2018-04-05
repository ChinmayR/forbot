package constants

const (
	ACCOUNT_ID = "509983"
	AUTH_TOKEN = "cb9b68b71444f1c24aa5c9e9e40ae409-3603616a28fd4e42580d23d90856ecd5"

	RFC3339      = "2006-01-02T15:04:05Z07:00"
	IMAGE_FORMAT = "01-02T15-04"
)

const (
	EUR_USD = "EUR_USD"
	GBP_USD = "GBP_USD"
	USD_JPY = "USD_JPY"
	EUR_JPY = "EUR_JPY"
)

var Symbols = []string{EUR_USD, GBP_USD, USD_JPY, EUR_JPY}

type HourRange struct {
	StartHour int
	EndHour   int
}

var ValidHours = map[string]HourRange{
	"EUR_USD": HourRange{
		StartHour: 0,
		EndHour:   24,
	},
	"GBP_USD": HourRange{
		StartHour: 0,
		EndHour:   24,
	},
	"USD_JPY": HourRange{
		StartHour: 0,
		EndHour:   24,
	},
	"EUR_JPY": HourRange{
		StartHour: 0,
		EndHour:   24,
	},
}

func GetValidHourForSymbol(symbol string) (int, int) {
	validHoursForSymbol := ValidHours[symbol]
	return validHoursForSymbol.StartHour, validHoursForSymbol.EndHour
}

type StopRunPoint struct {
	Symbol string
	Year   int
	Month  int
	Day    int
	Values []float64
}

var StopRunPoints_EUR_USD = []StopRunPoint{
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 15, Values: []float64{1.2521, 1.2388, 1.2303}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 16, Values: []float64{1.2521, 1.2458, 1.2388}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 19, Values: []float64{1.2456}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 20, Values: []float64{1.2435}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 21, Values: []float64{1.2372}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 22, Values: []float64{1.2372, 1.2435, 1.2212}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 23, Values: []float64{1.2319, 1.2297}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 26, Values: []float64{}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 27, Values: []float64{1.2354, 1.2435, 1.2260, 1.2212}},
	{Symbol: EUR_USD, Year: 2018, Month: 2, Day: 28, Values: []float64{1.2260, 1.2282}},
	{Symbol: EUR_USD, Year: 2018, Month: 3, Day: 29, Values: []float64{1.2355, 1.2375}},
}

var StopRunPoints_GBP_USD = []StopRunPoint{
	{Symbol: GBP_USD, Year: 2018, Month: 4, Day: 5, Values: []float64{1.4094, 1.4200, 1.4015}},
}

var StopRunPoints_USD_JPY = []StopRunPoint{}

var StopRunPoints_EUR_JPY = []StopRunPoint{
	{Symbol: EUR_JPY, Year: 2018, Month: 4, Day: 5, Values: []float64{131.77, 132.32, 130.87, 130.01}},
}

func GetStopRunPointsForSymbol(symbol string) []StopRunPoint {
	switch symbol {
	case EUR_USD:
		return StopRunPoints_EUR_USD
	case GBP_USD:
		return StopRunPoints_GBP_USD
	case USD_JPY:
		return StopRunPoints_USD_JPY
	case EUR_JPY:
		return StopRunPoints_EUR_JPY
	default:
		return StopRunPoints_EUR_USD
	}
}
