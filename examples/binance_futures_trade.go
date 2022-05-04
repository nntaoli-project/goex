package main

import (
	"log"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	"github.com/nntaoli-project/goex/builder"
)

const (
	BINANCE_TESTNET_API_KEY        = "***REMOVED***"
	BINANCE_TESTNET_API_KEY_SECRET = "***REMOVED***"
)

func main() {

	api, err := builder.DefaultAPIBuilder.APIKey(BINANCE_TESTNET_API_KEY).APISecretkey(BINANCE_TESTNET_API_KEY_SECRET).Endpoint(binance.TESTNET_FUTURE_USD_WS_BASE_URL).BuildFuturesWs(goex.BINANCE_FUTURES)

	if err != nil {
		log.Fatalln(err.Error())
	}
	api.TickerCallback(func(ticker *goex.FutureTicker) {
		log.Printf("%+v\n", *ticker.Ticker)
	})
	api.SubscribeTicker(goex.BTC_USD, goex.SWAP_USDT_CONTRACT)
	api.DepthCallback(func(depth *goex.Depth) {
		log.Printf("%+v\n", *depth)
	})
	api.SubscribeDepth(goex.BTC_USDT, goex.SWAP_USDT_CONTRACT)

	time.Sleep(time.Minute) // run for one minute

}
