package kucoin

import (
	"github.com/nntaoli-project/goex"
	"testing"
	"time"
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
	kLines, _ := kc.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 10, int(time.Now().Unix()-3600))
	t.Log(kLines)
}

func TestKuCoin_GetTrades(t *testing.T) {
	trades, _ := kc.GetTrades(goex.BTC_USDT, 0)
	t.Log(trades)
}

