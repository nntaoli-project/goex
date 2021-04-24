package binance

import (
	"github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"net/http"
	"testing"
)

var baDapi = NewBinanceFutures(&goex.APIConfig{
	HttpClient:   http.DefaultClient,
	ApiKey:       "",
	ApiSecretKey: "",
})

func init() {
	logger.SetLevel(logger.DEBUG)
}

func TestBinanceFutures_GetFutureDepth(t *testing.T) {
	t.Log(baDapi.GetFutureDepth(goex.ETH_USD, goex.QUARTER_CONTRACT, 10))
}

func TestBinanceSwap_GetFutureTicker(t *testing.T) {
	ticker, err := baDapi.GetFutureTicker(goex.LTC_USD, goex.SWAP_CONTRACT)
	t.Log(err)
	t.Logf("%+v", ticker)
}

func TestBinance_GetExchangeInfo(t *testing.T) {
	baDapi.GetExchangeInfo()
}

func TestBinanceFutures_GetFutureUserinfo(t *testing.T) {
	t.Log(baDapi.GetFutureUserinfo())
}

func TestBinanceFutures_PlaceFutureOrder(t *testing.T) {
	//1044675677
	t.Log(baDapi.PlaceFutureOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "19990", "2", goex.OPEN_SELL, 0, 10))
}

func TestBinanceFutures_LimitFuturesOrder(t *testing.T) {
	t.Log(baDapi.LimitFuturesOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "20001", "2", goex.OPEN_SELL))
}

func TestBinanceFutures_MarketFuturesOrder(t *testing.T) {
	t.Log(baDapi.MarketFuturesOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "2", goex.OPEN_SELL))
}

func TestBinanceFutures_GetFutureOrder(t *testing.T) {
	t.Log(baDapi.GetFutureOrder("1045208666", goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestBinanceFutures_FutureCancelOrder(t *testing.T) {
	t.Log(baDapi.FutureCancelOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "1045328328"))
}

func TestBinanceFutures_GetFuturePosition(t *testing.T) {
	t.Log(baDapi.GetFuturePosition(goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestBinanceFutures_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(baDapi.GetUnfinishFutureOrders(goex.BTC_USD , goex.QUARTER_CONTRACT))
}
