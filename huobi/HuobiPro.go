package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type HuobiPro struct {
	*HuoBi_V2
	ws                *WsConn
	createWsLock      sync.Mutex
	wsTickerHandleMap map[string]func(*Ticker)
	wsDepthHandleMap  map[string]func(*Depth)
}

func NewHuobiPro(client *http.Client, apikey, secretkey, accountId string) *HuobiPro {
	hbv2 := new(HuoBi_V2)
	hbv2.accountId = accountId
	hbv2.accessKey = apikey
	hbv2.secretKey = secretkey
	hbv2.httpClient = client
	hbv2.baseUrl = "https://api.huobipro.com"
	return &HuobiPro{HuoBi_V2: hbv2, wsTickerHandleMap: make(map[string]func(*Ticker)), wsDepthHandleMap: make(map[string]func(*Depth))}
}

func (hbpro *HuobiPro) createWsConn() {
	if hbpro.ws == nil {
		//connect wsx
		hbpro.createWsLock.Lock()
		defer hbpro.createWsLock.Unlock()

		if hbpro.ws == nil {
			hbpro.ws = NewWsConn("wss://api.huobipro.com/ws")
			hbpro.ws.Heartbeat(func() interface{} {
				return map[string]interface{}{
					"ping": time.Now().Unix()}
			}, 5*time.Second)
			hbpro.ws.ReConnect()
			hbpro.ws.ReceiveMessage(func(msg []byte) {
				gzipreader, _ := gzip.NewReader(bytes.NewReader(msg))
				data, _ := ioutil.ReadAll(gzipreader)
				datamap := make(map[string]interface{})
				err := json.Unmarshal(data, &datamap)
				if err != nil {
					log.Println("json unmarshal error for ", string(data))
					return
				}

				if datamap["ping"] != nil {
					hbpro.ws.WriteJSON(map[string]interface{}{
						"pong": datamap["ping"]}) // 回应心跳
					return
				}

				if datamap["pong"] != nil { //
					return
				}

				if datamap["id"] != nil { //忽略订阅成功的回执消息
					log.Println(string(data))
					return
				}

				ch, isok := datamap["ch"].(string)
				if !isok {
					log.Println("error:", string(data))
					return
				}

				tick := datamap["tick"].(map[string]interface{})
				if hbpro.wsTickerHandleMap[ch] != nil {
					return
				}

				if hbpro.wsDepthHandleMap[ch] != nil {
					(hbpro.wsDepthHandleMap[ch])(hbpro.parseDepthData(tick))
					return
				}

				log.Println(string(data))
			})
		}
	}
}

func (hbpro *HuobiPro) parseDepthData(tick map[string]interface{}) *Depth {
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)
	}

	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)
	}

	return depth
}

func (hbpro *HuobiPro) GetExchangeName() string {
	return HUOBI_PRO
}

func (hbpro *HuobiPro) GetTickerWithWs(pair CurrencyPair, handle func(ticker *Ticker)) error {
	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.detail", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsTickerHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  1,
		"sub": sub})
}

func (hbpro *HuobiPro) GetDepthWithWs(pair CurrencyPair, handle func(dep *Depth)) error {
	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.depth.step0", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsDepthHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  2,
		"sub": sub})
}
