package goex

import (
	"net/url"
)

type IHttpClient interface {
	DoRequest(method, rqUrl string, params *url.Values, headers map[string]string) (data []byte, err error)
}

// IProtocolParser 协议解析器
type IProtocolParser interface {
}

// IMarketRest 行情接口，不需要授权
type IMarketRest interface {
	GetInfo() ([]byte, error) //获取交易规则信息
	GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error)
	GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error)
	GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, error)
}
