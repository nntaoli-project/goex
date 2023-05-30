package futures

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	"github.com/nntaoli-project/goex/v2/options"
	"github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type PrvApi struct {
	*common.Prv
	Isolated *IsolatedPrvApi
	Cross    *CrossPrvApi
}

func NewPrvApi(v5 *common.OKxV5, apiOpts ...options.ApiOption) *PrvApi {
	prvApi := new(PrvApi)
	prvApi.Prv = v5.NewPrvApi(apiOpts...)

	prvApi.Isolated = new(IsolatedPrvApi)
	prvApi.Isolated.PrvApi = prvApi

	prvApi.Cross = new(CrossPrvApi)
	prvApi.Cross.PrvApi = prvApi

	return prvApi
}

func (prv *PrvApi) GetFuturesAccount(coin string) (map[string]model.FuturesAccount, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.OKxV5.UriOpts.Endpoint, prv.OKxV5.UriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	acc, err := prv.OKxV5.UnmarshalOpts.GetFuturesAccountResponseUnmarshaler(data)
	return acc, responseBody, err
}

func (prv *PrvApi) GetPositions(pair model.CurrencyPair, opts ...model.OptionParameter) ([]model.FuturesPosition, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.OKxV5.UriOpts.Endpoint, prv.OKxV5.UriOpts.GetPositionsUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	util.MergeOptionParams(&params, opts...)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	positions, err := prv.OKxV5.UnmarshalOpts.GetPositionsResponseUnmarshaler(data)
	return positions, responseBody, err
}

func (prv *PrvApi) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	opt = append(opt, model.OptionParameter{
		Key:   "instType",
		Value: "SWAP",
	})
	return prv.Prv.GetHistoryOrders(pair, opt...)
}
