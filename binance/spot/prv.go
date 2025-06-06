package spot

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/binance/common"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type PrvApi struct {
	*Spot
	apiOpts options.ApiOptions
}

func NewPrvApi(apiOpts ...options.ApiOption) *PrvApi {
	s := new(PrvApi)
	for _, opt := range apiOpts {
		opt(&s.apiOpts)
	}
	return s
}

func (s *PrvApi) GetAccount(coin string) (map[string]Account, []byte, error) {
	params := url.Values{}
	params.Set("omitZeroBalances", "true")
	reqUrl := fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.GetAccountUri)
	data, err := s.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, data, err
	}
	accounts, err := s.UnmarshalerOpts.GetAccountResponseUnmarshaler(data)
	return accounts, data, err
}

func (s *PrvApi) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opt ...OptionParameter) (*Order, []byte, error) {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	params.Set("side", adaptOrderSide(side))
	params.Set("type", adaptOrderType(orderTy))
	params.Set("timeInForce", "GTC")
	params.Set("quantity", FloatToString(qty, pair.QtyPrecision))
	params.Set("price", FloatToString(price, pair.PricePrecision))
	params.Set("newOrderRespType", "ACK")

	MergeOptionParams(&params, opt...)
	common.AdaptOrderClientIDOptionParameter(&params)

	data, err := s.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.NewOrderUri), &params, nil)
	if err != nil {
		return nil, data, err
	}

	ord, err := s.UnmarshalerOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	ord.Pair = pair
	ord.Price = price
	ord.Qty = qty
	ord.Status = OrderStatus_Pending
	ord.Side = side
	ord.OrderTy = orderTy

	return ord, data, nil
}

func (s *PrvApi) GetOrderInfo(pair CurrencyPair, id string, opt ...OptionParameter) (*Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.GetOrderUri)
	params := url.Values{}
	params.Set("symbol", pair.Symbol)

	if id != "" {
		params.Set("orderId", id)
	}

	MergeOptionParams(&params, opt...)
	adaptClientOrderId(&params)

	resp, err := s.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, resp, err
	}

	ord, err := s.UnmarshalerOpts.GetOrderInfoResponseUnmarshaler(resp)
	if ord != nil {
		ord.Pair = pair
	}

	return ord, resp, err
}

func (s *PrvApi) GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, []byte, error) {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	MergeOptionParams(&params, opt...)
	data, err := s.DoAuthRequest(http.MethodGet, fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.GetPendingOrdersUri), &params, nil)
	if err != nil {
		return nil, data, err
	}
	orders, err := s.UnmarshalerOpts.GetPendingOrdersResponseUnmarshaler(data)
	return orders, data, err
}

func (s *PrvApi) GetHistoryOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, []byte, error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)
	params.Set("limit", "100")
	MergeOptionParams(&params, opt...)
	reqUrl := fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.GetHistoryOrdersUri)
	data, err := s.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, data, err
	}
	orders, err := s.UnmarshalerOpts.GetHistoryOrdersResponseUnmarshaler(data)
	return orders, data, err
}

func (s *PrvApi) CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) ([]byte, error) {
	var params = url.Values{}
	params.Set("symbol", pair.Symbol)
	if id != "" {
		params.Set("orderId", id)
	}

	MergeOptionParams(&params, opt...)
	adaptClientOrderId(&params)

	data, err := s.DoAuthRequest(http.MethodDelete, fmt.Sprintf("%s%s", s.UriOpts.Endpoint, s.UriOpts.CancelOrderUri), &params, nil)
	if err != nil {
		return data, err
	}
	return data, s.UnmarshalerOpts.CancelOrderResponseUnmarshaler(data)
}

func (s *PrvApi) DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error) {
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
