package config

import (
	"net/url"
	"time"
)

type httpConf struct {
	Proxy   *url.URL
	Timeout time.Duration
}

type conf struct {
	HttpConf *httpConf
}

var C = &conf{
	HttpConf: &httpConf{
		Timeout: 5 * time.Second,
	},
}
