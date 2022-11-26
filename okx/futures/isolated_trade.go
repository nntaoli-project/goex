package futures

import (
	"errors"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type fIsolatedTrade struct {
	*common.Trade
}

func newFIsolatedTrade(apiOpts ...goex.ApiOption) *fIsolatedTrade {
	ft := new(fIsolatedTrade)
	ft.Trade = common.NewTrade(apiOpts...)
	return ft
}

func (f *fIsolatedTrade) CreateOrder(order goex.Order, opts ...goex.OptionParameter) (*goex.Order, error) {
	if order.Side != goex.Futures_OpenBuy &&
		order.Side != goex.Futures_OpenSell &&
		order.Side != goex.Futures_CloseBuy &&
		order.Side != goex.Futures_CloseSell {
		return nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts, goex.OptionParameter{
		Key:   "tdMode",
		Value: "isolated",
	})
	return f.Trade.CreateOrder(order, opts...)
}
