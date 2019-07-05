package okex

import (
	. "github.com/nntaoli-project/GoEx"
)

type OKExFuture struct {
	*OKEx
}

func (ok *OKExFuture) GetExchangeName() string {
	return OKEX_FUTURE
}

func (ok *OKExFuture) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("")
}
func (ok *OKExFuture) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	panic("")
}
func (ok *OKExFuture) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	panic("")
}
func (ok *OKExFuture) GetFutureIndex(currencyPair CurrencyPair) (float64, error) { panic("") }
func (ok *OKExFuture) GetFutureUserinfo() (*FutureAccount, error)                { panic("") }
func (ok *OKExFuture) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	panic("")
}
func (ok *OKExFuture) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("")
}
func (ok *OKExFuture) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("")
}
func (ok *OKExFuture) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("")
}
func (ok *OKExFuture) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("")
}
func (ok *OKExFuture) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("")
}
func (ok *OKExFuture) GetFee() (float64, error)                                    { panic("") }
func (ok *OKExFuture) GetContractValue(currencyPair CurrencyPair) (float64, error) { panic("") }
func (ok *OKExFuture) GetDeliveryTime() (int, int, int, int)                       { panic("") }
func (ok *OKExFuture) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("")
}
func (ok *OKExFuture) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("")
}
