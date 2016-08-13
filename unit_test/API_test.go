package unit

import (
	"testing"
	"../huobi"
	. "../"
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

	depth, err := huobiApi.GetDepth(3, LTC_CNY);
	t.Log("bids:", depth.AskList);
	t.Log("asks:", depth.AskList);
}
