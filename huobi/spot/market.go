package spot

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	. "github.com/nntaoli-project/goex/v2/model"
	"net/http"
	"net/url"
)

func (s *spotImpl) GetName() string {
	return "huobi.com"
}

func (s *spotImpl) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spotImpl) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error) {
	data, err := s.doNoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s?symbol=%s", s.uriOpts.Endpoint, s.uriOpts.TickerUri, pair.Symbol), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, errors.New(string(data)))
	}

	tk, err := s.unmarshalerOpts.TickerUnmarshaler(data)
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
	if method == http.MethodGet && params != nil {
		reqUrl += "?" + params.Encode()
	}

	responseData, err := GetHttpCli().DoRequest(method, reqUrl, "", headers)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, errors.New(string(responseData)))
	}

	var resp BaseResponse

	err = s.unmarshalerOpts.ResponseUnmarshaler(responseData, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, errors.New(string(responseData))
	}

	return responseData, nil
}
