package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type Trade struct {
	*V5
	apiOpts goex.ApiOptions
}

func (t *Trade) CreateOrder(order goex.Order, opts ...goex.OptionParameter) (*goex.Order, error) {
	reqUrl := fmt.Sprintf("%s%s", t.uriOpts.Endpoint, t.uriOpts.NewOrderUri)
	params := url.Values{}

	params.Set("instId", order.Pair.Symbol)
	//params.Set("tdMode", "cash")
	//params.Set("posSide", "")
	params.Set("ordType", adaptOrderTypeToSym(order.OrderTy))
	params.Set("px", util.FloatToString(order.Price, order.Pair.PricePrecision))
	params.Set("sz", util.FloatToString(order.Qty, order.Pair.QtyPrecision))

	side, posSide := adaptOrderSideToSym(order.Side)
	params.Set("side", side)
	if posSide != "" {
		params.Set("posSide", posSide)
	}

	if order.CId != "" {
		params.Set("clOrdId", order.CId)
	}
	util.MergeOptionParams(&params, opts...)

	data, err := t.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[CreateOrder] err=%s, response=%s", err.Error(), string(data))
		return nil, err
	}

	ord, err := t.unmarshalOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	ord.Pair = order.Pair
	ord.Price = order.Price
	ord.Qty = order.Qty
	ord.Side = order.Side
	ord.OrderTy = order.OrderTy
	ord.Status = goex.OrderStatus_Pending

	return ord, err
}

func (t *Trade) GetOrderInfo(pair goex.CurrencyPair, id string, opt ...goex.OptionParameter) (*goex.Order, error) {
	reqUrl := fmt.Sprintf("%s%s", t.uriOpts.Endpoint, t.uriOpts.GetOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)

	data, err := t.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, err
	}

	ord, err := t.unmarshalOpts.GetOrderInfoResponseUnmarshaler(data[1 : len(data)-1])
	if err != nil {
		return nil, err
	}

	ord.Pair = pair
	ord.Origin = data

	return ord, nil
}

func (t *Trade) GetPendingOrders(pair goex.CurrencyPair, opt ...goex.OptionParameter) ([]goex.Order, error) {
	reqUrl := fmt.Sprintf("%s%s", t.uriOpts.Endpoint, t.uriOpts.GetPendingOrdersUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	data, err := t.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, err
	}

	return t.unmarshalOpts.GetPendingOrdersResponseUnmarshaler(data)
}

func (t *Trade) GetHistoryOrders(pair goex.CurrencyPair, opt ...goex.OptionParameter) ([]goex.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Trade) CancelOrder(pair goex.CurrencyPair, id string, opt ...goex.OptionParameter) error {
	reqUrl := fmt.Sprintf("%s%s", t.uriOpts.Endpoint, t.uriOpts.CancelOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)
	util.MergeOptionParams(&params, opt...)

	data, err := t.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if data != nil && len(data) > 0 {
		return t.unmarshalOpts.CancelOrderResponseUnmarshaler(data)
	}

	return err
}

func (t *Trade) DoAuthRequest(httpMethod, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	return t.V5.DoAuthRequest(httpMethod, reqUrl, params, t.apiOpts, headers)
}

func NewTrade(opts ...goex.ApiOption) *Trade {
	var api = new(Trade)
	for _, opt := range opts {
		opt(&api.apiOpts)
	}
	return api
}
