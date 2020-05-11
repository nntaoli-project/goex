package coinbig

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var cb = New(http.DefaultClient, "", "")


func TestCoinBig_BuildSigned(t *testing.T) {
	return
	//params := url.Values{}
	//
	//params.Set("apikey", cb.accessKey)
	//params.Set("symbol", "btc_usdt")
	//params.Set("size", "20")
	//params.Set("type", "1")
	//t.Log(params)
	//t.Log(cb.Encode(params))
	//cb.buildSigned(&params)
	//t.Log(cb.Encode(params))
	//0272D6FAF94868C0E84B9FC86C1266AA
	//E94B142B9A11DE390F20D8AB193D7FE1
	param := "apikey=7E7D227B7CA207088BF275BACB3015B2&size=20&symbol=btc_usdt&type=1&secret_key=4AC5140A7CDDC7F9610A3C0E3284231C"

	//w := md5.New()
	//io.WriteString(w,param)
	//md5String := fmt.Sprintf("%x",w.Sum(nil))
	//t.Log(strings.ToUpper(md5String))
	s, _ := goex.GetParamMD5Sign("", param)
	t.Log(s)
}

func TestCoinBig_GetAccount(t *testing.T) {
	//return
	t.Log(cb.GetAccount())
}
func TestCoinBig_GetUnfinishOrders(t *testing.T) {
	//return
	t.Log(cb.GetUnfinishOrders(goex.BTC_USDT))
}
func TestCoinBig_GetOneOrder(t *testing.T) {
	return
	t.Log(cb.GetOneOrder("1111", goex.BTC_USDT))
}

func TestCoinBig_CancelOrder(t *testing.T) {
	return
	t.Log(cb.CancelOrder("1111", goex.BTC_USDT))

}
func TestCoinBig_GetTicker(t *testing.T) {
	return
	t.Log(cb.GetTicker(goex.BTC_USDT))
}

func TestCoinBig_GetDepth(t *testing.T) {
	return
	t.Log(cb.GetDepth(3, goex.BTC_USDT))
}

func TestCoinBig_LimitBuy(t *testing.T) {
	return
	t.Log(cb.LimitBuy("1", "1", goex.BTC_USDT))
}
