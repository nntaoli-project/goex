package okcoin

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
	"time"
)

func TestNewOKExFutureWs(t *testing.T) {
	okWs := NewOKExFutureWs()
	okWs.ErrorHandleFunc(func(err error) {
		t.Log(err)
	})
	okWs.SetCallbacks(func(ticker *goex.FutureTicker) {
		t.Log(ticker, ticker.Ticker)
	}, func(depth *goex.Depth) {
		t.Log(depth.ContractType, depth.Pair, depth.AskList, depth.BidList)
	}, func(trade *goex.Trade, contract string) {
		t.Log(contract, trade)
	})
	okWs.SubscribeTicker(goex.LTC_USD, goex.QUARTER_CONTRACT)
	okWs.SubscribeDepth(goex.LTC_USD, goex.QUARTER_CONTRACT, 5)
	okWs.SubscribeTrade(goex.LTC_USD, goex.QUARTER_CONTRACT)
	time.Sleep(3 * time.Minute)
}
