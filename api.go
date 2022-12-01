package goex

import (
	"github.com/nntaoli-project/goex/v2/model"
	"net/url"
)

type IHttpClient interface {
	DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error)
}

// IMarketRest 行情接口，不需要授权
type IMarketRest interface {
	GetName() string //获取交易所名字/域名
	GetDepth(pair model.CurrencyPair, limit int, opt ...model.OptionParameter) (*model.Depth, error)
	GetTicker(pair model.CurrencyPair, opt ...model.OptionParameter) (*model.Ticker, error)
	GetKline(pair model.CurrencyPair, period model.KlinePeriod, opt ...model.OptionParameter) ([]model.Kline, error)
}

// ITradeRest 交易相关的接口
type ITradeRest interface {
	CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opt ...model.OptionParameter) (*model.Order, error) //创建订单
	//CreateOrders 批量创建订单,考虑中，是否有意义
	//CreateOrders(orders []Order, opt ...OptionParameter) ([]Order, error)
	GetOrderInfo(pair model.CurrencyPair, id string, opt ...model.OptionParameter) (*model.Order, error)
	GetPendingOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, error)
	// GetHistoryOrders 获取历史委托订单列表
	GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, error)
	CancelOrder(pair model.CurrencyPair, id string, opt ...model.OptionParameter) error
	//CancelOrders(pair *CurrencyPair, id []string, opt ...OptionParameter) error
	DoAuthRequest(method, reqUrl string, params *url.Values, header map[string]string) ([]byte, error)
}

type IFuturesPosition interface {
	GetPositions(pair model.CurrencyPair, opts ...model.OptionParameter) ([]model.FuturesPosition, error)
}

// IAccount
// 获取账户资产相关的
type IAccount interface {
	GetAccount(coin string) (map[string]model.Account, error)
}

// IWallet 获取资产信息，划转资金等操作
type IWallet interface {
}
