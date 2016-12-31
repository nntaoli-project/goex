package unit

import (
	"testing"
	"github.com/nntaoli/crypto_coin_api/huobi"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
	"time"
)

func Test_Huobi_GetKline(t *testing.T) {
	var hb coinapi.API;
	hb = huobi.New(http.DefaultClient, "", "");
	hb.GetTicker(coinapi.BTC_CNY)
	klines, _ := hb.GetKlineRecords(coinapi.BTC_CNY, "1", 100, -1);
	for _, v := range klines {
		t.Log(time.Unix(v.Timestamp, 0));
	}

	loc , _ := time.LoadLocation("Local");
	t.Log(loc)
}
