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
