package futures

import (
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
	opts = append(opts, goex.OptionParameter{
		Key:   "tdMode",
		Value: "isolated",
	})
	return f.Trade.CreateOrder(order, opts...)
}
