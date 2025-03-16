package common

import (
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
	case model.Kline_60min, model.Kline_1h:
		return "1H"
	case model.Kline_4h:
		return "4H"
	case model.Kline_6h:
		return "6H"
	case model.Kline_1day:
		return "1D"
	case model.Kline_1week:
		return "1W"
	default:
		return string(period)
	}
}

func adaptOrderSideToSym(s model.OrderSide) (side, posSide string) {
	switch s {
	case model.Spot_Buy:
		return "buy", ""
	case model.Spot_Sell:
		return "sell", ""
	case model.Futures_OpenBuy:
		return "buy", "long"
	case model.Futures_OpenSell:
		return "sell", "short"
	case model.Futures_CloseBuy:
		return "sell", "long"
	case model.Futures_CloseSell:
		return "buy", "short"
	}
	return "", ""
}

func adaptOrderTypeToSym(ty model.OrderType) string {
	switch ty {
	case model.OrderType_Limit:
		return "limit"
	case model.OrderType_Market:
		return "market"
	}
	return string(ty)
}

func adaptSymToOrderSide(side, posSide string) model.OrderSide {
	if side == "buy" {
		switch posSide {
		case "long":
			return model.Futures_OpenBuy
		case "short":
			return model.Futures_CloseSell
		default:
			return model.Spot_Buy
		}
	} else if side == "sell" {
		switch posSide { //现货
		case "long":
			return model.Futures_CloseBuy
		case "short":
			return model.Futures_OpenSell
		default:
			return model.Spot_Sell
		}
	}
	return model.OrderSide("unknown")
}

func adaptSymToOrderTy(st string) model.OrderType {
	switch st {
	case "limit":
		return model.OrderType_Limit
	case "market":
		return model.OrderType_Market
	default:
		return model.OrderType(st)
	}
}

func adaptSymToOrderStatus(st string) model.OrderStatus {
	switch st {
	case "live":
		return model.OrderStatus_Pending
	case "filled":
		return model.OrderStatus_Finished
	case "canceled":
		return model.OrderStatus_Canceled
	case "partially_filled":
		return model.OrderStatus_PartFinished
	default:
		return model.OrderStatus(-1)
	}
}

func AdaptQtyOrPricePrecision(sz string) int {
	if sz == "1" {
		return 0
	}
	return len(sz) - 2
}

func AdaptOrderClientIDOptionParameter(params *url.Values) {
	cid := params.Get(model.Order_Client_ID__Opt_Key)
	if cid != "" {
		params.Set("clOrdId", cid) //clOrdId
		params.Del(model.Order_Client_ID__Opt_Key)
	}
}
