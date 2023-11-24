package futures

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/huobi/common"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Ts      int64  `json:"ts"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type TradeBaseResponse struct {
	*BaseResponse
	Data json.RawMessage `json:"data"`
}

type USDTSwapPrvApi struct {
	*USDTSwap
	apiOpts options.ApiOptions
}

func NewUSDTSwapPrvApi(apiOpts ...options.ApiOption) *USDTSwapPrvApi {
	f := &USDTSwapPrvApi{}
	for _, opt := range apiOpts {
		opt(&f.apiOpts)
	}
	return f
}

func (f *USDTSwapPrvApi) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, []byte, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	params.Set("price", FloatToString(price, pair.PricePrecision))
	params.Set("volume", FloatToString(qty, pair.QtyPrecision))
	params.Set("order_price_type", string(orderTy))

	direction, offset := AdaptSideToDirectionAndOffset(side)
	params.Set("direction", direction)
	params.Set("offset", offset)

	MergeOptionParams(&params, opts...)

	if params.Get("lever_rate") == "" {
		logger.Warnf("[create order] set default lever rate 10")
		params.Set("lever_rate", "10") //set default 10 lever rate
	}

	data, err := f.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.NewOrderUri), &params, nil)
	if err != nil {
		return nil, data, err
	}

	logger.Debugf("[create order] response data=%s", string(data))

	ord, err := f.unmarshalerOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	ord.Status = OrderStatus_Pending

	return ord, data, nil
}

func (f *USDTSwapPrvApi) GetOrderInfo(pair CurrencyPair, id string, opts ...OptionParameter) (*Order, []byte, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)

	if id != "" {
		params.Set("order_id", id)
	}

	MergeOptionParams(&params, opts...)

	data, err := f.DoAuthRequest(http.MethodPost, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetOrderUri), &params, nil)
	if err != nil {
		return nil, data, err
	}

	logger.Debugf("[GetOrderInfo] %s", string(data))
	if data == nil || len(data) == 0 ||
		bytes.Equal(data, []byte{110, 117, 108, 108}) {
		return nil, data, nil
	}

	order, err := f.unmarshalerOpts.GetOrderInfoResponseUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}
	order.Pair = pair
	return order, data, nil
}

func (f *USDTSwapPrvApi) GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, []byte, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	params.Set("page_size", "50")
	data, err := f.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetPendingOrdersUri), &params, nil)
	if err != nil {
		return nil, data, err
	}
	logger.Debugf("[GetPendingOrders] %s", string(data))
	orders, err := f.unmarshalerOpts.GetPendingOrdersResponseUnmarshaler(data)
	return orders, data, err
}

func (f *USDTSwapPrvApi) GetHistoryOrders(pair CurrencyPair, opts ...OptionParameter) ([]Order, []byte, error) {
	params := url.Values{}
	params.Set("contract", pair.Symbol)
	params.Set("trade_type", "0")
	params.Set("type", "1")
	params.Set("status", "0")
	MergeOptionParams(&params, opts...)

	data, err := f.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetHistoryOrdersUri), &params, nil)
	if err != nil {
		return nil, data, err
	}
	logger.Debugf("[GetHistoryOrders] %s", string(data))
	orders, err := f.unmarshalerOpts.GetHistoryOrdersResponseUnmarshaler(data)
	return orders, data, err
}

func (f *USDTSwapPrvApi) CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) ([]byte, error) {
	params := url.Values{}
	params.Set("order_id", id)
	params.Set("contract_code", pair.Symbol)

	MergeOptionParams(&params, opt...)

	if params["client_order_id"] != nil {
		params.Del("order_id")
	}

	data, err := f.DoAuthRequest(http.MethodPost, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.CancelOrderUri), &params, nil)
	if err != nil {
		return data, err
	}

	return data, f.unmarshalerOpts.CancelOrderResponseUnmarshaler(data)
}

func (f *USDTSwapPrvApi) CancelOrders(pair *CurrencyPair, id []string, opt ...OptionParameter) error {
	//TODO implement me
	panic("implement me")
}

func (f *USDTSwapPrvApi) GetFuturesAccount(coin string) (acc map[string]FuturesAccount, responseBody []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *USDTSwapPrvApi) GetPositions(pair CurrencyPair, opts ...OptionParameter) (positions []FuturesPosition, responseBody []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *USDTSwapPrvApi) DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error) {
	///////////////////// 参数签名 ////////////////////////
	signParams := common.DoSignParam(method, reqUrl, f.apiOpts)

	if header == nil {
		header = make(map[string]string, 1)
	}
	header["Content-Type"] = "application/json"

	reqBody, _ := ValuesToJson(*params)
	logger.Debugf("request body: %s", string(reqBody))

	respBodyData, err := Cli.DoRequest(method, reqUrl+"?"+signParams.Encode(), string(reqBody), header)

	if err != nil {
		return nil, err
	}

	var baseResp TradeBaseResponse
	err = f.unmarshalerOpts.ResponseUnmarshaler(respBodyData, &baseResp)
	if err != nil {
		return nil, err
	}

	if baseResp.Status == "ok" || baseResp.Code == 200 {
		return baseResp.Data, nil
	}

	return nil, errors.New(string(respBodyData))
}
