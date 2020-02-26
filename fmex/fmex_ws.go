package fmex

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	. "github.com/nntaoli-project/goex"
)

const (
	FMexWSTicker        = "ticker.%s"
	FMexWSOrderBook     = "depth.L%d.%s"
	FMexWSOrderBookL20  = "depth.L20.%s"
	FMexWSOrderBookL150 = "depth.L150.%s"
	FMexWSOrderBookFull = "depth.full.%s"
	FMexWSTrades        = "trade.%s"
	FMexWSKLines        = "candle.%s.%s"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FMexWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	klineCallback  func(*Kline, int)

	clientId      string
	subcribeTypes []string
	timeoffset    int64
}

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "M1",
	KLINE_PERIOD_3MIN:   "M3",
	KLINE_PERIOD_5MIN:   "M5",
	KLINE_PERIOD_15MIN:  "M15",
	KLINE_PERIOD_30MIN:  "M30",
	KLINE_PERIOD_60MIN:  "H1",
	KLINE_PERIOD_4H:     "H4",
	KLINE_PERIOD_6H:     "H6",
	KLINE_PERIOD_1DAY:   "D1",
	KLINE_PERIOD_1WEEK:  "W1",
	KLINE_PERIOD_1MONTH: "MN",
}
var _INERNAL_KLINE_PERIOD_REVERTER = map[string]int{
	"M1":  KLINE_PERIOD_1MIN,
	"M3":  KLINE_PERIOD_3MIN,
	"M5":  KLINE_PERIOD_5MIN,
	"M15": KLINE_PERIOD_15MIN,
	"M30": KLINE_PERIOD_30MIN,
	"H1":  KLINE_PERIOD_60MIN,
	"H4":  KLINE_PERIOD_4H,
	"H6":  KLINE_PERIOD_6H,
	"D1":  KLINE_PERIOD_1DAY,
	"W1":  KLINE_PERIOD_1WEEK,
	"MN":  KLINE_PERIOD_1MONTH,
}

func NewFMexWs(client *http.Client) *FMexWs {
	fmWs := &FMexWs{}
	fmWs.clientId = getRandomString(8)
	fmWs.WsBuilder = NewWsBuilder().
		WsUrl("wss://api.fmex.com/v2/ws").
		AutoReconnect().
		Heartbeat(func() []byte {
			ts := time.Now().Unix()*1000 + fmWs.timeoffset*1000
			args := make([]interface{}, 0)
			args = append(args, ts)
			heartbeatData := map[string]interface{}{
				"cmd":  "ping",
				"id":   fmWs.clientId,
				"args": args}
			data, _ := json.Marshal(heartbeatData)
			return data
		}, 25*time.Second).
		ProtoHandleFunc(fmWs.handle)

	return fmWs
}

//生成随机字符串
func getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (fmWs *FMexWs) SetCallbacks(
	tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, int),
) {
	fmWs.tickerCallback = tickerCallback
	fmWs.depthCallback = depthCallback
	fmWs.tradeCallback = tradeCallback
	fmWs.klineCallback = klineCallback
}

func (fmWs *FMexWs) subscribe(sub map[string]interface{}) error {
	fmWs.connectWs()
	return fmWs.wsConn.Subscribe(sub)
}

