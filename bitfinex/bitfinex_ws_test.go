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
	
	bitfinexWs.SetCallbacks(handleTicker)

	t.Log(bitfinexWs.SubscribeTicker(goex.BTC_USD))
	t.Log(bitfinexWs.SubscribeTicker(goex.LTC_USD))
	time.Sleep(time.Minute)
}
