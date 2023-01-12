package spot

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	"net/http"
	"net/url"
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

func (api *PrvApi) GetAccount(coin string) (map[string]Account, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", api.UriOpts.Endpoint, api.UriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, responseBody, err := api.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	acc, err := api.UnmarshalOpts.GetAccountResponseUnmarshaler(data)
	return acc, responseBody, err
}
