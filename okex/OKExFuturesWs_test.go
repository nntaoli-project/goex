package okex

import (
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	client *http.Client
)

func init() {
	logger.SetLevel(logger.DEBUG)
}

func TestNewOKExV3FuturesWs(t *testing.T) {
	os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	ok := NewOKEx(&goex.APIConfig{
		HttpClient: http.DefaultClient,
	})
	ok.OKExV3FutureWs.TickerCallback(func(ticker *goex.FutureTicker) {
		t.Log(ticker.Ticker, ticker.ContractType)
	})
	ok.OKExV3FutureWs.DepthCallback(func(depth *goex.Depth) {
		t.Log(depth)
	})
	ok.OKExV3FutureWs.TradeCallback(func(trade *goex.Trade, s string) {
		t.Log(s, trade)
	})
	//ok.OKExV3FutureWs.SubscribeTicker(goex.EOS_USD, goex.QUARTER_CONTRACT)
	ok.OKExV3FutureWs.SubscribeDepth(goex.EOS_USD, goex.QUARTER_CONTRACT)
	//ok.OKExV3FutureWs.SubscribeTrade(goex.EOS_USD, goex.QUARTER_CONTRACT)
	time.Sleep(1 * time.Minute)
}
