package okex

import (
	"github.com/nntaoli-project/goex"
	log "github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"testing"
)

func init() {
	log.SetLevel(log.DEBUG)
	//os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:2222") //local socks5 proxy
}

func TestOKExV5Swap_GetFutureTicker(t *testing.T) {
	swap := NewOKExV5Swap(&goex.APIConfig{
		HttpClient:    http.DefaultClient,
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
		Lever:         0,
	})
	t.Log(swap.GetFutureTicker(goex.BTC_USDT, goex.SWAP_CONTRACT))
}

func TestOKExV5Swap_GetFutureDepth(t *testing.T) {
	swap := NewOKExV5Swap(&goex.APIConfig{
		HttpClient: http.DefaultClient,
	})

	dep, err := swap.GetFutureDepth(goex.BTC_USDT, goex.SWAP_CONTRACT, 2)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
