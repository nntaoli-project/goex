package goex

import (
	"github.com/nntaoli-project/goex/v2/model"
)

// IPubRest is a public interface that does not require authorization."
type IPubRest interface {
	GetName() string //获取交易所名字/域名
	// GetDepth
	//    bids: 降序
	//    asks: 升序
	GetDepth(pair model.CurrencyPair, limit int, opt ...model.OptionParameter) (depth *model.Depth, responseBody []byte, err error)
	GetTicker(pair model.CurrencyPair, opt ...model.OptionParameter) (ticker *model.Ticker, responseBody []byte, err error)
	GetKline(pair model.CurrencyPair, period model.KlinePeriod, opt ...model.OptionParameter) (klines []model.Kline, responseBody []byte, err error)
	GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error)
	// NewCurrencyPair 同时支持现货和期货
	//@parameter
	//  - bashSym
	//  - quoteSym
	//	- opts 交割合约的时候传入contract alias name: this_week,next_week, quarter ...
	NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error)
}

//IPrvRest is a private interface specification that requires authorization to call.
type IPrvRest interface {
	GetAccount(coin string) (map[string]model.Account, []byte, error)
	//CreateOrder
	//@returns
	//  order        包含订单ID信息
	//  responseBody 交易所接口返回的原始字节数据
	//  err          错误
	CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opt ...model.OptionParameter) (order *model.Order, responseBody []byte, err error)
	GetOrderInfo(pair model.CurrencyPair, id string, opt ...model.OptionParameter) (order *model.Order, responseBody []byte, err error)
	GetPendingOrders(pair model.CurrencyPair, opt ...model.OptionParameter) (orders []model.Order, responseBody []byte, err error)
	GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) (orders []model.Order, responseBody []byte, err error)
	CancelOrder(pair model.CurrencyPair, id string, opt ...model.OptionParameter) (responseBody []byte, err error)
}

type ISpotPrvRest interface {
	IPrvRest
}

type IFuturesPubRest interface {
	IPubRest
	//GetFundingRate
	//获取资金费率，仅适用于永续合约
	GetFundingRate(pair model.CurrencyPair, opts ...model.OptionParameter) (rate *model.FundingRate, responseBody []byte, err error)
}

// IFuturesPrvRest includes some special interface implementations for futures supplement.
type IFuturesPrvRest interface {
	IPrvRest
	GetFuturesAccount(coin string) (acc map[string]model.FuturesAccount, responseBody []byte, err error)
	//GetPositions 获取持仓数据
	//@returns
	//	positions    仓位数据
	//	responseBody 交易所接口返回的原始字节数据
	//	err          错误
	GetPositions(pair model.CurrencyPair, opts ...model.OptionParameter) (positions []model.FuturesPosition, responseBody []byte, err error)
}
