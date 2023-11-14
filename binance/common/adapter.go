package common

import (
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
	"net/url"
)

func AdaptKlinePeriodToSymbol(period model.KlinePeriod) string {
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

func AdaptOrderTypeToString(ty model.OrderType) string {
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

func AdaptOrderSideToString(s model.OrderSide) string {
	switch s {
	case model.Spot_Buy, model.Futures_OpenBuy, model.Futures_CloseSell:
		return "BUY"
	case model.Spot_Sell, model.Futures_OpenSell, model.Futures_CloseBuy:
		return "SELL"
	default:
		logger.Warnf("[adapt side] order side:%s error", s)
	}
	return string(s)
}

func AdaptStringToOrderStatus(st string) model.OrderStatus {
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

func AdaptStringToFuturesOrderSide(side, positionSide string) model.OrderSide {
	switch side {
	case "BUY":
		if positionSide == "LONG" {
			return model.Futures_OpenBuy
		}
		if positionSide == "SHORT" {
			return model.Futures_CloseSell
		}
	case "SELL":
		if positionSide == "LONG" {
			return model.Futures_CloseBuy
		}
		if positionSide == "SHORT" {
			return model.Futures_OpenSell
		}
	default:
		logger.Warnf("[adaptOrderOrigSide] unknown order origin side: %s", side)
	}
	return model.OrderSide(side)
}

func AdaptStringToOrderType(ty string) model.OrderType {
	switch ty {
	case "LIMIT":
		return model.OrderType_Limit
	case "MARKET":
		return model.OrderType_Market
	default:
		return model.OrderType(ty)
	}
}

func AdaptOrderClientIDOptionParameter(params *url.Values) {
	cid := params.Get(model.Order_Client_ID__Opt_Key)
	if cid != "" {
		params.Set("newClientOrderId", cid) //clOrdId
		params.Del(model.Order_Client_ID__Opt_Key)
	}
}
