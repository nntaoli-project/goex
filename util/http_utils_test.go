package util

import (
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"testing"
)

func TestDoHttpRequest(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	//goex.SetHttpProxy("socks5://127.0.0.1:2220") //local proxy
	goex.SetHttpLib(goex.HttpLib_FastHttpClient)
	goex.SetHttpTimeout(5)

	ret, err := DoHttpRequest(http.MethodGet, "http://www.baidu.com", "", map[string]string{})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(ret))
}
