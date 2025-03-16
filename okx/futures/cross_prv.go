package futures

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2/model"
)

type CrossPrvApi struct {
	*PrvApi
}

func (f *CrossPrvApi) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, []byte, error) {
	if side != Futures_OpenBuy &&
		side != Futures_OpenSell &&
		side != Futures_CloseBuy &&
		side != Futures_CloseSell {
		return nil, nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts,
		OptionParameter{
			Key:   "tdMode",
			Value: "cross",
		})

	return f.Prv.CreateOrder(pair, qty, price, side, orderTy, opts...)
}
