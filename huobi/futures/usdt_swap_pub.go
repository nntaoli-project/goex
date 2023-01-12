package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

func (f *USDTSwap) GetName() string {
	return "hbdm.com"
}

func (f *USDTSwap) DoNoAuthRequest(method, reqUrl string, params *url.Values) ([]byte, error) {
	if method == http.MethodGet {
		reqUrl += "?" + params.Encode()
	}

	respBodyData, err := Cli.DoRequest(method, reqUrl, "", map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		return respBodyData, err
	}

	var baseResp BaseResponse
	err = json.Unmarshal(respBodyData, &baseResp)
	if err != nil {
		logger.Errorf("[DoNoAuthRequest] err=%s", err.Error())
		return nil, err
	}

	if baseResp.Status != "ok" {
		return respBodyData, errors.New(string(respBodyData))
	}

	return respBodyData, nil
}

func (f *USDTSwap) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, []byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *USDTSwap) GetTicker(pair CurrencyPair, opts ...OptionParameter) (*Ticker, []byte, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	MergeOptionParams(&params, opts...)

	data, err := f.DoNoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.TickerUri), &params)
	if err != nil {
		return nil, data, err
	}

	tk, err := f.unmarshalerOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	tk.Pair = pair

	return tk, data, nil
}

func (f *USDTSwap) GetKline(pair CurrencyPair, period KlinePeriod, opts ...OptionParameter) ([]Kline, []byte, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	params.Set("period", AdaptKlinePeriod(period))

	MergeOptionParams(&params, opts...)

	if params.Get("size") == "" && params.Get("from") == "" {
		params.Set("size", "100")
	}

	data, err := f.DoNoAuthRequest(http.MethodGet, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.KlineUri), &params)
	if err != nil {
		return nil, data, err
	}
	logger.Debugf("[GetKline] data=%s", string(data))

	klines, err := f.unmarshalerOpts.KlineUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	for i, _ := range klines {
		klines[i].Pair = pair
	}

	return klines, data, err
}
