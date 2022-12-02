package lib

import (
	"context"
	"errors"
	"fmt"
	"github.com/nntaoli-project/goex/v2/internal/config"
	"github.com/nntaoli-project/goex/v2/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DefaultHttpClient struct {
	cli *http.Client
}

func NewDefaultHttpClient() *DefaultHttpClient {
	cli := new(DefaultHttpClient)
	cli.init()
	return cli
}

func (cli *DefaultHttpClient) init() {
	logger.Info("[http utils] setup lib http cli")

	httpTimeout := config.C.HttpConf.Timeout
	cli.cli = &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			IdleConnTimeout: 2 * httpTimeout,
		},
	}

	proxyUrl := config.C.HttpConf.Proxy
	if proxyUrl != nil {
		logger.Info("[http utils] proxy=", proxyUrl.String())
		cli.cli.Transport.(*http.Transport).Proxy = func(r *http.Request) (*url.URL, error) {
			return proxyUrl, nil
		}
	}
}

func (cli *DefaultHttpClient) DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error) {
	logger.Debugf("[http utils] [%s] request url: %s", method, rqUrl)

	reqTimeoutCtx, _ := context.WithTimeout(context.TODO(), config.C.HttpConf.Timeout)
	req, _ := http.NewRequestWithContext(reqTimeoutCtx, method, rqUrl, strings.NewReader(reqBody))

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := cli.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("[http utils] close response body error:", err.Error())
		}
	}(resp.Body)

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode != 200 {
		return bodyData, errors.New(resp.Status)
	}

	return bodyData, nil
}
