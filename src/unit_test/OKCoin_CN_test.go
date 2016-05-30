package unit

import (
	. "rest"
	"github.com/stretchr/testify/assert"
	"testing"
    "rest/okcoin"
)

func Test_OKCoin_CN(t *testing.T) {
    var api API;
    api = okcoin.New("okcoin_cn", "", "");
    tk, err := api.GetTicker(BTC_CNY);
    assert.True(t, err == nil);
    t.Logf("last:%s buy:%s sell:%s high:%s low:%s vol:%s date:%s",
        tk.Last, tk.Buy, tk.Sell, tk.High, tk.Low, tk.Vol, tk.Date);
}