package util

//http request 工具函数
import (
	"context"
	"errors"
	"fmt"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/config"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	fastHttpClient *fasthttp.Client
	socksDialer    fasthttp.DialFunc

	defaultHttpClient *http.Client
)

func setupFastHttpClient() {
	if fastHttpClient == nil {
		logger.Log.Info("[http utils] init fast http client")
		httpTimeout := config.C.HttpConf.Timeout
		fastHttpClient = &fasthttp.Client{
			Name:               "goex-http-utils",
			MaxConnsPerHost:    512,
			MaxConnWaitTimeout: 4 * httpTimeout,
			WriteTimeout:       httpTimeout,
			ReadTimeout:        httpTimeout,
		}

		proxyUrl := config.C.HttpConf.Proxy
		if proxyUrl != nil && socksDialer == nil {
			if proxyUrl.Scheme != "socks5" {
				logger.Log.Error("[http utils] fasthttp only support the socks5 proxy")
				return
			}
			logger.Log.Info("[http utils] proxy=", proxyUrl.String())
			socksDialer = fasthttpproxy.FasthttpSocksDialer(proxyUrl.String())
			fastHttpClient.Dial = socksDialer
		}
	}

}

func doHttpRequestWithFasthttp(reqMethod, reqUrl, postData string, headers map[string]string) ([]byte, error) {
	logger.Log.Debug("[http utils] use fasthttp client")

	setupFastHttpClient()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.SetMethod(reqMethod)
	req.SetRequestURI(reqUrl)
	req.SetBodyString(postData)

	err := fastHttpClient.DoTimeout(req, resp, config.C.HttpConf.Timeout)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode(), string(resp.Body())))
	}

	return resp.Body(), nil
}

func setupDefaultHttpClient() {
	if defaultHttpClient == nil {
		logger.Log.Info("[http utils] init default http client")
		httpTimeout := config.C.HttpConf.Timeout
		defaultHttpClient = &http.Client{
			Timeout: httpTimeout,
			Transport: &http.Transport{
				IdleConnTimeout:       2 * httpTimeout,
				ResponseHeaderTimeout: httpTimeout,
				ExpectContinueTimeout: httpTimeout,
			},
		}

		proxyUrl := config.C.HttpConf.Proxy
		if proxyUrl != nil {
			logger.Log.Info("[http utils] proxy=", proxyUrl.String())
			defaultHttpClient.Transport.(*http.Transport).Proxy = func(r *http.Request) (*url.URL, error) {
				return proxyUrl, nil
			}
		}
	}
}

func DoHttpRequest(reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	logger.Log.Debugf("[%s] request url: %s", reqType, reqUrl)

	lib := os.Getenv("HTTP_LIB")
	if lib == goex.HttpLib_FastHttpClient || config.C.HttpConf.Lib == goex.HttpLib_FastHttpClient {
		return doHttpRequestWithFasthttp(reqType, reqUrl, postData, requstHeaders)
	}

	setupDefaultHttpClient()

	reqTimeoutCtx, _ := context.WithTimeout(context.TODO(), config.C.HttpConf.Timeout)
	req, _ := http.NewRequestWithContext(reqTimeoutCtx, reqType, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}

	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := defaultHttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}
