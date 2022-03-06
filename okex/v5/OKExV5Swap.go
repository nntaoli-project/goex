package okex

import (
	"fmt"
	. "github.com/nntaoli-project/goex"
	"sort"
	"time"
)

type OKExV5Swap struct {
	*OKExV5
}

func NewOKExV5Swap(config *APIConfig) *OKExV5Swap {
	v5 := new(OKExV5Swap)
	v5.OKExV5 = NewOKExV5(config)
	return v5
}

func (O *OKExV5Swap) GetExchangeName() string {
	return OKEX_SWAP
}

func (O *OKExV5Swap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	t, err := O.OKExV5.GetTickerV5(fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-")))

	if err != nil {
		return nil, err
	}

	return &Ticker{
		Pair: currencyPair,
		Last: t.Last,
		Buy:  t.BuyPrice,
		Sell: t.SellPrice,
		High: t.High,
		Low:  t.Low,
		Vol:  t.Vol,
		Date: t.Timestamp,
	}, nil
}

func (O *OKExV5Swap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	instId := fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	dep, err := O.OKExV5.GetDepthV5(instId, size)

	if err != nil {
		return nil, err
	}

	depth := &Depth{}

	for _, ask := range dep.Asks {
		depth.AskList = append(depth.AskList, DepthRecord{Price: ToFloat64(ask[0]), Amount: ToFloat64(ask[1])})
	}

	for _, bid := range dep.Bids {
		depth.BidList = append(depth.BidList, DepthRecord{Price: ToFloat64(bid[0]), Amount: ToFloat64(bid[1])})
	}

	sort.Sort(sort.Reverse(depth.AskList))

	depth.Pair = currencyPair
	depth.UTime = time.Unix(0, int64(dep.Timestamp)*1000000)

	return depth, nil
}

func (O *OKExV5Swap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	panic("implement me")
}

func (O *OKExV5Swap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	panic("implement me")
}

func (O *OKExV5Swap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetFee() (float64, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}

func (O OKExV5Swap) GetDeliveryTime() (int, int, int, int) {
	panic("implement me")
}

func (O *OKExV5Swap) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]FutureKline, error) {
	panic("implement me")
}

func (O *OKExV5Swap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("implement me")
}
