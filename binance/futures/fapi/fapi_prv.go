package fapi

import (
	"errors"
	"github.com/nntaoli-project/goex/v2/binance/common"
	"github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
	"github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type Prv struct {
	*FApi
	apiOpts options.ApiOptions
}

func (p *Prv) GetAccount(currency string) (map[string]Account, []byte, error) {
	param := &url.Values{}
	responseBody, err := p.DoAuthRequest(http.MethodGet, p.UriOpts.Endpoint+p.UriOpts.GetAccountUri, param, nil)
	if err != nil {
		return nil, responseBody, err
	}
	accounts, err := p.UnmarshalOpts.GetAccountResponseUnmarshaler(responseBody)
	return accounts, responseBody, err
}

func (p *Prv) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opt ...OptionParameter) (order *Order, responseBody []byte, err error) {
	if orderTy == OrderType_Limit && qty*price < 5.0 { //币安规则
		return nil, nil, errors.New("MIN NOTIONAL must >= 5.0 USDT")
	}

	var param = url.Values{}
	param.Set("symbol", pair.Symbol)
	param.Set("price", util.FloatToString(price, pair.PricePrecision))
	param.Set("quantity", util.FloatToString(qty, pair.QtyPrecision))
	param.Set("type", common.AdaptOrderTypeToString(orderTy))
	param.Set("side", common.AdaptOrderSideToString(side))
	param.Set("timeInForce", "GTC")
	param.Set("newOrderRespType", "ACK")

	switch side {
	case Futures_OpenSell, Futures_CloseSell:
		param.Set("positionSide", "SHORT")
	case Futures_OpenBuy, Futures_CloseBuy:
		param.Set("positionSide", "LONG")
	}

	util.MergeOptionParams(&param, opt...)           //合并参数
	common.AdaptOrderClientIDOptionParameter(&param) //client id

	responseBody, err = p.DoAuthRequest(http.MethodPost, p.UriOpts.Endpoint+p.UriOpts.NewOrderUri, &param, nil)
	if err != nil {
		return nil, responseBody, err
	}

	ord, err := p.UnmarshalOpts.CreateOrderResponseUnmarshaler(responseBody)
	if ord != nil {
		ord.Pair = pair
		ord.Price = price
		ord.Qty = qty
		ord.Side = side
		ord.OrderTy = orderTy
	}

	return ord, responseBody, err
}

func (p *Prv) GetOrderInfo(pair CurrencyPair, id string, opt ...OptionParameter) (order *Order, responseBody []byte, err error) {
	param := &url.Values{}
	param.Set("symbol", pair.Symbol)
	param.Set("orderId", id)

	util.MergeOptionParams(param, opt...)

	data, err := p.DoAuthRequest(http.MethodGet, p.UriOpts.Endpoint+p.UriOpts.GetOrderUri, param, nil)
	if err != nil {
		return nil, data, err
	}

	order, err = p.UnmarshalOpts.GetOrderInfoResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	order.Pair = pair

	return
}

func (p *Prv) GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) (orders []Order, responseBody []byte, err error) {
	param := &url.Values{}
	param.Set("symbol", pair.Symbol)

	util.MergeOptionParams(param, opt...)

	data, err := p.DoAuthRequest(http.MethodGet, p.UriOpts.Endpoint+p.UriOpts.GetPendingOrdersUri, param, nil)
	if err != nil {
		return nil, data, err
	}

	orders, err = p.UnmarshalOpts.GetPendingOrdersResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	for i, _ := range orders {
		orders[i].Pair = pair
	}

	return orders, data, nil
}

func (p *Prv) GetHistoryOrders(pair CurrencyPair, opt ...OptionParameter) (orders []Order, responseBody []byte, err error) {
	param := &url.Values{}
	param.Set("symbol", pair.Symbol)
	param.Set("limit", "500")

	util.MergeOptionParams(param, opt...)

	data, err := p.DoAuthRequest(http.MethodGet, p.UriOpts.Endpoint+p.UriOpts.GetHistoryOrdersUri, param, nil)
	if err != nil {
		return nil, data, err
	}

	orders, err = p.UnmarshalOpts.GetHistoryOrdersResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	for i, _ := range orders {
		orders[i].Pair = pair
	}

	return orders, data, nil
}

func (p *Prv) CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) (responseBody []byte, err error) {
	param := &url.Values{}
	param.Set("symbol", pair.Symbol)
	param.Set("orderId", id)

	util.MergeOptionParams(param, opt...)

	data, err := p.DoAuthRequest(http.MethodDelete, p.UriOpts.Endpoint+p.UriOpts.CancelOrderUri, param, nil)
	if err != nil {
		return data, err
	}

	err = p.UnmarshalOpts.CancelOrderResponseUnmarshaler(data)

	return data, err
}

func (p *Prv) GetFuturesAccount(currency string) (acc map[string]FuturesAccount, responseBody []byte, err error) {
	panic("not implement")
}

func (p *Prv) GetPositions(pair CurrencyPair, opts ...OptionParameter) (positions []FuturesPosition, responseBody []byte, err error) {
	param := &url.Values{}
	param.Set("symbol", pair.Symbol)

	util.MergeOptionParams(param, opts...)

	data, err := p.DoAuthRequest(http.MethodGet, p.UriOpts.Endpoint+p.UriOpts.GetPositionsUri, param, nil)
	if err != nil {
		return nil, data, err
	}

	pos, err := p.UnmarshalOpts.GetPositionsResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	for i, _ := range pos {
		pos[i].Pair = pair
	}

	return pos, data, nil
}

func (p *Prv) DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error) {
	if header == nil {
		header = make(map[string]string, 2)
	}
	header["X-MBX-APIKEY"] = p.apiOpts.Key
	common.SignParams(params, p.apiOpts.Secret)
	//if http.MethodGet == method {
	reqUrl += "?" + params.Encode()
	//}
	respBody, err := httpcli.Cli.DoRequest(method, reqUrl, "", header)
	logger.Debugf("[DoAuthRequest] response body: %s", string(respBody))
	return respBody, err
}

func NewPrvApi(fapi *FApi, opts ...options.ApiOption) *Prv {
	var prv = new(Prv)
	prv.FApi = fapi
	for _, opt := range opts {
		opt(&prv.apiOpts)
	}
	return prv
}
