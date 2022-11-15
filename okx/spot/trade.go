package spot

import (
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type spotTradeImp struct {
	*common.Trade
}

func newSpotTradeImp(apiOpts ...goex.ApiOption) *spotTradeImp {
	s := new(spotTradeImp)
	s.Trade = common.NewTrade(apiOpts...)
	return s
}

func (s *spotTradeImp) CreateOrder(order goex.Order, opts ...goex.OptionParameter) (*goex.Order, error) {
	opts = append(opts, goex.OptionParameter{
		Key:   "tdMode",
		Value: "cash",
	})
	return s.Trade.CreateOrder(order, opts...)
}
