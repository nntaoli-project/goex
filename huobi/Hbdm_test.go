package huobi

import (
	"github.com/nntaoli-project/goex"
	"testing"
	"time"
)

var dm = NewHbdm(&goex.APIConfig{
	Endpoint:     "https://api.hbdm.com",
	HttpClient:   httpProxyClient,
	ApiKey:       "",
	ApiSecretKey: ""})

func TestHbdm_GetFutureUserinfo(t *testing.T) {
	t.Log(dm.GetFutureUserinfo())
}

func TestHbdm_GetFuturePosition(t *testing.T) {
	t.Log(dm.GetFuturePosition(goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestHbdm_PlaceFutureOrder(t *testing.T) {
	t.Log(dm.PlaceFutureOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "3800", "1", goex.OPEN_BUY, 0, 20))
}

func TestHbdm_FutureCancelOrder(t *testing.T) {
	t.Log(dm.FutureCancelOrder(goex.BTC_USD, goex.QUARTER_CONTRACT, "6"))
}

func TestHbdm_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(dm.GetUnfinishFutureOrders(goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureOrders(t *testing.T) {
	t.Log(dm.GetFutureOrders([]string{"6", "5"}, goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureOrder(t *testing.T) {
	t.Log(dm.GetFutureOrder("6", goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureTicker(t *testing.T) {
	t.Log(dm.GetFutureTicker(goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureDepth(t *testing.T) {
	dep, err := dm.GetFutureDepth(goex.BTC_USD, goex.QUARTER_CONTRACT, 0)
	t.Log(err)
	t.Logf("%+v\n%+v", dep.AskList, dep.BidList)
}
func TestHbdm_GetFutureIndex(t *testing.T) {
	t.Log(dm.GetFutureIndex(goex.BTC_USD))
}

func TestHbdm_GetFutureEstimatedPrice(t *testing.T) {
	t.Log(dm.GetFutureEstimatedPrice(goex.BTC_USD))
}

func TestHbdm_GetKlineRecords(t *testing.T) {
	klines, _ := dm.GetKlineRecords(goex.QUARTER_CONTRACT, goex.EOS_USD, goex.KLINE_PERIOD_1MIN, 20, 0)
	for _, k := range klines {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(k.Pair, tt, k.Open, k.Close, k.High, k.Low, k.Vol, k.Vol2)
	}
}
