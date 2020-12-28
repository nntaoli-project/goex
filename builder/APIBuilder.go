package builder

import (
	"context"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/bigone"
	"github.com/nntaoli-project/goex/binance"
	"github.com/nntaoli-project/goex/bitfinex"
	"github.com/nntaoli-project/goex/bithumb"
	"github.com/nntaoli-project/goex/bitmex"
	"github.com/nntaoli-project/goex/bitstamp"
	"github.com/nntaoli-project/goex/bittrex"
	"github.com/nntaoli-project/goex/coinbene"
	"github.com/nntaoli-project/goex/kucoin"

	"github.com/nntaoli-project/goex/atop"
	//"github.com/nntaoli-project/goex/coin58"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/nntaoli-project/goex/coinex"
	"github.com/nntaoli-project/goex/gdax"
	"github.com/nntaoli-project/goex/hitbtc"
	"github.com/nntaoli-project/goex/huobi"
	"github.com/nntaoli-project/goex/kraken"
	"github.com/nntaoli-project/goex/okex"
	"github.com/nntaoli-project/goex/poloniex"
	"github.com/nntaoli-project/goex/zb"
)

type APIBuilder struct {
	HttpClientConfig *HttpClientConfig
	client           *http.Client
	httpTimeout      time.Duration
	apiKey           string
	secretkey        string
	clientId         string
	apiPassphrase    string
	futuresEndPoint  string
	endPoint         string
}

type HttpClientConfig struct {
	HttpTimeout  time.Duration
	Proxy        *url.URL
	MaxIdleConns int
}

func (c HttpClientConfig) String() string {
	return fmt.Sprintf("{ProxyUrl:\"%s\",HttpTimeout:%s,MaxIdleConns:%d}", c.Proxy, c.HttpTimeout.String(), c.MaxIdleConns)
}

func (c *HttpClientConfig) SetHttpTimeout(timeout time.Duration) *HttpClientConfig {
	c.HttpTimeout = timeout
	return c
}

func (c *HttpClientConfig) SetProxyUrl(proxyUrl string) *HttpClientConfig {
	if proxyUrl == "" {
		return c
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return c
	}
	c.Proxy = proxy
	return c
}

func (c *HttpClientConfig) SetMaxIdleConns(max int) *HttpClientConfig {
	c.MaxIdleConns = max
	return c
}

var (
	DefaultHttpClientConfig = &HttpClientConfig{
		Proxy:        nil,
		HttpTimeout:  5 * time.Second,
		MaxIdleConns: 10}
	DefaultAPIBuilder = NewAPIBuilder()
)

func NewAPIBuilder() (builder *APIBuilder) {
	return NewAPIBuilder2(DefaultHttpClientConfig)
}

func NewAPIBuilder2(config *HttpClientConfig) *APIBuilder {
	if config == nil {
		config = DefaultHttpClientConfig
	}

	return &APIBuilder{
		HttpClientConfig: config,
		client: &http.Client{
			Timeout: config.HttpTimeout,
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return config.Proxy, nil
				},
				MaxIdleConns:          config.MaxIdleConns,
				IdleConnTimeout:       5 * config.HttpTimeout,
				MaxConnsPerHost:       2,
				MaxIdleConnsPerHost:   2,
				TLSHandshakeTimeout:   config.HttpTimeout,
				ResponseHeaderTimeout: config.HttpTimeout,
				ExpectContinueTimeout: config.HttpTimeout,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, config.HttpTimeout)
				}},
		}}
}

func NewCustomAPIBuilder(client *http.Client) (builder *APIBuilder) {
	return &APIBuilder{client: client}
}

func (builder *APIBuilder) GetHttpClientConfig() *HttpClientConfig {
	return builder.HttpClientConfig
}

func (builder *APIBuilder) GetHttpClient() *http.Client {
	return builder.client
}

func (builder *APIBuilder) HttpProxy(proxyUrl string) (_builder *APIBuilder) {
	if proxyUrl == "" {
		return builder
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return builder
	}
	builder.HttpClientConfig.Proxy = proxy
	transport := builder.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxy)
	return builder
}

func (builder *APIBuilder) HttpTimeout(timeout time.Duration) (_builder *APIBuilder) {
	builder.HttpClientConfig.HttpTimeout = timeout
	builder.httpTimeout = timeout
	builder.client.Timeout = timeout
	transport := builder.client.Transport.(*http.Transport)
	if transport != nil {
		//transport.ResponseHeaderTimeout = timeout
		//transport.TLSHandshakeTimeout = timeout
		transport.IdleConnTimeout = timeout
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}
	return builder
}

func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}

func (builder *APIBuilder) APISecretkey(key string) (_builder *APIBuilder) {
	builder.secretkey = key
	return builder
}

func (builder *APIBuilder) ClientID(id string) (_builder *APIBuilder) {
	builder.clientId = id
	return builder
}

func (builder *APIBuilder) ApiPassphrase(apiPassphrase string) (_builder *APIBuilder) {
	builder.apiPassphrase = apiPassphrase
	return builder
}

func (builder *APIBuilder) FuturesEndpoint(endpoint string) (_builder *APIBuilder) {
	builder.futuresEndPoint = endpoint
	return builder
}

func (builder *APIBuilder) Endpoint(endpoint string) (_builer *APIBuilder) {
	builder.endPoint = endpoint
	return builder
}

