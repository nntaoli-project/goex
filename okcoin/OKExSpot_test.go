package okcoin

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var okexSpot = NewOKExSpot(http.DefaultClient, "", "")

func TestOKExSpot_GetTicker(t *testing.T) {
	ticker, err := okexSpot.GetTicker(goex.ETC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestOKExSpot_GetDepth(t *testing.T) {
	dep, err := okexSpot.GetDepth(2, goex.ETC_BTC)
	assert.Nil(t, err)
	t.Log(dep)
}
