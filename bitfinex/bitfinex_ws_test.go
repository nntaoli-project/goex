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

	bitfinexWs.SetCallbacks(handleTicker, handleTrade)

	//Ticker
	t.Log(bitfinexWs.SubscribeTicker(goex.BTC_USD))
	t.Log(bitfinexWs.SubscribeTicker(goex.LTC_USD))

	//Trades
	t.Log(bitfinexWs.SubscribeTrade(goex.BTC_USD))

	time.Sleep(time.Minute)
}
