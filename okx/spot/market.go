package spot

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"net/url"
)

func (s *spotImpl) GetName() string {
	return "okx.com"
}

func (s *spotImpl) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spotImpl) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	data, err := s.doNoAuthRequest("GET", s.defaultUriOpts.Endpoint+s.defaultUriOpts.TickerUri, &params, nil)
	if err != nil {
		return nil, err
	}

	tk, err := s.tickerUnmarshaler.UnmarshalTicker(data)
	if err != nil {
		return nil, err
	}

	tk.Pair = pair
	tk.Origin = data

	return tk, err
}

func (s *spotImpl) GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, error) {
	//TODO implement me

	panic("implement me")
}

func (s *spotImpl) doNoAuthRequest(method, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	responseData, err := GetHttpCli().DoRequest(method, reqUrl, params, headers)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, errors.New(string(responseData)))
	}

	var resp struct {
		Code int             `json:"code,string"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}

	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Msg)
	}

	return resp.Data, err
}
