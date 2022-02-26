package okex

import (
	log "github.com/nntaoli-project/goex/internal/logger"
	"testing"

	"github.com/nntaoli-project/goex"
)

func newOKExV5SpotClient() *OKExV5Spot {
	return NewOKExV5Spot(&goex.APIConfig{
		//HttpClient: &http.Client{
		//	Transport: &http.Transport{
		//		Proxy: func(req *http.Request) (*url.URL, error) {
		//			return &url.URL{
		//				Scheme: "socks5",
		//				Host:   "192.168.1.29:2222"}, nil
		//		},
		//	},
		//},
		Endpoint:      "https://www.okex.com",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	})
}

func init() {
	log.SetLevel(log.DEBUG)
}

func TestOKExV5Spot_GetTicker(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetTicker(goex.BTC_USDT))
}

func TestOKExV5Spot_GetDepth(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetDepth(5, goex.BTC_USDT))
}

func TestOKExV5SpotGetKlineRecords(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 10))
}

func TestOKExV5Spot_LimitBuy(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.LimitBuy("1", "1.0", goex.XRP_USDT))
	//{"code":"0","data":[{"clOrdId":"0bf60374efe445BC258eddf46df044c3","ordId":"305267682086109184","sCode":"0","sMsg":"","tag":""}],"msg":""}}
}

func TestOKExV5Spot_CancelOrder(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.CancelOrder("305267682086109184", goex.XRP_USDT))
}

func TestOKExV5Spot_GetUnfinishOrders(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetUnfinishOrders(goex.XRP_USDT))
}

func TestOKExV5Spot_GetOneOrder(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetOneOrder("305267682086109184", goex.XRP_USDT))
}

func TestOKExV5Spot_GetAccount(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetAccount())
}
