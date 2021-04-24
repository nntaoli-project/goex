package kucoin

import (
	"github.com/Jameslu041/goex"
	"testing"
)

var kc = New("", "", "")

func TestKuCoin_GetTicker(t *testing.T) {
	ticker, _ := kc.GetTicker(goex.BTC_USDT)
	t.Log(ticker)
}

func TestKuCoin_GetDepth(t *testing.T) {
	depth, _ := kc.GetDepth(10, goex.BTC_USDT)
	t.Log(depth)
}

func TestKuCoin_GetKlineRecords(t *testing.T) {
	kLines, _ := kc.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 10)
	t.Log(kLines)
}

func TestKuCoin_GetTrades(t *testing.T) {
	trades, _ := kc.GetTrades(goex.BTC_USDT, 0)
	t.Log(trades)
}

func TestKuCoin_GetAccount(t *testing.T) {
	acc, _ := kc.GetAccount()
	t.Log(acc)
}
