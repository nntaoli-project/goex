package common

import (
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

type Market struct {
	*V5
}

func (m *Market) GetName() string {
	return "okx.com"
}

func (m *Market) GetDepth(pair CurrencyPair, size int, opt ...OptionParameter) (*Depth, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("sz", fmt.Sprint(size))
	MergeOptionParams(&params, opt...)

	data, err := m.DoNoAuthRequest("GET", m.uriOpts.Endpoint+m.uriOpts.DepthUri, &params)
	if err != nil {
		return nil, err
	}

	dep, err := m.unmarshalOpts.DepthUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	dep.Pair = pair
	//dep.Origin = data

	return dep, err
}

func (m *Market) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	data, err := m.DoNoAuthRequest("GET", m.uriOpts.Endpoint+m.uriOpts.TickerUri, &params)
	if err != nil {
		return nil, err
	}

	tk, err := m.unmarshalOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	tk.Pair = pair
	tk.Origin = data

	return tk, err
}

func (m *Market) GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, error) {
	reqUrl := fmt.Sprintf("%s%s", m.uriOpts.Endpoint, m.uriOpts.KlineUri)
	param := url.Values{}
	param.Set("instId", pair.Symbol)
	param.Set("bar", AdaptKlinePeriodToSymbol(period))
	param.Set("limit", "100")
	MergeOptionParams(&param, opt...)

	data, err := m.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, err
	}
	return m.unmarshalOpts.KlineUnmarshaler(data)
}

func (m *Market) DoNoAuthRequest(httpMethod, reqUrl string, params *url.Values) ([]byte, error) {
	reqBody := ""
	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	data, err := GetHttpCli().DoRequest(httpMethod, reqUrl, reqBody, nil)
	if err != nil {
		return data, err
	}

	var baseResp BaseResp
	err = m.unmarshalOpts.ResponseUnmarshaler(data, &baseResp)
	if err != nil {
		return data, err
	}

	if baseResp.Code == 0 {
		logger.Debugf("[DoNoAuthRequest] response=%s", string(data))
		return baseResp.Data, nil
	}

	logger.Debugf("[DoNoAuthRequest] error=%s", baseResp.Msg)
	return data, err
}
