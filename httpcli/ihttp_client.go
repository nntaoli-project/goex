package httpcli

type IHttpClient interface {
	SetTimeout(sec int64)
	SetProxy(proxy string) error
	SetHeaders(key, value string) //添加全局http header
	DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error)
}
