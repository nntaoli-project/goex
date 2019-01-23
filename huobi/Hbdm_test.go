package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
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
	t.Logf("%+v\n%+v", dep.AskList , dep.BidList)
}
