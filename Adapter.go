package goex

import "strings"

func AdaptTradeSide(side string) TradeSide {
	side2 := strings.ToUpper(side)
	switch side2 {
	case "SELL":
		return SELL
	case "BUY":
		return BUY
	case "BUY_MARKET":
		return BUY_MARKET
	case "SELL_MARKET":
		return SELL_MARKET
	default:
		return -1
	}
}

func AdaptKlinePeriodForOKEx(period int) string {
	switch period {
	case KLINE_PERIOD_1MIN:
		return "1min"
	case KLINE_PERIOD_5MIN:
		return "5min"
	case KLINE_PERIOD_15MIN:
		return "15min"
	case KLINE_PERIOD_30MIN:
		return "30min"
	case KLINE_PERIOD_1H:
		return "1hour"
	case KLINE_PERIOD_4H:
		return "4hour"
	case KLINE_PERIOD_1DAY:
		return "day"
	case KLINE_PERIOD_2H:
		return "2hour"
	case KLINE_PERIOD_1WEEK:
		return "week"
	default:
		return "day"
	}
}
