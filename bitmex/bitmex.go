package bitmex

import (
	"fmt"
	"net/http"
	"time"

	apiclient "github.com/nntaoli-project/GoEx/bitmex/client"
	//"github.com/nntaoli-project/GoEx/bitmex/client/instrument"
	"github.com/nntaoli-project/GoEx/bitmex/client/order_book"
	//"github.com/nntaoli-project/GoEx/bitmex/client/position"
	//"github.com/nntaoli-project/GoEx/bitmex/client/trade"
	//apiuser "github.com/nntaoli-project/GoEx/bitmex/client/user"
	//"github.com/nntaoli-project/GoEx/bitmex/models"
	//"github.com/go-openapi/strfmt"
	"github.com/json-iterator/go"

	. "github.com/nntaoli-project/GoEx"

	"io/ioutil"
	"os"
)

const (
	BaseURL     = "www.bitmex.com"
	BasePath    = "/api/v1"
	TestBaseURL = "testnet.bitmex.com"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

//bitmex register link  https://www.bitmex.com/register/0fcQP7

type Bitmex struct {
	accessKey,
	secretKey string
	lever float64
	trans *Transport
	api   *apiclient.APIClient
}

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

func New(client *http.Client, accesskey, secretkey, baseUrl, proxyUrl string) *Bitmex {
	b := new(Bitmex)
	b.setProxy(proxyUrl)
	cfg := &apiclient.TransportConfig{}
	cfg.Host = baseUrl
	cfg.BasePath = BasePath
	cfg.Schemes = []string{"https"}

	b.api = apiclient.NewHTTPClientWithConfig(nil, cfg)
	b.trans = NewTransport(cfg.Host, cfg.BasePath, accesskey, secretkey, cfg.Schemes)
	b.api.SetTransport(b.trans)
	b.setTimeOffset()

	return b
}

func (b *Bitmex) setProxy(proxyUrl string) {
	os.Setenv("HTTP_PROXY", proxyUrl)
	os.Setenv("HTTPS_PROXY", proxyUrl)
}

func (b *Bitmex) setTimeOffset() error {
	info, err := b.info()
	if err != nil {
		fmt.Println(err)
		return err
	}
	nonce := time.Now().UnixNano()
	b.trans.timeOffset = nonce/1000000 - info.Timestamp
	return nil
}

// Info get server information
func (b *Bitmex) info() (info Info, err error) {
	url := fmt.Sprintf("https://%v%v", b.trans.Host, b.trans.BasePath)
	var response *http.Response
	response, err = http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &info)
	return
}

/**
 * 期货行情
 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 */
func (b *Bitmex) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	panic("not implements")
}

/**
 * 期货深度
 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param size 获取深度档数
 * @return
 */
func (b *Bitmex) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	nDepth := int32(size)
	ret, err := b.api.OrderBook.OrderBookGetL2(&order_book.OrderBookGetL2Params{Depth: &nDepth, Symbol: b.pairToSymbol(currencyPair)})
	if err != nil {
		return nil, err
	}
	depth := new(Depth)
	for _, v := range ret.Payload {
		if *v.Side == "Sell" {
			depth.AskList = append(depth.AskList,
				DepthRecord{Price: float64(v.Price),
					Amount: float64(v.Size)})
		} else {
			depth.BidList = append(depth.BidList,
				DepthRecord{Price: float64(v.Price),
					Amount: float64(v.Size)})
		}
	}
	depth.UTime = time.Now()
	return depth, nil
}

/**
 *获取交易所名字
 */
func (b *Bitmex) GetExchangeName() string {
	return BITMEX
}

/**
 *获取交割预估价
 */
func (b *Bitmex) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (b *Bitmex) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *全仓账户
 */
func (b *Bitmex) GetFutureUserinfo() (*FutureAccount, error) {
	panic("not implements")
}

/**
 * 期货下单
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param price  价格
 * @param amount  委托数量
 * @param openType   1:开多   2:开空   3:平多   4:平空
 * @param matchPrice  是否为对手价 0:不是    1:是   ,当取值为1时,price无效
 */
func (b *Bitmex) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	panic("not implements")
}

/**
 * 取消订单
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param orderId   订单ID

 */
func (b *Bitmex) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("not implements")
}

/**
 * 用户持仓查询
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @return
 */
func (b *Bitmex) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("not implements")
}

/**
 *获取订单信息
 */
func (b *Bitmex) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not implements")
}

/**
 *获取单个订单信息
 */
func (b *Bitmex) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("not implements")
}

/**
 *获取未完成订单信息
 */
func (b *Bitmex) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not implements")
}

/**
 *获取交易费
 */
func (b *Bitmex) GetFee() (float64, error) {
	panic("not implements")
}

/**
 *获取交易所的美元人民币汇率
 */
func (b *Bitmex) GetExchangeRate() (float64, error) {
	panic("not implements")
}

/**
 *获取每张合约价值
 */
func (b *Bitmex) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (b *Bitmex) GetDeliveryTime() (int, int, int, int) {
	panic("not implements")
}

/**
 * 获取K线数据
 */
func (b *Bitmex) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("not implements")
}

/**
 * 获取Trade数据
 *非个人，整个交易所的交易记录
 */
func (b *Bitmex) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (b *Bitmex) pairToSymbol(pair CurrencyPair) string {
	if pair.CurrencyA.Symbol == BTC.Symbol {
		return NewCurrencyPair(XBT, USD).ToSymbol("")
	}
	return pair.AdaptUsdtToUsd().ToSymbol("")
}
