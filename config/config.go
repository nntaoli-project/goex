package config

import (
	"net/url"
	"time"
)

type conf struct {
	httpConf struct {
		proxy   *url.URL
		timeout time.Duration
	}
}

var _conf = new(conf)

func init() {
	_conf.httpConf.timeout = 5 * time.Second
}

func SetHttpProxy(uri string) error {
	proxyU, err := url.Parse(uri)
	if err != nil {
		return err
	}
	_conf.httpConf.proxy = proxyU
	return nil
}

func SetHttpTimeout(sec int) {
	_conf.httpConf.timeout = time.Duration(sec) * time.Second
}

func GetHttpTimeout() time.Duration {
	return _conf.httpConf.timeout
}

func GetHttpProxy() *url.URL {
	return _conf.httpConf.proxy
}
