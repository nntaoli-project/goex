package unit

import (
	"testing"
	"github.com/nntaoli/crypto_coin_api/poloniex"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"os"
)

const accesskey = ""
const secretkey = ""

func Test_Poloniex_GetTicker(t *testing.T) {
	//pUrl, _ := url.Parse("http://127.0.0.1:8787");
	//client := &http.Client{Transport:&http.Transport{Proxy:http.ProxyURL(pUrl)}};
	//client.Timeout = 30 * time.Second;
	//client.Timeout = 10 * time.Second;
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	ticker, err := poloniex.GetTicker(coinapi.XCN_BTC);
	if !assert.Empty(t, err) {
		os.Exit(-1);
	}
	t.Logf("last:%.8f", ticker.Last);
}

func Test_Poloniex_GetOrderBook(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	orderbook, err := poloniex.GetDepth(3, coinapi.ETC_BTC);
	assert.Empty(t, err);
	t.Log(orderbook);
}

func Test_Poloniex_GetAccount(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	acc, err := poloniex.GetAccount();
	assert.Empty(t, err);
	t.Log(acc);
}

func Test_Poloniex_GetUnfinishOrders(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	orders, err := poloniex.GetUnfinishOrders(coinapi.ETC_BTC);
	assert.Empty(t, err);
	t.Log(orders);
}

func Test_Poloniex_LimitPlaceOrder(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	orid, err := poloniex.LimitBuy("0.08", "0.00194902", coinapi.ETC_BTC);
	assert.Empty(t, err);
	t.Log(orid)
}

func Test_Poloniex_CancelOrder(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	r, err := poloniex.CancelOrder("17523364266", coinapi.ETC_BTC);
	assert.Empty(t, err);
	assert.True(t, r);
}

func Test_Poloniex_GetOneOrder(t *testing.T) {
	os.Setenv("http_proxy", "http://127.0.0.1:8787");
	poloniex := poloniex.New(http.DefaultClient, accesskey, secretkey);
	order, err := poloniex.GetOneOrder("17523364266", coinapi.ETC_BTC);
	assert.Empty(t, err);
	t.Log(order);
}
