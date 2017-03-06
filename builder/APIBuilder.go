package builder

import (
	. "github.com/nntaoli/crypto_coin_api"
	"github.com/nntaoli/crypto_coin_api/chbtc"
	"github.com/nntaoli/crypto_coin_api/huobi"
	"github.com/nntaoli/crypto_coin_api/okcoin"
	"github.com/nntaoli/crypto_coin_api/poloniex"
	"github.com/nntaoli/crypto_coin_api/yunbi"
	"net/http"
)

type APIBuilder struct {
	apiKey    string
	secretkey string
}

func NewAPIBuilder() (builder *APIBuilder) {
	return &APIBuilder{}
}

func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}

func (builder *APIBuilder) APISecretkey(key string) (_builder *APIBuilder) {
	builder.secretkey = key
	return builder
}

func (builder *APIBuilder) Build(exName string) (api API) {
	var _api API
	switch exName {
	case "okcoin.cn":
		_api = okcoin.New(http.DefaultClient, builder.apiKey, builder.secretkey)
	case "huobi.com":
		_api = huobi.New(http.DefaultClient, builder.apiKey, builder.secretkey)
	case "chbtc.com":
		_api = chbtc.New(http.DefaultClient, builder.apiKey, builder.secretkey)
	case "yunbi.com":
		_api = yunbi.New(http.DefaultClient, builder.apiKey, builder.secretkey)
	case "poloniex.com":
		_api = poloniex.New(http.DefaultClient, builder.apiKey, builder.secretkey)
	case "okcoin.com":
		_api = okcoin.NewCOM(http.DefaultClient, builder.apiKey, builder.secretkey)

	}
	return _api
}
