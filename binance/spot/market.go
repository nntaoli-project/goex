package spot

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"net/url"
)

func (s *spotImpl) GetName() string {
	return "binance.com"
}

func (s *spotImpl) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spotImpl) GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error) {
	cli := GetHttpCli()
	params := url.Values{}
	params.Set("symbol", fmt.Sprintf("%s%s", pair.Symbol, pair.Market))

	if len(opt) > 0 {
		for _, p := range opt {
			if p.Key == "symbols" {
				params.Del("symbol") //only symbol or symbols
			}
			params.Add(p.Key, p.Value)
		}
	}

	data, err := cli.DoRequest("GET", s.defaultUriOpts.Endpoint+s.defaultUriOpts.TickerUri, &params, nil)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, errors.New(string(data)))
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
