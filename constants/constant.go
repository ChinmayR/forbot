package constants

const (
	ACCOUNT_ID = "509983"
	AUTH_TOKEN = "91e0ecb7a2d464feb06769ee342b14b0-d8b4f92f7590add749377505842e6a69"

	RFC3339      = "2006-01-02T15:04:05Z07:00"
	IMAGE_FORMAT = "01-02T15-04"
)

type StopRunPoint struct {
	Year   int
	Month  int
	Day    int
	Values []float64
}

var StopRunPoints = []StopRunPoint{
	{Year: 2018, Month: 2, Day: 15, Values: []float64{1.2521, 1.2388, 1.2303}},
	{Year: 2018, Month: 2, Day: 16, Values: []float64{1.2521, 1.2458, 1.2388}},
	{Year: 2018, Month: 2, Day: 19, Values: []float64{1.2456}},
	{Year: 2018, Month: 2, Day: 20, Values: []float64{1.2435}},
	{Year: 2018, Month: 2, Day: 21, Values: []float64{1.2372}},
	{Year: 2018, Month: 2, Day: 22, Values: []float64{1.2372, 1.2435, 1.2212}},
	{Year: 2018, Month: 2, Day: 23, Values: []float64{1.2319, 1.2297}},
	{Year: 2018, Month: 2, Day: 26, Values: []float64{}},
	{Year: 2018, Month: 2, Day: 27, Values: []float64{1.2354, 1.2435, 1.2260, 1.2212}},
	{Year: 2018, Month: 2, Day: 28, Values: []float64{1.2260, 1.2282}},
}
