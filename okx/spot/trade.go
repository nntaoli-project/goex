package spot

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	. "github.com/nntaoli-project/goex/v2/options"
)

type spotTradeImp struct {
	*common.Trade
}

func newSpotTradeImp(apiOpts ...ApiOption) *spotTradeImp {
	s := new(spotTradeImp)
	s.Trade = common.NewTrade(apiOpts...)
	return s
}

func (s *spotTradeImp) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, error) {
	//check params
	if Spot_Buy != side && side != Spot_Sell {
		return nil, errors.New("spot order side is error")
	}

	opts = append(opts,
		OptionParameter{
			Key:   "tdMode",
			Value: "cash",
		})

	return s.Trade.CreateOrder(Order{
		Pair:    pair,
		Price:   price,
		Qty:     qty,
		Side:    side,
		OrderTy: orderTy,
	}, opts...)
}
