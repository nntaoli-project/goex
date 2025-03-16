package fapi

import (
	"errors"
	"fmt"
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

	}

	return responseBody, responseBody, err
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
	//TODO implement me
	panic("implement me")
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