func (builder *APIBuilder) Build(exName string) (api API) {
	var _api API
	switch exName {
	case KUCOIN:
		_api = kucoin.New(builder.apiKey, builder.secretkey, builder.apiPassphrase)
	//case OKCOIN_CN:
	//	_api = okcoin.New(builder.client, builder.apiKey, builder.secretkey)
	case POLONIEX:
		_api = poloniex.New(builder.client, builder.apiKey, builder.secretkey)
	//case OKCOIN_COM:
	//	_api = okcoin.NewCOM(builder.client, builder.apiKey, builder.secretkey)
	case BITSTAMP:
		_api = bitstamp.NewBitstamp(builder.client, builder.apiKey, builder.secretkey, builder.clientId)
	case HUOBI_PRO:
		//_api = huobi.NewHuoBiProSpot(builder.client, builder.apiKey, builder.secretkey)
		_api = huobi.NewHuobiWithConfig(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey})
	case OKEX_V3, OKEX:
		_api = okex.NewOKEx(&APIConfig{
			HttpClient:    builder.client,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretkey,
			ApiPassphrase: builder.apiPassphrase,
			Endpoint:      builder.endPoint,
		})
	case BITFINEX:
		_api = bitfinex.New(builder.client, builder.apiKey, builder.secretkey)
	case KRAKEN:
		_api = kraken.New(builder.client, builder.apiKey, builder.secretkey)
	case BINANCE:
		//_api = binance.New(builder.client, builder.apiKey, builder.secretkey)
		_api = binance.NewWithConfig(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey})
	case BITTREX:
		_api = bittrex.New(builder.client, builder.apiKey, builder.secretkey)
	case BITHUMB:
		_api = bithumb.New(builder.client, builder.apiKey, builder.secretkey)
	case GDAX:
		_api = gdax.New(builder.client, builder.apiKey, builder.secretkey)
	case ZB:
		_api = zb.New(builder.client, builder.apiKey, builder.secretkey)
	case COINEX:
		_api = coinex.New(builder.client, builder.apiKey, builder.secretkey)
	case BIGONE:
		_api = bigone.New(builder.client, builder.apiKey, builder.secretkey)
	case HITBTC:
		_api = hitbtc.New(builder.client, builder.apiKey, builder.secretkey)
	case ATOP:
		_api = atop.New(builder.client, builder.apiKey, builder.secretkey)
	default:
		println("exchange name error [" + exName + "].")

	}
	return _api
}

func (builder *APIBuilder) BuildFuture(exName string) (api FutureRestAPI) {
	switch exName {
	case BITMEX:
		return bitmex.New(&APIConfig{
			//Endpoint:     "https://www.bitmex.com/",
			Endpoint:     builder.futuresEndPoint,
			HttpClient:   builder.client,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey})
	case BITMEX_TEST:
		return bitmex.New(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     "https://testnet.bitmex.com",
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		})
	case OKEX_FUTURE, OKEX_V3:
		//return okcoin.NewOKEx(builder.client, builder.apiKey, builder.secretkey)
		return okex.NewOKEx(&APIConfig{
			HttpClient: builder.client,
			//	Endpoint:      "https://www.okex.com",
			Endpoint:      builder.futuresEndPoint,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretkey,
			ApiPassphrase: builder.apiPassphrase}).OKExFuture
	case HBDM:
		return huobi.NewHbdm(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey})
	case HBDM_SWAP:
		return huobi.NewHbdmSwap(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		})
	case OKEX_SWAP:
		return okex.NewOKEx(&APIConfig{
			HttpClient:    builder.client,
			Endpoint:      builder.futuresEndPoint,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretkey,
			ApiPassphrase: builder.apiPassphrase}).OKExSwap
	case COINBENE:
		return coinbene.NewCoinbeneSwap(APIConfig{
			HttpClient: builder.client,
			//	Endpoint:     "http://openapi-contract.coinbene.com",
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		})

	case BINANCE_SWAP:
		return binance.NewBinanceSwap(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		})
	case BINANCE, BINANCE_FUTURES:
		return binance.NewBinanceFutures(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		})
	default:
		println(fmt.Sprintf("%s not support future", exName))
		return nil
	}
}

func (builder *APIBuilder) BuildFuturesWs(exName string) (FuturesWsApi, error) {
	switch exName {
	case OKEX_V3, OKEX, OKEX_FUTURE:
		return okex.NewOKExV3FuturesWs(okex.NewOKEx(&APIConfig{
			HttpClient: builder.client,
			Endpoint:   builder.futuresEndPoint,
		})), nil
	case HBDM:
		return huobi.NewHbdmWs(), nil
	case HBDM_SWAP:
		return huobi.NewHbdmSwapWs(), nil
	case BINANCE, BINANCE_FUTURES, BINANCE_SWAP:
		return binance.NewFuturesWs(), nil
	case BITMEX:
		return bitmex.NewSwapWs(), nil
	}
	return nil, errors.New("not support the exchange " + exName)
}

func (builder *APIBuilder) BuildSpotWs(exName string) (SpotWsApi, error) {
	switch exName {
	case OKEX_V3, OKEX:
		return okex.NewOKExSpotV3Ws(nil), nil
	case HUOBI_PRO, HUOBI:
		return huobi.NewSpotWs(), nil
	case BINANCE:
		return binance.NewSpotWs(), nil
	}
	return nil, errors.New("not support the exchange " + exName)
}

func (builder *APIBuilder) BuildWallet(exName string) (WalletApi, error) {
	switch exName {
	case OKEX_V3, OKEX:
		return okex.NewOKEx(&APIConfig{
			HttpClient:    builder.client,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretkey,
			ApiPassphrase: builder.apiPassphrase,
		}).OKExWallet, nil
	case HUOBI_PRO:
		return huobi.NewWallet(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		}), nil
	case BINANCE:
		return binance.NewWallet(&APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretkey,
		}), nil
	}
	return nil, errors.New("not support the wallet api for  " + exName)
}
