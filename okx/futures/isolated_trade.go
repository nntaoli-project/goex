package futures

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type fIsolatedTrade struct {
	*common.Trade
}

func newFIsolatedTrade(apiOpts ...ApiOption) *fIsolatedTrade {
	ft := new(fIsolatedTrade)
	ft.Trade = common.NewTrade(apiOpts...)
	return ft
}

func (f *fIsolatedTrade) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, error) {
	if side != Futures_OpenBuy &&
		side != Futures_OpenSell &&
		side != Futures_CloseBuy &&
		side != Futures_CloseSell {
		return nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts,
		OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Trade.CreateOrder(Order{
		Pair:    pair,
		Qty:     qty,
		Price:   price,
		Side:    side,
		OrderTy: orderTy,
	}, opts...)
}
