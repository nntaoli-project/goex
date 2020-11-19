package binance

import (
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"testing"
)

var baDapi = NewBinanceFutures(&goex.APIConfig{
	Endpoint: "https://dapi.binancezh.pro",
	//HttpClient: &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: func(req *http.Request) (*url.URL, error) {
	//			return url.Parse("socks5://127.0.0.1:1080")
	//			return nil, nil
	//		},
	//	},
	//	Timeout: 10 * time.Second,
	//},
	HttpClient:   http.DefaultClient,
	ApiKey:       "",
	ApiSecretKey: "",
})

func init() {
	logger.SetLevel(logger.DEBUG)
}

func TestBinanceFutures_GetFutureDepth(t *testing.T) {
	t.Log(baDapi.GetFutureDepth(goex.ETH_USD, goex.BI_QUARTER_CONTRACT, 10))
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
