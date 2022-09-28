package goex

import (
	"github.com/nntaoli-project/goex/v2/internal/config"
	"github.com/nntaoli-project/goex/v2/internal/lib"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/url"
	"time"
)

var httpCli IHttpClient

func SetHttpProxy(uri string) error {
	proxyU, err := url.Parse(uri)
	if err != nil {
		return err
	}
	config.C.HttpConf.Proxy = proxyU
	return nil
}

func SetHttpTimeout(sec int) {
	config.C.HttpConf.Timeout = time.Duration(sec) * time.Second
}

func SetHttpCli(cli IHttpClient) {
	httpCli = cli
}

func GetHttpCli() IHttpClient {
	return httpCli
}

func SetupDefaultLibs() {
	SetHttpCli(lib.NewDefaultHttpClient())
}

func SetDebugLogger() {
	logger.Log.SetLevel(logger.DEBUG)
}
