package binance

import (
	"github.com/nntaoli-project/goex"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var ba = NewWithConfig(
	&goex.APIConfig{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse("socks5://127.0.0.1:1080")
					return nil, nil
				},
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
			},
			Timeout: 10 * time.Second,
		},
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       "q6y6Gr7fF3jSJLncpfn2PmAA0xu4XRiRFHpFkyJy3d7K68WUxY0Gt8rrajCDUfbI",
		ApiSecretKey: "AP8C2kh4RyISN3fpRCFMZJddf233XbPcYWQ1S7gBan3pGjCQg2JnyQFSJrIaNzRh",
	})

func TestBinance_GetTicker(t *testing.T) {
	ticker, err := ba.GetTicker(goex.NewCurrencyPair2("USDT_USD"))
	t.Log(ticker, err)
}

func TestBinance_LimitBuy(t *testing.T) {
	order, err := ba.LimitBuy("0.005", "8000", goex.BTC_USDT)
	t.Log(order, err)
}

func TestBinance_LimitSell(t *testing.T) {
	order, err := ba.LimitSell("0.01", "0.1", goex.LTC_BTC)
	t.Log(order, err)
}

func TestBinance_CancelOrder(t *testing.T) {
	t.Log(ba.CancelOrder("1156274704", goex.BTC_USDT))
}

func TestBinance_GetOneOrder(t *testing.T) {
	t.Log(ba.GetOneOrder("1156274704", goex.BTC_USDT))
}

func TestBinance_GetDepth(t *testing.T) {
	//return
	dep, err := ba.GetDepth(5, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
}

func TestBinance_GetAccount(t *testing.T) {
	account, err := ba.GetAccount()
	t.Log(account, err)
}

func TestBinance_GetUnfinishOrders(t *testing.T) {
	orders, err := ba.GetUnfinishOrders(goex.ETH_BTC)
	t.Log(orders, err)
}

func TestBinance_GetKlineRecords(t *testing.T) {
	before := time.Now().Add(-time.Hour).Unix() * 1000
	kline, _ := ba.GetKlineRecords(goex.ETH_BTC, goex.KLINE_PERIOD_5MIN, 100, int(before))
	for _, k := range kline {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(tt, k.Open, k.Close, k.High, k.Low, k.Vol)
	}
}

func TestBinance_GetTrades(t *testing.T) {
	t.Log(ba.GetTrades(goex.BTC_USDT, 0))
}

func TestBinance_GetTradeSymbols(t *testing.T) {
	t.Log(ba.GetTradeSymbol(goex.BTC_USDT))
}

func TestBinance_SetTimeOffset(t *testing.T) {
	t.Log(ba.setTimeOffset())
	t.Log(ba.timeOffset)
}

func TestBinance_GetOrderHistorys(t *testing.T) {
	t.Log(ba.GetOrderHistorys(goex.BTC_USDT, 1, 1))
}
