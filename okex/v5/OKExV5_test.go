package okex

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/nntaoli-project/goex"
)

func newOKExV5Client() *OKExV5 {
	return NewOKExV5(&goex.APIConfig{
		//HttpClient: &http.Client{
		//	Transport: &http.Transport{
		//		Proxy: func(req *http.Request) (*url.URL, error) {
		//			return &url.URL{
		//				Scheme: "socks5",
		//				Host:   "127.0.0.1:2222"}, nil
		//		},
		//	},
		//},
		Endpoint:      "https://www.okex.com",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	})
}

func TestOKExV5_GetTicker(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetTickerV5("BTC-USD-SWAP"))
}

func TestOKExV5_GetDepth(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetDepthV5("BTC-USD-SWAP", 0))
}

func TestOKExV5_GetKlineRecordsV5(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetKlineRecordsV5("BTC-USD-SWAP", goex.KLINE_PERIOD_1H, &url.Values{}))

}
