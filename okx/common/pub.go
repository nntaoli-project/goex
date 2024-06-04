package common

import (
	"fmt"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

func (okx *OKxV5) GetName() string {
	return "okx.com"
}

func (okx *OKxV5) GetDepth(pair CurrencyPair, size int, opt ...OptionParameter) (*Depth, []byte, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("sz", fmt.Sprint(size))
	MergeOptionParams(&params, opt...)

	data, responseBody, err := okx.DoNoAuthRequest("GET", okx.UriOpts.Endpoint+okx.UriOpts.DepthUri, &params)
	if err != nil {
		return nil, responseBody, err
	}

	dep, err := okx.UnmarshalOpts.DepthUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	dep.Pair = pair

	return dep, responseBody, err
}

func (okx *OKxV5) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, []byte, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	data, responseBody, err := okx.DoNoAuthRequest("GET", okx.UriOpts.Endpoint+okx.UriOpts.TickerUri, &params)
	if err != nil {
		return nil, data, err
	}

	tk, err := okx.UnmarshalOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, nil, err
	}

	tk.Pair = pair

	return tk, responseBody, err
}

func (okx *OKxV5) GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.KlineUri)
	param := url.Values{}
	param.Set("instId", pair.Symbol)
	param.Set("bar", AdaptKlinePeriodToSymbol(period))
	param.Set("limit", "100")
	MergeOptionParams(&param, opt...)

	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, nil, err
	}
	klines, err := okx.UnmarshalOpts.KlineUnmarshaler(data)
	return klines, responseBody, err
}

func (okx *OKxV5) GetExchangeInfo(instType string, opt ...OptionParameter) (map[string]CurrencyPair, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.GetExchangeInfoUri)
	param := url.Values{}
	param.Set("instType", instType)
	MergeOptionParams(&param, opt...)

	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, responseBody, err
	}

	currencyPairMap, err := okx.UnmarshalOpts.GetExchangeInfoResponseUnmarshaler(data)

	return currencyPairMap, responseBody, err
}

func (okx *OKxV5) GetFundingRate(pair CurrencyPair, opts ...OptionParameter) (rate *FundingRate, responseBody []byte, err error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.GetFundingRateUri)
	param := url.Values{}
	param.Set("instId", pair.Symbol)
	MergeOptionParams(&param, opts...)
	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, responseBody, err
	}
	rate, err = okx.UnmarshalOpts.GetFundingRateResponseUnmarshaler(data)
	if rate != nil && err == nil {
		rate.Symbol = pair.Symbol
	}
	return rate, nil, err
}

func (okx *OKxV5) GetFundingRateHistory(pair CurrencyPair, limit int, opts ...OptionParameter) (rates []FundingRate, responseBody []byte, err error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.GetFundingRateHistoryUri)
	param := url.Values{}
	param.Set("instId", pair.Symbol)
	param.Set("limit", fmt.Sprint(limit))
	MergeOptionParams(&param, opts...)
	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, responseBody, err
	}
	rates, err = okx.UnmarshalOpts.GetFundingRateHistoryResponseUnmarshaler(data)
	return rates, nil, err
}

func (okx *OKxV5) DoNoAuthRequest(httpMethod, reqUrl string, params *url.Values) ([]byte, []byte, error) {
	reqBody := ""
	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	responseBody, err := Cli.DoRequest(httpMethod, reqUrl, reqBody, nil)
	if err != nil {
		return nil, responseBody, err
	}

	var baseResp BaseResp
	err = okx.UnmarshalOpts.ResponseUnmarshaler(responseBody, &baseResp)
	if err != nil {
		return responseBody, responseBody, err
	}

	if baseResp.Code == 0 {
		logger.Debugf("[DoNoAuthRequest] response=%s", string(responseBody))
		return baseResp.Data, responseBody, nil
	}

	logger.Debugf("[DoNoAuthRequest] error=%s", baseResp.Msg)
	return nil, responseBody, err
}
