package goex

import "net/url"

type IHttpClient interface {
	DoRequest(method, rqUrl string, params *url.Values, headers map[string]string) (data []byte, err error)
}
