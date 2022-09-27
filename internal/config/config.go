package config

import (
	"net/url"
	"time"
)

type httpConf struct {
	Lib     string
	Proxy   *url.URL
	Timeout time.Duration
}

type conf struct {
	HttpConf *httpConf
}

var C = &conf{
	HttpConf: &httpConf{
		Lib:     "default",
		Timeout: 5 * time.Second,
	},
}
