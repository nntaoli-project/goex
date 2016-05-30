package unit

import (
	"testing"
	"rest/huobi"
	. "rest"
	"net/http"
)

func Test_HuoBiApi(t *testing.T)  {
	client := http.DefaultClient;
	hbCfg := APIConfig{client ,"http://api.huobi.com/","",""};
	huobiApi := huobi.New(hbCfg);
	ticker , err := huobiApi.GetTicker(LTC_CNY);
	if err != nil {
		t.Error(err);
		return ;
	}

	t.Log(ticker);
}
