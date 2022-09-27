package goex

import (
	"github.com/nntaoli-project/goex/v2/internal/config"
	"net/url"
	"time"
)

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

func SetHttpLib(lib string) {
	config.C.HttpConf.Lib = lib
}
