package huobi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hb2 = NewV2(http.DefaultClient, "ae42a322-c7a61291-1fc13ecd-b6401", "7f498c3c-19a993d3-72363e3a-0ee8b")

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
