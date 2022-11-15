package futures

import (
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
	opts = append(opts, goex.OptionParameter{
		Key:   "tdMode",
		Value: "cross",
	})
	return f.Trade.CreateOrder(order, opts...)
}
