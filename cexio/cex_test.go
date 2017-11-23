package cexio

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

func TestCex_GetTicker(t *testing.T) {
	cex := New(http.DefaultClient, "", "", "")
	e, err := cex.GetTicker(BTC_USD)
	if err != nil {
		t.FailNow()
	}
	if e.Last == 0.0 {
		t.FailNow()
	}
}

func TestCex_GetDepth(t *testing.T) {
	cex := New(http.DefaultClient, "", "", "")
	e, err := cex.GetDepth(3, BTC_USD)
	if err != nil {
		t.FailNow()
	}
	t.Log(e)
}
