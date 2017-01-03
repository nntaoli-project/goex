package unit
import (
	"testing"
	"github.com/nntaoli/crypto_coin_api/huobi"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	hb *huobi.HuoBi = huobi.New(http.DefaultClient, "", "")
)

func Test_Huobi_GetKline(t *testing.T) {
	klines, _ := hb.GetKlineRecords(coinapi.BTC_CNY, "1", 100, -1);
	for _, v := range klines {
		t.Log(time.Unix(v.Timestamp, 0));
	}
}

func Test_Huobi_GetOrder(t *testing.T) {
	ord, err := hb.GetOneOrder("3399439745", coinapi.BTC_CNY)
	assert.NoError(t, err)
	t.Log(ord)
}

func Test_Huobi_Market(t *testing.T)  {
	ord ,err := hb.MarketBuy("0" , "1" , coinapi.LTC_CNY)
	assert.NoError(t , err)
	t.Log(ord)

	ord2 , err := hb.MarketSell("0.02" , "0" , coinapi.LTC_CNY)
	assert.NoError(t , err)
	t.Log(ord2)
}
