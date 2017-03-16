package builder

import (
	. "github.com/nntaoli/crypto_coin_api"
	"github.com/nntaoli/crypto_coin_api/chbtc"
	"github.com/nntaoli/crypto_coin_api/coincheck"
	"context"
	"github.com/nntaoli/crypto_coin_api/huobi"
	"github.com/nntaoli/crypto_coin_api/okcoin"
	"github.com/nntaoli/crypto_coin_api/poloniex"
	"github.com/nntaoli/crypto_coin_api/yunbi"
	"github.com/nntaoli/crypto_coin_api/zaif"
	"github.com/qiniu/log"
	"net"
	"net/http"
	"time"
)

type APIBuilder struct {
	client      *http.Client
	httpTimeout time.Duration
	apiKey      string
	secretkey   string
}

func NewAPIBuilder() (builder *APIBuilder) {
	_client := http.DefaultClient
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 4 * time.Second,
	}
	_client.Transport = transport
	return &APIBuilder{client: _client}
}

func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}

func (builder *APIBuilder) APISecretkey(key string) (_builder *APIBuilder) {
	builder.secretkey = key
	return builder
}

func (builder *APIBuilder) HttpTimeout(timeout time.Duration) (_builder *APIBuilder) {
	builder.httpTimeout = timeout
	builder.client.Timeout = timeout
	transport := builder.client.Transport.(*http.Transport)
	if transport != nil {
		transport.ResponseHeaderTimeout = timeout
		transport.TLSHandshakeTimeout = timeout
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}
	return builder
}

func (builder *APIBuilder) Build(exName string) (api API) {
	var _api API
	switch exName {
	case "okcoin.cn":
		_api = okcoin.New(builder.client, builder.apiKey, builder.secretkey)
	case "huobi.com":
		_api = huobi.New(builder.client, builder.apiKey, builder.secretkey)
	case "chbtc.com":
		_api = chbtc.New(builder.client, builder.apiKey, builder.secretkey)
	case "yunbi.com":
		_api = yunbi.New(builder.client, builder.apiKey, builder.secretkey)
	case "poloniex.com":
		_api = poloniex.New(builder.client, builder.apiKey, builder.secretkey)
	case "okcoin.com":
		_api = okcoin.NewCOM(builder.client, builder.apiKey, builder.secretkey)
	case "coincheck.com":
		_api = coincheck.New(builder.client, builder.apiKey, builder.secretkey)
	case "zaif.jp":
		_api = zaif.New(builder.client, builder.apiKey, builder.secretkey)
	default:
		log.Println("error")

	}
	return _api
}
