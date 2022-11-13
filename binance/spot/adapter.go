package spot

import (
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
)

func adaptKlinePeriod(period goex.KlinePeriod) string {
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

func adaptOrderSide(s goex.OrderSide) string {
	switch s {
	case goex.Spot_Buy:
		return "BUY"
	case goex.Spot_Sell:
		return "SELL"
	default:
		logger.Warnf("[adapt side] order side:%+v error", s)
	}
	return s.String()
}

func adaptOrderType(ty goex.OrderType) string {
	switch ty {
	case goex.OrderType_Limit:
		return "LIMIT"
	case goex.OrderType_Market:
		return "MARKET"
	default:
		logger.Warnf("[adapt order type] order typ unknown")
	}
	return ty.String()
}

func adaptOrderStatus(st string) goex.OrderStatus {
	switch st {
	case "NEW":
		return goex.OrderStatus_Pending
	case "FILLED":
		return goex.OrderStatus_Finished
	case "CANCELED":
		return goex.OrderStatus_Canceled
	case "PARTIALLY_FILLED":
		return goex.OrderStatus_PartFinished
	}
	return goex.OrderStatus(-1)
}

func adaptOrderOrigSide(side string) goex.OrderSide {
	switch side {
	case "BUY":
		return goex.Spot_Buy
	case "SELL":
		return goex.Spot_Sell
	default:
		logger.Warnf("[adaptOrderOrigSide] unknown order origin side: %s", side)
	}
	return goex.OrderSide{
		Code: -1,
		Type: side,
	}
}

func adaptOrderOrigType(ty string) goex.OrderType {
	switch ty {
	case "LIMIT":
		return goex.OrderType_Limit
	case "MARKET":
		return goex.OrderType_Market
	default:
		return goex.OrderType{
			Code: -1,
			Type: ty,
		}
	}
}

func adaptOrderOrigStatus(st string) goex.OrderStatus {
	switch st {
	case "NEW":
		return goex.OrderStatus_Pending
	case "FILLED":
		return goex.OrderStatus_Finished
	case "CANCELED":
		return goex.OrderStatus_Canceled
	case "PARTIALLY_FILLED":
		return goex.OrderStatus_PartFinished
	default:
		return goex.OrderStatus(-1)
	}
}
