package cexio

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
	"os"
)

var (
	accessKey = ""
	secretKey = ""
	userId = ""
	cex *Cex
)

func init() {
	accessKey = os.Getenv("CEXIO_API_KEY")
	secretKey = os.Getenv("CEXIO_API_SECRET")
	userId = os.Getenv("CEXIO_USERID")
	cex = New(http.DefaultClient, accessKey, secretKey, userId)
}

func TestCex_GetTicker(t *testing.T) {
	e, err := cex.GetTicker(BTC_USD)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if e.Last == 0.0 {
		t.Log("Last == 0.0")
		t.FailNow()
	}
}

func TestCex_GetDepth(t *testing.T) {
	e, err := cex.GetDepth(3, BTC_USD)
	if err != nil {
		t.FailNow()
	}
	t.Log(e)
}

func TestCex_GetAccount(t *testing.T) {
	account, err := cex.GetAccount()
	if err != nil {
		t.FailNow()
	}
	t.Log(account)
}