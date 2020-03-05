package goex

//http request 工具函数
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nntaoli-project/goex/internal/logger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var fastHttpClient fasthttp.Client

func init() {
	fastHttpClient.MaxConnsPerHost = 2
	fastHttpClient.ReadTimeout = 10 * time.Second
	fastHttpClient.WriteTimeout = 10 * time.Second
}

func NewHttpRequestWithFasthttp(client *http.Client, reqMethod, reqUrl, postData string, headers map[string]string) ([]byte, error) {
	logger.Log.Debug("use fasthttp client")
	transport := client.Transport
	if transport != nil {
		if proxy, err := transport.(*http.Transport).Proxy(nil); err == nil && proxy != nil {
			proxyUrl := proxy.String()
			logger.Log.Debug("proxy url: ", proxyUrl)
			if proxy.Scheme != "socks5" {
				logger.Log.Error("fasthttp only support the socks5 proxy")
			} else {
				fastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(strings.TrimPrefix(proxyUrl, proxy.Scheme+"://"))
			}
		}
	}
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
	err := fastHttpClient.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode(), string(resp.Body())))
	}

	return resp.Body(), nil
}

func NewHttpRequest(client *http.Client, reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	logger.Log.Debugf("[%s] request url: %s", reqType, reqUrl)
	lib := os.Getenv("HTTP_LIB")
	if lib == "fasthttp" {
		return NewHttpRequestWithFasthttp(client, reqType, reqUrl, postData, requstHeaders)
	}

	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
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

func HttpGet(client *http.Client, reqUrl string) (map[string]interface{}, error) {
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", nil)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println(string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpGet2(client *http.Client, reqUrl string, headers map[string]string) (map[string]interface{}, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println("respData", string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpGet3(client *http.Client, reqUrl string, headers map[string]string) ([]interface{}, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	//println(string(respData))
	var bodyDataMap []interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println("respData", string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpGet4(client *http.Client, reqUrl string, headers map[string]string, result interface{}) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respData, result)
	if err != nil {
		log.Printf("HttpGet4 - json.Unmarshal failed : %v, resp %s", err, string(respData))
		return err
	}

	return nil
}
func HttpGet5(client *http.Client, reqUrl string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostForm2(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostForm3(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}

func HttpPostForm4(client *http.Client, reqUrl string, postData map[string]string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"
	data, _ := json.Marshal(postData)
	return NewHttpRequest(client, "POST", reqUrl, string(data), headers)
}

func HttpDeleteForm(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
}

func HttpPut(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "PUT", reqUrl, postData.Encode(), headers)
}
