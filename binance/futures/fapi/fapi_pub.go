package fapi

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/nntaoli-project/goex/v2/binance/common"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

func (f *FApi) DoNoAuthRequest(httpMethod, reqUrl string, params *url.Values) ([]byte, []byte, error) {
	reqBody := ""
	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	responseBody, err := Cli.DoRequest(httpMethod, reqUrl, reqBody, nil)
	if err != nil {
		logger.Errorf("[DoNoAuthRequest] http request error, body: %s", string(responseBody))
		return responseBody, responseBody, err
	}

	if err = checkStatusError(responseBody); err != nil {
		return responseBody, responseBody, err
	}

	return responseBody, responseBody, nil
}

// checkStatusError 币安fapi部分接口参数校验失败时返回HTTP 200 + 错误信息, 只能从响应体判断:
//
//	{"status":"ERROR","type":"GENERAL","code":"99099990","errorData":"illegal params.","data":null,...}
//
// 这些接口成功响应的顶层没有status字段, 所以按等值判断不会误伤。
func checkStatusError(data []byte) error {
	if len(data) == 0 || data[0] != '{' { //数组等非错误响应, 交给各自的unmarshaler处理
		return nil
	}

	status, err := jsonparser.GetString(data, "status")
	if err != nil || status != "ERROR" {
		return nil
	}

	code, _ := jsonparser.GetString(data, "code")
	msg, _ := jsonparser.GetString(data, "errorData")
	if msg == "" {
		msg = string(data)
	}

	return fmt.Errorf("binance fapi error: code=%s, msg=%s", code, msg)
}

func (f *FApi) GetName() string {
	return "binance.com"
}

func (f *FApi) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	data, body, err := f.DoNoAuthRequest(http.MethodGet, f.UriOpts.Endpoint+f.UriOpts.GetExchangeInfoUri, &url.Values{})
	if err != nil {
		logger.Errorf("[GetExchangeInfo] http request error, body: %s", string(body))
		return nil, body, err
	}

	m, err := f.UnmarshalOpts.GetExchangeInfoResponseUnmarshaler(data)
	if err != nil {
		logger.Errorf("[GetExchangeInfo] unmarshaler data error, err: %s", err.Error())
		return nil, body, err
	}

	f.currencyPairM = m

	return m, body, err
}

func (f *FApi) NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error) {
	var (
		contractAlias string
		currencyPair  model.CurrencyPair
	)

	if len(opts) == 0 {
		contractAlias = "PERPETUAL"
	} else if opts[0].Key == "contractAlias" {
		contractAlias = opts[0].Value
	}

	currencyPair = f.currencyPairM[baseSym+quoteSym+contractAlias]
	if currencyPair.Symbol == "" {
		return currencyPair, errors.New("not found currency pair")
	}

	return currencyPair, nil
}

func (f *FApi) GetDepth(pair model.CurrencyPair, limit int, opt ...model.OptionParameter) (depth *model.Depth, responseBody []byte, err error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)
	params.Set("limit", fmt.Sprint(limit))

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := f.DoNoAuthRequest(http.MethodGet, f.UriOpts.Endpoint+f.UriOpts.DepthUri, &params)
	if err != nil {
		return nil, responseBody, err
	}

	dep, err := f.UnmarshalOpts.DepthUnmarshaler(data)
	dep.Pair = pair

	return dep, responseBody, err
}

func (f *FApi) GetTicker(pair model.CurrencyPair, opt ...model.OptionParameter) (ticker *model.Ticker, responseBody []byte, err error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)

	data, responseBody, err := f.DoNoAuthRequest(http.MethodGet, f.UriOpts.Endpoint+f.UriOpts.TickerUri, &params)
	if err != nil {
		return nil, responseBody, err
	}

	ticker, err = f.UnmarshalOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, responseBody, err
	}

	ticker.Pair = pair

	return ticker, responseBody, err
}

func (f *FApi) GetFundingRateHistory(pair model.CurrencyPair, limit int, opt ...model.OptionParameter) (rates []model.FundingRate, responseBody []byte, err error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)
	if limit > 0 {
		params.Set("limit", fmt.Sprint(limit))
	}

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := f.DoNoAuthRequest(http.MethodGet, f.UriOpts.Endpoint+f.UriOpts.GetFundingRateHistoryUri, &params)
	if err != nil {
		return nil, responseBody, err
	}

	rates, err = f.UnmarshalOpts.GetFundingRateHistoryResponseUnmarshaler(data)
	if err != nil {
		return nil, responseBody, err
	}

	for i := range rates {
		rates[i].Symbol = pair.Symbol
	}

	return rates, responseBody, nil
}

func (f *FApi) GetKline(pair model.CurrencyPair, period model.KlinePeriod, opt ...model.OptionParameter) (klines []model.Kline, responseBody []byte, err error) {
	var param = url.Values{}
	param.Set("symbol", pair.Symbol)
	param.Set("interval", common.AdaptKlinePeriodToSymbol(period))
	param.Set("limit", "100")

	util.MergeOptionParams(&param, opt...)

	data, responseBody, err := f.DoNoAuthRequest(http.MethodGet, f.UriOpts.Endpoint+f.UriOpts.KlineUri, &param)
	if err != nil {
		return nil, responseBody, err
	}

	klines, err = f.UnmarshalOpts.KlineUnmarshaler(data)

	for i, _ := range klines {
		klines[i].Pair = pair
	}

	return klines, responseBody, err
}
