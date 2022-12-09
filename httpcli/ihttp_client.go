package httpcli

type IHttpClient interface {
	SetTimeout(sec int64)
	SetProxy(proxy string) error
	DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error)
}
