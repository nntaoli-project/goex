package config

import (
	"net/url"
	"time"
)

const (
	Lib_FastHttpClient = "fasthttp"
)

type httpConf struct {
	lib     string
	proxy   *url.URL
	timeout time.Duration
}

type conf struct {
	httpC *httpConf
}

var _conf = &conf{
	httpC: &httpConf{
		lib:     "default",
		timeout: 5 * time.Second,
	},
}

func SetHttpProxy(uri string) error {
	proxyU, err := url.Parse(uri)
	if err != nil {
		return err
	}
	_conf.httpC.proxy = proxyU
	return nil
}

func SetHttpTimeout(sec int) {
	_conf.httpC.timeout = time.Duration(sec) * time.Second
}

func SetHttpLib(lib string) {
	_conf.httpC.lib = lib
}

func GetHttpTimeout() time.Duration {
	return _conf.httpC.timeout
}

func GetHttpProxy() *url.URL {
	return _conf.httpC.proxy
}

func GetHttpLib() string {
	return _conf.httpC.lib
}
