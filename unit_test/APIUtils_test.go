package unit_test

import (
	. "github.com/nntaoli/crypto_coin_api"
	"github.com/nntaoli/crypto_coin_api/builder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallAllUnfinishedOrders(t *testing.T) {
	api := builder.NewAPIBuilder().
		APIKey("aaa").
		APISecretkey("bbb").
		Build("huobi.com")

	t.Log(api.GetExchangeName())

	c := CancelAllUnfinishedOrders(api, LTC_CNY)
	t.Logf("cancel order count [%d]" , c)

	//verify
	orders, _ := api.GetUnfinishOrders(LTC_CNY)
	assert.True(t, len(orders) == 0)
}
