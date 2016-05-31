package unit

import (
	"testing"
	"rest/huobi"
	. "rest"
	"net/http"
)

func Test_HuoBiApi(t *testing.T) {
	client := http.DefaultClient;
	huobiApi := huobi.New(client, "", "");
	ticker, err := huobiApi.GetTicker(LTC_CNY);
	if err != nil {
		t.Error(err);
		return;
	}

	t.Log(ticker);
}
