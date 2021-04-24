package huobi

import (
	"github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"testing"
	"time"
)

func TestNewHbdmSwapWs(t *testing.T) {
	logger.SetLevel(logger.DEBUG)

	ws := NewHbdmSwapWs()

	ws.DepthCallback(func(depth *goex.Depth) {
		t.Log(depth)
	})
	ws.TickerCallback(func(ticker *goex.FutureTicker) {
		t.Log(ticker.Date, ticker.Last, ticker.Buy, ticker.Sell, ticker.High, ticker.Low, ticker.Vol)
	})
	ws.TradeCallback(func(trade *goex.Trade, contract string) {
		t.Log(trade, contract)
	})

	//t.Log(ws.SubscribeDepth(goex.BTC_USD, goex.SWAP_CONTRACT))
	//t.Log(ws.SubscribeTicker(goex.BTC_USD, goex.SWAP_CONTRACT))
	t.Log(ws.SubscribeTrade(goex.BTC_USD , goex.SWAP_CONTRACT))
	
	time.Sleep(time.Minute)
}
