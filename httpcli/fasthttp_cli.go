package httpcli

import (
	"errors"
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"time"
)

type FastHttpCli struct {
	fastHttpClient *fasthttp.Client
	//socksDialer    fasthttp.DialFunc
	timeout time.Duration
}

func NewFastHttpCli() *FastHttpCli {
	cli := &fasthttp.Client{
		Name:                "github.com/nntaoli-project/goex/v2/",
		MaxConnsPerHost:     512,
		MaxIdleConnDuration: 20 * time.Second,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
	}
	return &FastHttpCli{fastHttpClient: cli, timeout: 10 * time.Second}
}

func (cli *FastHttpCli) SetTimeout(sec int64) {
	cli.timeout = time.Duration(sec) * time.Second
	cli.fastHttpClient.WriteTimeout = cli.timeout
	cli.fastHttpClient.ReadTimeout = cli.timeout
}

func (cli *FastHttpCli) SetProxy(proxy string) error {
	logger.Infof("[fast http cli] proxy=%s", proxy)
	cli.fastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
	return nil
}

func (cli *FastHttpCli) DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error) {
	//logger.Info("[fast http cli] use fasthttp client")
	logger.Debug("[fast http cli]  req url:", rqUrl)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.SetMethod(method)
	req.SetRequestURI(rqUrl)
	req.SetBodyString(reqBody)

	err = cli.fastHttpClient.DoTimeout(req, resp, cli.timeout)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New(resp.String())
	}

	// 拷贝响应的 body
	responseBody := make([]byte, len(resp.Body()))
	copy(responseBody, resp.Body())

	return responseBody, nil
}
