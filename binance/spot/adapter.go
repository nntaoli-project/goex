package spot

import (
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
)

func adaptKlinePeriod(period model.KlinePeriod) string {
	switch period {
	case model.Kline_1min:
		return "1m"
	case model.Kline_5min:
		return "5m"
	case model.Kline_15min:
		return "15m"
	case model.Kline_30min:
		return "30m"
	case model.Kline_1h, model.Kline_60min:
		return "1h"
	case model.Kline_1day:
		return "1d"
	case model.Kline_1week:
		return "1w"
	}
	return string(period)
}

func adaptOrderSide(s model.OrderSide) string {
	switch s {
	case model.Spot_Buy:
		return "BUY"
	case model.Spot_Sell:
		return "SELL"
	default:
		logger.Warnf("[adapt side] order side:%s error", s)
	}
	return string(s)
}

func adaptOrderType(ty model.OrderType) string {
	switch ty {
	case model.OrderType_Limit:
		return "LIMIT"
	case model.OrderType_Market:
		return "MARKET"
	default:
		logger.Warnf("[adapt order type] order typ unknown")
	}
	return string(ty)
}

func adaptOrderStatus(st string) model.OrderStatus {
	switch st {
	case "NEW":
		return model.OrderStatus_Pending
	case "FILLED":
		return model.OrderStatus_Finished
	case "CANCELED":
		return model.OrderStatus_Canceled
	case "PARTIALLY_FILLED":
		return model.OrderStatus_PartFinished
	}
	return model.OrderStatus(-1)
}

func adaptOrderOrigSide(side string) model.OrderSide {
	switch side {
	case "BUY":
		return model.Spot_Buy
	case "SELL":
		return model.Spot_Sell
	default:
		logger.Warnf("[adaptOrderOrigSide] unknown order origin side: %s", side)
	}
	return model.OrderSide(side)
}

func adaptOrderOrigType(ty string) model.OrderType {
	switch ty {
	case "LIMIT":
		return model.OrderType_Limit
	case "MARKET":
		return model.OrderType_Market
	default:
		return model.OrderType(ty)
	}
}

func adaptOrderOrigStatus(st string) model.OrderStatus {
	switch st {
	case "NEW":
		return model.OrderStatus_Pending
	case "FILLED":
		return model.OrderStatus_Finished
	case "CANCELED":
		return model.OrderStatus_Canceled
	case "PARTIALLY_FILLED":
		return model.OrderStatus_PartFinished
	default:
		return model.OrderStatus(-1)
	}
}
