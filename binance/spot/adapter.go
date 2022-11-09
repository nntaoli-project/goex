package spot

import "github.com/nntaoli-project/goex/v2"

func adapterKlinePeriod(period goex.KlinePeriod) string {
	switch period {
	case goex.Kline_1min:
		return "1m"
	case goex.Kline_5min:
		return "5m"
	case goex.Kline_15min:
		return "15m"
	case goex.Kline_30min:
		return "30m"
	case goex.Kline_1h, goex.Kline_60min:
		return "1h"
	case goex.Kline_1day:
		return "1d"
	case goex.Kline_1week:
		return "1w"
	}
	return string(period)
}
