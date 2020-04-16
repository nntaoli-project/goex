package bitfinex

import (
	"log"
	"testing"
	"time"

	"github.com/nntaoli-project/goex"
)

func TestNewBitfinexWs(t *testing.T) {
	bitfinexWs := NewWs()

	handleTicker := func(ticker *goex.Ticker) {
		log.Printf("Ticker: %+v: ", ticker)
	}

	handleTrade := func(trade *goex.Trade) {
		log.Printf("Trade: %+v: ", trade)
	}

	handleCandle := func(candle *goex.Kline) {
		log.Printf("Candle: %+v: ", candle)
	}

	bitfinexWs.SetCallbacks(handleTicker, handleTrade, handleCandle)

	//Ticker
	t.Log(bitfinexWs.SubscribeTicker(goex.BTC_USD))
	t.Log(bitfinexWs.SubscribeTicker(goex.LTC_USD))

	//Trades
	t.Log(bitfinexWs.SubscribeTrade(goex.BTC_USD))

	//Candles
	t.Log(bitfinexWs.SubscribeCandle(goex.BTC_USD, goex.KLINE_PERIOD_1MIN))
	
	time.Sleep(time.Minute)
}
