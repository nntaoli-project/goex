package main

import (
	"log"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	"github.com/nntaoli-project/goex/builder"
)

const (
	BINANCE_TESTNET_API_KEY        = "e78d5b2ce4cf436da0194610c1da797bf6fa3ea7c4342b984f13d42fd17fd515"
	BINANCE_TESTNET_API_KEY_SECRET = "5cea0042c43db6438b1577860096c3ba83bcc20b34e76b920b31684de712d9d2"
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
