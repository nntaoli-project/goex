package spot

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type PrvApi struct {
	*common.Prv
}

func (api *PrvApi) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, []byte, error) {
	//check params
	if Spot_Buy != side && side != Spot_Sell {
		return nil, nil, errors.New("spot order side is error")
	}

	opts = append(opts,
		OptionParameter{
			Key:   "tdMode",
			Value: "cash",
		})

	return api.Prv.CreateOrder(pair, qty, price, side, orderTy, opts...)
}

func (api *PrvApi) GetHistoryOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, []byte, error) {
	opt = append(opt, OptionParameter{
		Key:   "instType",
		Value: "SPOT",
	})
	return api.Prv.GetHistoryOrders(pair, opt...)
}
