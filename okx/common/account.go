package common

import (
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"net/http"
	"net/url"
)

type AccountApi struct {
	*V5
	apiOpts ApiOptions
}

func NewAccountApi(apiOpts ApiOptions) *AccountApi {
	return &AccountApi{
		apiOpts: apiOpts,
	}
}

func (acc *AccountApi) GetAccount(coin string) (map[string]Account, error) {
	reqUrl := fmt.Sprintf("%s%s", acc.uriOpts.Endpoint, acc.uriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, err := acc.V5.DoAuthRequest(http.MethodGet, reqUrl, &params, acc.apiOpts, nil)
	if err != nil {
		return nil, err
	}
	return acc.unmarshalOpts.GetAccountResponseUnmarshaler(data)
}
