package goex

import (
	"net/url"
)

type IHttpClient interface {
	DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error)
}

// IMarketRest 行情接口，不需要授权
type IMarketRest interface {
	GetName() string //获取交易所名字/域名
	GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error)
	GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error)
	GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, error)
}

// ITradeRest 交易相关的接口
type ITradeRest interface {
	CreateOrder(order Order, opt ...OptionParameter) (*Order, error) //创建订单
	//CreateOrders 批量创建订单,考虑中，是否有意义
	//CreateOrders(orders []Order, opt ...OptionParameter) ([]Order, error)
	GetOrderInfo(pair CurrencyPair, id string, opt ...OptionParameter) (*Order, error)
	GetPendingOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, error)
	// GetHistoryOrders 获取历史委托订单列表
	GetHistoryOrders(pair CurrencyPair, opt ...OptionParameter) ([]Order, error)
	CancelOrder(pair CurrencyPair, id string, opt ...OptionParameter) error
	//CancelOrders(pair *CurrencyPair, id []string, opt ...OptionParameter) error
	DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error)
}

// IWallet 获取资产信息，划转资金等操作
type IWallet interface {
}