func (fmWs *FMexWs) SubscribeDepth(pair CurrencyPair, size int) error {
	if fmWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	arg := fmt.Sprintf(FMexWSOrderBook, size, adaptContractType(pair))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fmWs.subscribe(map[string]interface{}{
		"id":   fmWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fmWs *FMexWs) SubscribeTicker(pair CurrencyPair) error {
	if fmWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	arg := fmt.Sprintf(FMexWSTicker, adaptContractType(pair))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fmWs.subscribe(map[string]interface{}{
		"id":   fmWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fmWs *FMexWs) SubscribeTrade(pair CurrencyPair) error {
	if fmWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	arg := fmt.Sprintf(FMexWSTrades, adaptContractType(pair))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fmWs.subscribe(map[string]interface{}{
		"id":   fmWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fmWs *FMexWs) SubscribeKline(pair CurrencyPair, period int) error {
	if fmWs.klineCallback == nil {
		return errors.New("place set kline callback func")
	}
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "M1"
	}

	arg := fmt.Sprintf(FMexWSKLines, periodS, strings.ToLower(pair.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fmWs.subscribe(map[string]interface{}{
		"id":   fmWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fmWs *FMexWs) connectWs() {
	fmWs.Do(func() {
		fmWs.wsConn = fmWs.WsBuilder.Build()
	})
}

func (fmWs *FMexWs) parseTickerData(tickmap []interface{}) *RawTicker {
	t := new(RawTicker)
	t.Date = uint64(time.Now().UnixNano() / 1000000)
	t.Last = ToFloat64(tickmap[0])
	t.Vol = ToFloat64(tickmap[9])
	t.Low = ToFloat64(tickmap[8])
	t.High = ToFloat64(tickmap[7])
	t.Buy = ToFloat64(tickmap[2])
	t.Sell = ToFloat64(tickmap[4])
	t.SellAmount = ToFloat64(tickmap[5])
	t.BuyAmount = ToFloat64(tickmap[3])
	t.LastTradeVol = ToFloat64(tickmap[1])

	return t
}

func (fmWs *FMexWs) parseDepthData(bids, asks []interface{}) *Depth {
	depth := new(Depth)
	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, DepthRecord{ToFloat64(bids[i]), ToFloat64(bids[i+1])})
		i += 2
		n++
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, DepthRecord{ToFloat64(asks[i]), ToFloat64(asks[i+1])})
		i += 2
		n++
	}

	return depth
}

func (fmWs *FMexWs) handle(msg []byte) error {
	//fmt.Println("ws msg:", string(msg))
	datamap := make(map[string]interface{})
	err := json.Unmarshal(msg, &datamap)
	if err != nil {
		fmt.Println("json unmarshal error for ", string(msg))
		return err
	}

	msgType, isOk := datamap["type"].(string)
	if isOk {
		resp := strings.Split(msgType, ".")
		switch resp[0] {
		case "hello", "ping":
			stime := int64(ToInt(datamap["ts"]))
			st := time.Unix(0, stime*1000*1000)
			lt := time.Now()
			offset := st.Sub(lt).Seconds()
			fmWs.timeoffset = int64(offset)
		case "ticker":
			tick := fmWs.parseTickerData(datamap["ticker"].([]interface{}))
			pair, err := getPairFromType(resp[1])
			if err != nil {
				panic(err)
			}
			tick.Pair = pair
			fmWs.tickerCallback((*Ticker)(unsafe.Pointer(tick)))
			return nil
		case "depth":
			dep := fmWs.parseDepthData(datamap["bids"].([]interface{}), datamap["asks"].([]interface{}))
			stime := int64(ToInt(datamap["ts"]))
			dep.UTime = time.Unix(0, stime*1000000)
			pair, err := getPairFromType(resp[2])
			if err != nil {
				panic(err)
			}
			dep.Pair = pair

			fmWs.depthCallback(dep)
			return nil
		case "candle":
			period := _INERNAL_KLINE_PERIOD_REVERTER[resp[1]]
			kline := &Kline{
				Timestamp: int64(ToInt(datamap["id"])),
				Open:      ToFloat64(datamap["open"]),
				Close:     ToFloat64(datamap["close"]),
				High:      ToFloat64(datamap["high"]),
				Low:       ToFloat64(datamap["low"]),
				Vol:       ToFloat64(datamap["quote_vol"]),
			}
			pair, err := getPairFromType(resp[2])
			if err != nil {
				panic(err)
			}
			kline.Pair = pair
			fmWs.klineCallback(kline, period)
			return nil
		case "trade":
			side := BUY
			if datamap["side"] == "sell" {
				side = SELL
			}
			trade := &Trade{
				Tid:    int64(ToUint64(datamap["id"])),
				Type:   TradeSide(side),
				Amount: ToFloat64(datamap["amount"]),
				Price:  ToFloat64(datamap["price"]),
				Date:   int64(ToUint64(datamap["ts"])),
			}
			pair, err := getPairFromType(resp[1])
			if err != nil {
				panic(err)
			}
			trade.Pair = pair
			fmWs.tradeCallback(trade)
			return nil
		default:
			return errors.New("unknown message " + msgType)

		}
	}
	return nil
}

func getPairFromType(pair string) (CurrencyPair, error) {
	s := strings.Split(pair, "usd_")
	if len(s) == 2 {
		return NewCurrencyPair2(s[0] + "_" + "USD"), nil
	}
	return CurrencyPair{}, errors.New("pair not support :" + pair)
}
