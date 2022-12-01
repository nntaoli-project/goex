package spot

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	. "github.com/nntaoli-project/goex/v2/util"
	"net/http"
	"net/url"
)

func (s *spotImpl) GetName() string {
	return "binance.com"
}

func (s *spotImpl) GetDepth(pair CurrencyPair, size int, opts ...OptionParameter) (*Depth, error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)
	params.Set("limit", fmt.Sprint(size))
	MergeOptionParams(&params, opts...)

	reqUrl := fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.DepthUri)
	data, err := s.doNoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, err
	}
	logger.Debugf("[GetDepth] %s", string(data))
	return s.unmarshalerOpts.DepthUnmarshaler(data)
}

func (s *spotImpl) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error) {
	params := url.Values{}
	params.Set("symbol", pair.Symbol)

	if len(opt) > 0 {
		for _, p := range opt {
			if p.Key == "symbols" {
				params.Del("symbol") //only symbol or symbols
			}
			params.Add(p.Key, p.Value)
		}
	}

	data, err := s.doNoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.TickerUri), &params, nil)
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

func (s *spotImpl) GetKline(pair CurrencyPair, period KlinePeriod, opts ...OptionParameter) ([]Kline, error) {
	params := url.Values{}
	params.Set("limit", "1000")
	params.Set("symbol", pair.Symbol)
	params.Set("interval", adaptKlinePeriod(period))
	MergeOptionParams(&params, opts...)

	reqUrl := fmt.Sprintf("%s%s", s.uriOpts.Endpoint, s.uriOpts.KlineUri)
	respBody, err := s.doNoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, err
	}

	return s.unmarshalerOpts.KlineUnmarshaler(respBody)
}

func (s *spotImpl) doNoAuthRequest(method, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	var reqBody string

	if method == http.MethodGet {
		reqUrl += "?" + params.Encode()
	} else {
		reqBody = params.Encode()
	}

	responseData, err := GetHttpCli().DoRequest(method, reqUrl, reqBody, headers)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, errors.New(string(responseData)))
	}

	return responseData, err
}
