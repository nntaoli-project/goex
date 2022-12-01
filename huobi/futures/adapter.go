package futures

import (
	. "github.com/nntaoli-project/goex/v2/model"
)

func AdaptSideToDirectionAndOffset(side OrderSide) (direction, offset string) {
	switch side {
	case Futures_OpenBuy:
		return "buy", "open"
	case Futures_OpenSell:
		return "sell", "open"
	case Futures_CloseBuy:
		return "sell", "close"
	case Futures_CloseSell:
		return "buy", "close"
	default:
		panic("order side error")
	}
}

func AdaptOffsetDirectionToOrderSide(offset, direction string) OrderSide {
	if offset == "open" {
		if direction == "sell" {
			return Futures_OpenSell
		}
		return Futures_OpenBuy
	}

	if offset == "close" {
		if direction == "buy" {
			return Futures_CloseSell
		}
		return Futures_CloseBuy
	}

	return OrderSide("unknown")
}

func AdaptKlinePeriod(period KlinePeriod) string {
	switch period {
	case Kline_1h:
		return "60min"
	case Kline_4h:
		return "4hour"
	default:
		return string(period)
	}
}

func AdaptStatus(s int) OrderStatus {
	switch s {
	case 1, 2, 3:
		return OrderStatus_Pending
	case 4:
		return OrderStatus_PartFinished
	case 5, 6:
		return OrderStatus_Finished
	case 7:
		return OrderStatus_Canceled
	case 11:
		return OrderStatus_Canceling
	}
	return 1000
}
