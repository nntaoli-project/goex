package okex

import (
	"context"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var (
	client *http.Client
)

func init() {
	logger.SetLevel(logger.DEBUG)
	client = &http.Client{
		Transport: &http.Transport{Proxy: func(request *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1080")
		},
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, e = net.DialTimeout(network, addr, 5*time.Second)
				return conn, e
			},
		},
		Timeout: 10 * time.Second,
	}
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
	ok.OKExV3FutureWs.OrderCallback(func(order *goex.FutureOrder, s string) {
		t.Log(s, order)
	})
	//ok.OKExV3FutureWs.SubscribeTicker(goex.EOS_USD, goex.QUARTER_CONTRACT)
	ok.OKExV3FutureWs.SubscribeDepth(goex.EOS_USDT, goex.QUARTER_CONTRACT, 5)
	//ok.OKExV3FutureWs.SubscribeTrade(goex.EOS_USD, goex.QUARTER_CONTRACT)
	//ok.OKExV3FutureWs.SubscribeOrder(goex.BSV_USD, goex.NEXT_WEEK_CONTRACT)
	time.Sleep(1 * time.Minute)
}
