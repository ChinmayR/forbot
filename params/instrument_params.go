package params

import (
	"net/url"
	"time"

	"github.com/ChinmayR/forbot/constants"
)

// http://developer.oanda.com/rest-live-v20/instrument-ep/
type InstrumentCandlesParams struct {
	//Price       string
	Granularity string ////http://developer.oanda.com/rest-live-v20/instrument-df/#CandlestickGranularity
	//Count          int
	From time.Time
	To   time.Time
	//DailyAlignment int
}

func (p InstrumentCandlesParams) ToValues() url.Values {
	v := url.Values{}
	//v.Add("price", p.Price)
	v.Add("granularity", p.Granularity)
	//v.Add("count", strconv.Itoa(p.Count))
	v.Add("from", p.From.Format(constants.RFC3339))
	v.Add("to", p.To.Format(constants.RFC3339))
	//v.Add("dailyAlignment", strconv.Itoa(p.DailyAlignment))

	return v
}
