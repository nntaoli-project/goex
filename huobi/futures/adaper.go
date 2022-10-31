package futures

import . "github.com/nntaoli-project/goex/v2"

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

func AdaptKlinePeriod(period KlinePeriod) string {
	switch period {
	case Kline_1h:
		return "1hour"
	case Kline_4h:
		return "4hour"
	default:
		return string(period)
	}
}
