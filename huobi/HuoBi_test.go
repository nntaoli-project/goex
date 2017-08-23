package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var hb = New(http.DefaultClient, "", "")

func TestHuoBi_GetDepth(t *testing.T) {
	dep, err := hb.GetDepth(2, goex.BTC_CNY)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestHuoBi_GetKlineRecords(t *testing.T) {
	klines, err := hb.GetKlineRecords(goex.BTC_CNY, goex.KLINE_PERIOD_4H, 1, -1)
	assert.Nil(t, err)
	t.Log(klines)
}
