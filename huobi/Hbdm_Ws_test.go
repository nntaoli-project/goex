package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"log"
	"testing"
	"time"
)

func TestNewHbdmWs(t *testing.T) {
	ws := NewHbdmWs()
	ws.ProxyUrl("socks5://127.0.0.1:1080")

	ws.SetCallbacks(func(ticker *goex.FutureTicker) {
		log.Println(ticker.Ticker)
	}, func(depth *goex.Depth) {
		log.Println(">>>>>>>>>>>>>>>")
		log.Println(depth.ContractType, depth.Pair)
		log.Println(depth.BidList)
		log.Println(depth.AskList)
		log.Println("<<<<<<<<<<<<<<")
	}, func(trade *goex.Trade, s string) {
		log.Println(s, trade)
	})

	t.Log(ws.SubscribeTicker(goex.BTC_USD, goex.QUARTER_CONTRACT))
	t.Log(ws.SubscribeDepth(goex.BTC_USD, goex.NEXT_WEEK_CONTRACT, 0))
	t.Log(ws.SubscribeTrade(goex.LTC_USD, goex.THIS_WEEK_CONTRACT))
	time.Sleep(time.Minute)
}
