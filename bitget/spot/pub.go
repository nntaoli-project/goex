package spot

import (
	"errors"
	. "github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
	"net/http"
	"net/url"
)

func (s *Spot) NewCurrencyPair(baseSym, quoteSym string) (model.CurrencyPair, error) {
	currencyPair := s.currencyPairM[baseSym+quoteSym]
	if currencyPair.Symbol == "" {
		return currencyPair, errors.New("not found currency pair")
	}
	return currencyPair, nil
}

func (s *Spot) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	body, err := s.DoNoAuthRequest(http.MethodGet, s.UriOpts.Endpoint+s.UriOpts.GetExchangeInfoUri, &url.Values{}, nil)
	if err != nil {
		logger.Errorf("[GetExchangeInfo] http request error, body: %s", string(body))
		return nil, body, err
	}

	m, err := s.UnmarshalerOpts.GetExchangeInfoResponseUnmarshaler(body)
	if err != nil {
		logger.Errorf("[GetExchangeInfo] unmarshaler data error, err: %s", err.Error())
		return nil, body, err
	}

	s.currencyPairM = m

	return m, body, err
}

func (s *Spot) DoNoAuthRequest(method, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	var reqBody string

	if method == http.MethodGet {
		reqUrl += "?" + params.Encode()
	} else {
		reqBody = params.Encode()
	}

	responseData, err := Cli.DoRequest(method, reqUrl, reqBody, headers)
	if err != nil {
		return responseData, err
	}

	return responseData, err
}
