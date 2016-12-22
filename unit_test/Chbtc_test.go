package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nntaoli/crypto_coin_api/chbtc"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
)

func Test_GetExchangeName(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "", "")
	assert.Equal(t, "chbtc", chbtc.GetExchangeName())
}

func Test_Chbtc_GetTicker(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "", "")
	ticker, err := chbtc.GetTicker(coinapi.ETH_BTC);
	assert.Empty(t, err)
	t.Log(ticker)
}

func Test_Chbtc_GetDepth(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "", "")
	depth, err := chbtc.GetDepth(1, coinapi.BTC_CNY)
	assert.Empty(t, err)
	t.Log(depth.BidList[0].Amount , depth.BidList[0].Price)
}

func Test_Chbtc_GetAcoount(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "accesskey", "secretkey")
	acc, err := chbtc.GetAccount();
	assert.Empty(t, err)
	t.Log(acc)
}

func Test_Chbtc_GetOneOrder(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "accesskey", "secretkey")
	order, err := chbtc.GetOneOrder("201609235068899", coinapi.ETC_CNY);
	assert.Empty(t, err);
	t.Log(order)
}

func Test_Chbtc_GetUnfinishedOrder(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "accesskey", "secretkey")
	order, err := chbtc.GetUnfinishOrders(coinapi.ETC_CNY);
	assert.Empty(t, err);
	t.Log(order)
}

func Test_Chbtc_CancelOrder(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "accesskey", "secretkey")
	ret, err := chbtc.CancelOrder("201609235077252", coinapi.ETC_CNY);
	//t.Log(err == errors.New("3001"))
	//assert.Equal(t, err, errors.New("3001"));
	assert.Empty(t, err);
	assert.True(t, ret);
}

func Test_Chbtc_PlaceOrder(t *testing.T) {
	chbtc := chbtc.New(http.DefaultClient, "accesskey", "secretkey")
	order, err := chbtc.LimitBuy("0.1", "8.0", coinapi.ETC_CNY);
	assert.Empty(t, err);
	t.Log(order);
}
