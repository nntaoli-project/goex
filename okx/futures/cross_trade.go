package futures

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	. "github.com/nntaoli-project/goex/v2/options"
)

type fCrossTrade struct {
	*common.Trade
}

func newFCrossTrade(apiOpts ...ApiOption) *fCrossTrade {
	ft := new(fCrossTrade)
	ft.Trade = common.NewTrade(apiOpts...)
	return ft
}

func (f *fCrossTrade) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, error) {
	if side != Futures_OpenBuy &&
		side != Futures_OpenSell &&
		side != Futures_CloseBuy &&
		side != Futures_CloseSell {
		return nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts,
		OptionParameter{
			Key:   "tdMode",
			Value: "cross",
		})

	return f.Trade.CreateOrder(Order{
		Pair:    pair,
		Qty:     qty,
		Price:   price,
		Side:    side,
		OrderTy: orderTy,
	}, opts...)
}
