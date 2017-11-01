package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var hb2 = NewV2(http.DefaultClient, "", "")

func TestHuoBi_V2_GetTicker(t *testing.T) {
	ticker, err := hb2.GetTicker(goex.BTS_CNY)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestHuoBi_V2_GetDepth(t *testing.T) {
	depth, err := hb2.GetDepth(2, goex.BCC_CNY)
	assert.Nil(t, err)
	t.Log("asks: ", depth.AskList)
	t.Log("bids: ", depth.BidList)
}
