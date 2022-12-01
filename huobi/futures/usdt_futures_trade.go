package futures

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"net/url"
	"time"
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

type usdtFuturesTrade struct {
	*USDTFutures
	apiOpts ApiOptions
}

func (f *usdtFuturesTrade) CreateOrder(pair CurrencyPair, qty, price float64, side OrderSide, orderTy OrderType, opts ...OptionParameter) (*Order, error) {
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
		return nil, err
	}

	logger.Debugf("[create order] response data=%s", string(data))

	ord, err := f.unmarshalerOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	ord.Status = OrderStatus_Pending

	return ord, nil
}

func (f *usdtFuturesTrade) GetOrderInfo(pair CurrencyPair, id string, opts ...OptionParameter) (*Order, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)

	if id != "" {
		params.Set("order_id", id)
	}

	MergeOptionParams(&params, opts...)

	data, err := f.DoAuthRequest(http.MethodPost, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetOrderUri), &params, nil)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[GetOrderInfo] %s", string(data))
	if data == nil || len(data) == 0 ||
		bytes.Compare(data, []byte{110, 117, 108, 108}) == 0 {
		return nil, nil
	}

	order, err := f.unmarshalerOpts.GetOrderInfoResponseUnmarshaler(data)
	if err != nil {
		return nil, err
	}
	order.Pair = pair
	return order, nil
}

func (f *usdtFuturesTrade) GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	params.Set("page_size", "50")
	data, err := f.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetPendingOrdersUri), &params, nil)
	if err != nil {
		return nil, err
	}
	logger.Debugf("[GetPendingOrders] %s", string(data))
	return f.unmarshalerOpts.GetPendingOrdersResponseUnmarshaler(data)
}

func (f *usdtFuturesTrade) GetHistoryOrders(pair CurrencyPair, opts ...OptionParameter) ([]Order, error) {
	params := url.Values{}
	params.Set("contract", pair.Symbol)
	params.Set("trade_type", "0")
	params.Set("type", "1")
	params.Set("status", "0")
	MergeOptionParams(&params, opts...)

	data, err := f.DoAuthRequest(http.MethodPost,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.GetHistoryOrdersUri), &params, nil)
	if err != nil {
		return nil, err
	}
	logger.Debugf("[GetHistoryOrders] %s", string(data))
	return f.unmarshalerOpts.GetHistoryOrdersResponseUnmarshaler(data)
}

func (f *usdtFuturesTrade) CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) error {
	params := url.Values{}
	params.Set("order_id", id)
	params.Set("contract_code", pair.Symbol)

	MergeOptionParams(&params, opt...)

	if params["client_order_id"] != nil {
		params.Del("order_id")
	}

	data, err := f.DoAuthRequest(http.MethodPost, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.CancelOrderUri), &params, nil)
	if err != nil {
		return err
	}

	return f.unmarshalerOpts.CancelOrderResponseUnmarshaler(data)
}

func (f *usdtFuturesTrade) CancelOrders(pair *CurrencyPair, id []string, opt ...OptionParameter) error {
	//TODO implement me
	panic("implement me")
}

func (f *usdtFuturesTrade) DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error) {
	///////////////////// 参数签名 ////////////////////////
	signParams := url.Values{}
	signParams.Set("AccessKeyId", f.apiOpts.Key)
	signParams.Set("SignatureMethod", "HmacSHA256")
	signParams.Set("SignatureVersion", "2")
	signParams.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))

	reqURL, _ := url.Parse(reqUrl)
	path := reqURL.RequestURI()
	domain := reqURL.Hostname()

	payload := fmt.Sprintf("%s\n%s\n%s\n%s", method, domain, path, signParams.Encode())
	sign, _ := HmacSHA256Base64Sign(f.apiOpts.Secret, payload)

	signParams.Set("Signature", sign)
	logger.Debugf("[DoAuthRequest] params=%s", signParams.Encode())
	///////////////////签名结束////////////////////

	if header == nil {
		header = make(map[string]string, 1)
	}
	header["Content-Type"] = "application/json"

	reqBody, _ := ValuesToJson(*params)
	logger.Debugf("request body: %s", string(reqBody))

	respBodyData, err := GetHttpCli().DoRequest(method, reqUrl+"?"+signParams.Encode(), string(reqBody), header)

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
