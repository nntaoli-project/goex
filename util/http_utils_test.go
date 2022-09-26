package util

import (
	"github.com/nntaoli-project/goex/v2/internal/config"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"testing"
)

func TestDoHttpRequest(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	config.SetHttpProxy("socks5://127.0.0.1:2220") //local proxy
	//	config.SetHttpLib(config.Lib_FastHttpClient)
	config.SetHttpTimeout(5)

	ret, err := DoHttpRequest(http.MethodGet, "https://www.google.com", "", map[string]string{})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(ret))
}
