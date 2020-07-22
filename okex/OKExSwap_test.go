package okex

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var config = &goex.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   "127.0.0.1:1080"}, nil
			},
		},
	},
	Endpoint:      "https://www.okex.com",
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}

var okExSwap = NewOKExSwap(config)

func TestOKExSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(okExSwap.GetFutureUserinfo())
}

func TestOKExSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(okExSwap.PlaceFutureOrder(goex.XRP_USD, goex.SWAP_CONTRACT, "0.2", "10", goex.OPEN_BUY, 0, 0))
}

func TestOKExSwap_FutureCancelOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelOrder(goex.XRP_USD, goex.SWAP_CONTRACT, "309935122485305344"))
}

func TestOKExSwap_GetFutureOrder(t *testing.T) {
	t.Log(okExSwap.GetFutureOrder("309935122485305344", goex.XRP_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetFuturePosition(t *testing.T) {
	t.Log(okExSwap.GetFuturePosition(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetFutureDepth(t *testing.T) {
	t.Log(okExSwap.GetFutureDepth(goex.LTC_USD, goex.SWAP_CONTRACT, 10))
}

func TestOKExSwap_GetFutureTicker(t *testing.T) {
	t.Log(okExSwap.GetFutureTicker(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetUnfinishFutureOrders(t *testing.T) {
	ords, _ := okExSwap.GetUnfinishFutureOrders(goex.XRP_USD, goex.SWAP_CONTRACT)
	for _, ord := range ords {
		t.Log(ord.OrderID2, ord.ClientOid)
	}

}

func TestOKExSwap_GetHistoricalFunding(t *testing.T) {
	for i := 1; ; i++ {
		funding, err := okExSwap.GetHistoricalFunding(goex.SWAP_CONTRACT, goex.BTC_USD, i)
		t.Log(err, len(funding))
	}
}

func TestOKExSwap_GetKlineRecords(t *testing.T) {
	since := time.Now().Add(-24 * time.Hour).Unix()
	kline, err := okExSwap.GetKlineRecords(goex.SWAP_CONTRACT, goex.BTC_USD, goex.KLINE_PERIOD_4H, 0, int(since))
	t.Log(err, kline[0].Kline)
}

func TestOKExSwap_GetKlineRecords2(t *testing.T) {
	start := time.Now().Add(time.Minute * -30).UTC().Format(time.RFC3339)
	t.Log(start)
	kline, err := okExSwap.GetKlineRecords2(goex.SWAP_CONTRACT, goex.BTC_USDT, start, "", "900")
	t.Log(err, kline[0].Kline)
}

func TestOKExSwap_GetInstruments(t *testing.T) {
	t.Log(okExSwap.GetInstruments())
}

func TestOKExSwap_SetMarginLevel(t *testing.T) {
	t.Log(okExSwap.SetMarginLevel(goex.EOS_USDT, 5, 3))
}

func TestOKExSwap_GetMarginLevel(t *testing.T) {
	t.Log(okExSwap.GetMarginLevel(goex.EOS_USDT))
}

func TestOKExSwap_GetFutureAccountInfo(t *testing.T) {
	t.Log(okExSwap.GetFutureAccountInfo(goex.BTC_USDT))
}

func TestOKExSwap_PlaceFutureAlgoOrder(t *testing.T) {
	ord := &goex.FutureOrder{
		ContractName: goex.SWAP_CONTRACT,
		Currency:     goex.BTC_USD,
		OType:        2, //开空
		OrderType:    1, //1：止盈止损 2：跟踪委托 3：冰山委托 4：时间加权
		Price:        9877,
		Amount:       1,

		TriggerPrice: 9877,
		AlgoType:     1,
	}
	t.Log(okExSwap.PlaceFutureAlgoOrder(ord))
}

func TestOKExSwap_FutureCancelAlgoOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelAlgoOrder(goex.BTC_USD, []string{"309935122485305344"}))

}

func TestOKExSwap_GetFutureAlgoOrders(t *testing.T) {
	t.Log(okExSwap.GetFutureAlgoOrders("", "2", goex.BTC_USD))
}
