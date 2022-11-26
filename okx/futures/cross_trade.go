package futures

import (
	"errors"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type fCrossTrade struct {
	*common.Trade
}

func newFCrossTrade(apiOpts ...goex.ApiOption) *fCrossTrade {
	ft := new(fCrossTrade)
	ft.Trade = common.NewTrade(apiOpts...)
	return ft
}

func (f *fCrossTrade) CreateOrder(order goex.Order, opts ...goex.OptionParameter) (*goex.Order, error) {
	if order.Side != goex.Futures_OpenBuy &&
		order.Side != goex.Futures_OpenSell &&
		order.Side != goex.Futures_CloseBuy &&
		order.Side != goex.Futures_CloseSell {
		return nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}
	opts = append(opts, goex.OptionParameter{
		Key:   "tdMode",
		Value: "cross",
	})
	return f.Trade.CreateOrder(order, opts...)
}
