package spot

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/binance/common"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

func (s *spotImpl) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opt ...OptionParameter) (*Order, error) {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	params.Set("side", adaptOrderSide(side))
	params.Set("type", adaptOrderType(orderTy))
	params.Set("timeInForce", "GTC")
	params.Set("quantity", FloatToString(qty, pair.QtyPrecision))
	params.Set("price", FloatToString(price, pair.PricePrecision))
	params.Set("newOrderRespType", "ACK")

	MergeOptionParams(&params, opt...)

	data, err := s.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.NewOrderUri), &params, nil)
	if err != nil {
		return nil, err
	}

	ord, err := s.unmarshalerOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	ord.Pair = order.Pair
	ord.Price = order.Price
	ord.Qty = order.Qty
	ord.Status = OrderStatus_Pending
	ord.Side = order.Side
	ord.OrderTy = order.OrderTy
	ord.Origin = data

	return ord, nil
}

func (s *spotImpl) GetOrderInfo(pair CurrencyPair, id string, opt ...OptionParameter) (*Order, error) {
	panic("")
}

func (s *spotImpl) GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, error) {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	MergeOptionParams(&params, opt...)
	data, err := s.DoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.GetPendingOrdersUri),
		&params, nil)
	if err != nil {
		return nil, err
	}
	return s.unmarshalerOpts.GetPendingOrdersResponseUnmarshaler(data)
}

func (s *spotImpl) GetHistoryOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spotImpl) CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) error {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	if id != "" {
		params.Set("orderId", id)
	}
	MergeOptionParams(&params, opt...)
	data, err := s.DoAuthRequest(http.MethodDelete, fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.CancelOrderUri), &params, nil)
	if err != nil {
		return err
	}
	return s.unmarshalerOpts.CancelOrderResponseUnmarshaler(data)
}

func (s *spotImpl) DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error) {
	if header == nil {
		header = make(map[string]string, 2)
	}
	header["X-MBX-APIKEY"] = s.apiOpts.Key
	common.SignParams(params, s.apiOpts.Secret)
	//if http.MethodGet == method {
	reqUrl += "?" + params.Encode()
	//}
	respBody, err := Cli.DoRequest(method, reqUrl, "", header)
	logger.Debugf("[DoAuthRequest] response body: %s", string(respBody))
	return respBody, err
}
