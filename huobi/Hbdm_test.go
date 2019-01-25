package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
	"time"
)

var dm = NewHbdm(&goex.APIConfig{
	HttpClient: httpProxyClient,
})

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
