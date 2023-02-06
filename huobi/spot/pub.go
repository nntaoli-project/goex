package spot

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	. "github.com/nntaoli-project/goex/v2/model"
	"net/http"
	"net/url"
)

func (s *Spot) GetName() string {
	return "huobi.com"
}

func (s *Spot) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, []byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Spot) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, []byte, error) {
	data, err := s.DoNoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s?symbol=%s", s.uriOpts.Endpoint, s.uriOpts.TickerUri, pair.Symbol), nil, nil)
	if err != nil {
		return nil, data, fmt.Errorf("%w%s", err, errors.New(string(data)))
	}

	tk, err := s.unmarshalerOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	tk.Pair = pair

	return tk, data, err
}

func (s *Spot) GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, []byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Spot) GetExchangeInfo() (map[string]CurrencyPair, []byte, error) {
	panic("not implement")
}

func (s *Spot) DoNoAuthRequest(method, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	if method == http.MethodGet && params != nil {
		reqUrl += "?" + params.Encode()
	}

	responseData, err := Cli.DoRequest(method, reqUrl, "", headers)
	if err != nil {
		return responseData, fmt.Errorf("%w%s", err, errors.New(string(responseData)))
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
