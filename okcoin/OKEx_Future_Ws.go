package okcoin

import (
	"encoding/json"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"strings"
	"time"
)

func (okFuture *OKEx) createWsConn() {
	if okFuture.ws == nil {
		//connect wsx
		okFuture.createWsLock.Lock()
		defer okFuture.createWsLock.Unlock()

		if okFuture.ws == nil {
			okFuture.wsTickerHandleMap = make(map[string]func(*Ticker))
			okFuture.wsDepthHandleMap = make(map[string]func(*Depth))

			okFuture.ws = NewWsConn("wss://real.okex.com:10440/websocket/okexapi")
			okFuture.ws.Heartbeat(func() interface{} { return map[string]string{"event": "ping"} }, 30*time.Second)
			okFuture.ws.ReConnect()
			okFuture.ws.ReceiveMessage(func(msg []byte) {
				if string(msg) == "{\"event\":\"pong\"}" {
					okFuture.ws.UpdateActivedTime()
					return
				}

				var data []interface{}
				err := json.Unmarshal(msg, &data)
				if err != nil {
					log.Print(err)
					return
				}

				if len(data) == 0 {
					return
				}

				datamap := data[0].(map[string]interface{})
				channel := datamap["channel"].(string)
				if channel == "addChannel" {
					return
				}

				tickmap := datamap["data"].(map[string]interface{})
				pair := okFuture.getPairFromChannel(channel)
				contractType := okFuture.getContractFromChannel(channel)

				if strings.HasSuffix(channel, "_ticker") {
					ticker := okFuture.parseTicker(tickmap)
					ticker.Pair = pair
					ticker.ContractType = contractType
					okFuture.wsTickerHandleMap[channel](ticker)
				} else if strings.Contains(channel, "depth_") {
					dep := okFuture.parseDepth(tickmap)
					dep.Pair = pair
					dep.ContractType = contractType
					okFuture.wsDepthHandleMap[channel](dep)
				}
			})
		}
	}
}

func (okFuture *OKEx) GetDepthWithWs(pair CurrencyPair, contractType string, handle func(*Depth)) error {
	okFuture.createWsConn()
	channel := fmt.Sprintf("ok_sub_futureusd_%s_depth_%s_5", strings.ToLower(pair.CurrencyA.Symbol), contractType)
	okFuture.wsDepthHandleMap[channel] = handle
	return okFuture.ws.WriteJSON(map[string]string{
		"event":   "addChannel",
		"channel": channel})
}

func (okFuture *OKEx) GetTickerWithWs(pair CurrencyPair, contractType string, handle func(*Ticker)) error {
	okFuture.createWsConn()
	channel := fmt.Sprintf("ok_sub_futureusd_%s_ticker_%s", strings.ToLower(pair.CurrencyA.Symbol), contractType)
	okFuture.wsTickerHandleMap[channel] = handle
	return okFuture.ws.WriteJSON(map[string]string{
		"event":   "addChannel",
		"channel": channel})
}

func (okFuture *OKEx) parseTicker(tickmap map[string]interface{}) *Ticker {
	return &Ticker{
		Last: ToFloat64(tickmap["last"]),
		Low:  ToFloat64(tickmap["low"]),
		High: ToFloat64(tickmap["high"]),
		Vol:  ToFloat64(tickmap["vol"]),
		Sell: ToFloat64(tickmap["sell"]),
		Buy:  ToFloat64(tickmap["buy"]),
		Date: ToUint64(tickmap["timestamp"])}
}

func (okFuture *OKEx) parseDepth(tickmap map[string]interface{}) *Depth {
	asks := tickmap["asks"].([]interface{})
	bids := tickmap["bids"].([]interface{})

	var depth Depth
	for _, v := range asks {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = ToFloat64(vv)
			case 1:
				dr.Amount = ToFloat64(vv)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bids {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = ToFloat64(vv)
			case 1:
				dr.Amount = ToFloat64(vv)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}
	return &depth
}

func (okFuture *OKEx) getPairFromChannel(channel string) CurrencyPair {
	metas := strings.Split(channel, "_")
	return NewCurrencyPair2(metas[3] + "_usd")
}

func (okFuture *OKEx) getContractFromChannel(channel string) string {
	if strings.Contains(channel, THIS_WEEK_CONTRACT) {
		return THIS_WEEK_CONTRACT
	}

	if strings.Contains(channel, NEXT_WEEK_CONTRACT) {
		return NEXT_WEEK_CONTRACT
	}

	if strings.Contains(channel, QUARTER_CONTRACT) {
		return QUARTER_CONTRACT
	}
	return ""
}
