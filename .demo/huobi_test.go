package demo

import (
	"fmt"
	"net/http"
	"testing"

	goex ".."
	huobi "../huobi"
	"github.com/stretchr/testify/assert"
)

var hbpro = huobi.NewHuobiPro(http.DefaultClient, HUOBI_API_KEY, HUOBI_SECRET_KEY, "")

var pair = goex.BCX_BTC

func TestHuoBi_Pro_GetDepth(t *testing.T) {
	depth, err := hbpro.GetDepth(50, pair)
	assert.Nil(t, err)
	t.Log("asks: ", depth.AskList)
	t.Log("bids: ", depth.BidList)
}
func TestHuoBi_Pro_LimitSell(t *testing.T) {
	accountid, err := hbpro.GetAccountId()
	fmt.Printf("accountid: %s err:%v", accountid, err)
	// account, err := hbpro.GetAccount()
	// fmt.Printf("account: %v err:%v", account, err)

	order, err := hbpro.LimitSell("10.0", "0.00101188", pair)
	assert.Nil(t, err)
	t.Logf("order: %v", order)
	works, err := hbpro.CancelOrder(order.OrderID2, pair)

	fmt.Printf("CancelOrder: %s  err:%v", works, err)

}
