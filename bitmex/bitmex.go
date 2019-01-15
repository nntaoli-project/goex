package bitmex

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
	"time"
)

var (
	BaseURL     = "https://www.bitmex.com/api/v1/"
	TestBaseURL = "https://testnet.bitmex.com/api/v1/"
)

//bitmex register link  https://www.bitmex.com/register/0fcQP7

type Bitmex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(client *http.Client, accesskey, secretkey, baseUrl string) *Bitmex {

	return &Bitmex{httpClient: client, accessKey: accesskey, secretKey: secretkey}
}

/**
 * 期货行情
 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 */
func (bm *Bitmex) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	panic("not implements")
}

/**
 * 期货深度
 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param size 获取深度档数
 * @return
 */
func (bm *Bitmex) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	uri := fmt.Sprintf("orderBook/L2?symbol=%s&depth=%d", bm.pairToSymbol(currencyPair), size)
	resp, err := HttpGet3(bm.httpClient, base_url+uri, nil)
	if err != nil {
		return nil, HTTP_ERR_CODE.OriginErr(err.Error())
	}

	//log.Println(resp)

	dep := new(Depth)
	dep.UTime = time.Now()
	dep.Pair = currencyPair

	for _, r := range resp {
		rr := r.(map[string]interface{})
		switch strings.ToLower(rr["side"].(string)) {
		case "sell":
			dep.AskList = append(dep.AskList, DepthRecord{Price: ToFloat64(rr["price"]), Amount: ToFloat64(rr["size"])})
		case "buy":
			dep.BidList = append(dep.BidList, DepthRecord{Price: ToFloat64(rr["price"]), Amount: ToFloat64(rr["size"])})
		}
	}

	return dep, nil
}

/**
 *获取交易所名字
 */
func (bm *Bitmex) GetExchangeName() string {
	return BITMEX
}

/**
 *获取交割预估价
 */
func (bm *Bitmex) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (bm *Bitmex) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *全仓账户
 */
func (bm *Bitmex) GetFutureUserinfo() (*FutureAccount, error) {
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
func (bm *Bitmex) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	panic("not implements")
}

/**
 * 取消订单
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param orderId   订单ID

 */
func (bm *Bitmex) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("not implements")
}

/**
 * 用户持仓查询
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @return
 */
func (bm *Bitmex) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("not implements")
}

/**
 *获取订单信息
 */
func (bm *Bitmex) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not implements")
}

/**
 *获取单个订单信息
 */
func (bm *Bitmex) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("not implements")
}

/**
 *获取未完成订单信息
 */
func (bm *Bitmex) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not implements")
}

/**
 *获取交易费
 */
func (bm *Bitmex) GetFee() (float64, error) {
	panic("not implements")
}

/**
 *获取交易所的美元人民币汇率
 */
func (bm *Bitmex) GetExchangeRate() (float64, error) {
	panic("not implements")
}

/**
 *获取每张合约价值
 */
func (bm *Bitmex) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (bm *Bitmex) GetDeliveryTime() (int, int, int, int) {
	panic("not implements")
}

/**
 * 获取K线数据
 */
func (bm *Bitmex) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("not implements")
}

/**
 * 获取Trade数据
 *非个人，整个交易所的交易记录
 */
func (bm *Bitmex) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

/**
 */
func (bm *Bitmex) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (bm *Bitmex) signature(req runtime.ClientRequest, operation *runtime.ClientOperation, formats strfmt.Registry) (expires, sign string) {

}

func (bm *Bitmex) pairToSymbol(pair CurrencyPair) string {
	if pair.CurrencyA.Symbol == BTC.Symbol {
		return NewCurrencyPair(XBT, USD).ToSymbol("")
	}
	return pair.AdaptUsdtToUsd().ToSymbol("")
}
