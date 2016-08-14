package unit

import (
	"testing"
	"../haobtc"
	"net/http"
	. "../"
	"github.com/stretchr/testify/assert"
//	"encoding/json"
)

func Test_GetTicker(t *testing.T) {
	api := haobtc.New(http.DefaultClient, "", "");
	_, err := api.GetTicker(LTC_CNY);
	assert.Error(t, err);

	ticker, err2 := api.GetTicker(BTC_CNY);
	assert.NoError(t, err2);
	t.Log(ticker);
}

func Test_GetDepth(t *testing.T) {
	api := haobtc.New(http.DefaultClient, "", "");
	_, err := api.GetDepth(5, LTC_CNY);
	assert.Error(t, err);

	depth, err2 := api.GetDepth(3, BTC_CNY);
	assert.NoError(t, err2);
	t.Log(depth);
}

func Test_GetAcount(t *testing.T) {
	api := haobtc.New(http.DefaultClient, "apikey", "secretkey");

	account, err := api.GetAccount();
	assert.NoError(t, err);

	t.Log(account);
}

func Test_LimitSell(t *testing.T)  {
	api := haobtc.New(http.DefaultClient, "apikey", "secretkey");

	order, err := api.LimitSell("0.01" , "4100" , BTC_CNY);
	assert.NoError(t, err);

	t.Log(order);
}

func  Test_HaoBtc_CancelOrder(t *testing.T){
	api := haobtc.New(http.DefaultClient, "apikey", "secret_key");
	ret , err := api.CancelOrder("-1" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(ret);
}

func Test_HaoBtc_GetOnOrder(t *testing.T)  {
	api := haobtc.New(http.DefaultClient, "", "");
	ret , err := api.GetOneOrder("123" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(ret);
}

func Test_HaoBtc_GetUnfinishOrder(t *testing.T)  {
	api := haobtc.New(http.DefaultClient, "", "");
	ret , err := api.GetUnfinishOrders(BTC_CNY);
	assert.NoError(t , err);
	t.Log(ret);
}
