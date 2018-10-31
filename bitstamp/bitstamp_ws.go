package bitstamp

import (
	"encoding/json"
	"fmt"
	"github.com/nntaoli-project/GoEx"
	"log"
	"sort"
	"strings"
	"time"
)

type Event struct {
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
}

func (bm *Bitstamp) createWsConn() {
	if bm.ws != nil {
		return
	}

	//connect wsx
	bm.createWsLock.Lock()
	defer bm.createWsLock.Unlock()

	if bm.ws == nil {
		bm.wsDepthHandleMap = make(map[string]func(*goex.Depth), 1)
		bm.wsTickerHandleMap = make(map[string]func(*goex.Ticker), 1)

		bm.ws = goex.NewWsConn("wss://ws.pusherapp.com/app/de504dc5763aeef9ff52?protocol=7&client=js&version=2.1.6&flash=false")
		bm.ws.Heartbeat(func() interface{} { return Event{Event: "pusher:ping"} }, 10*time.Second)
		bm.ws.ReConnect()
		bm.ws.ReceiveMessage(func(msg []byte) {
			var e Event
			err := json.Unmarshal(msg, &e)
			if err != nil {
				log.Println(err)
				return
			}
			switch e.Event {
			case "pusher:pong":
				bm.ws.UpdateActivedTime()
			case "data":
				pair := bm.getPairFromChannel(e.Channel)
				dep := bm.parseDepth(e.Data.(string))
				dep.Pair = pair
				if strings.HasPrefix(e.Channel, "order_book") {
					bm.wsDepthHandleMap[e.Channel](dep)
				}
			default:
				log.Printf("%+v", e)
			}
		})
	}
}

func (bm *Bitstamp) GetDepthWithWs(pair goex.CurrencyPair, handle func(*goex.Depth)) error {
	bm.createWsConn()
	channel := fmt.Sprintf("order_book_%s", strings.ToLower(pair.ToSymbol("")))
	if pair == goex.BTC_USD {
		channel = "order_book"
	}
	bm.wsDepthHandleMap[channel] = handle
	e := &Event{
		Event: "pusher:subscribe",
		Data: map[string]interface{}{
			"channel": channel}}
	return bm.ws.Subscribe(e)
}

func (bm *Bitstamp) parseDepth(dep string) *goex.Depth {
	var depthmap map[string]interface{}
	err := json.Unmarshal([]byte(dep), &depthmap)
	if err != nil {
		log.Println(err)
		return &goex.Depth{}
	}
	var depth goex.Depth

	bids, isok1 := depthmap["bids"].([]interface{})
	asks, isok2 := depthmap["asks"].([]interface{})
	if !isok1 || !isok2 {
		return &depth
	}

	for _, v := range bids {
		bid := v.([]interface{})
		depth.BidList = append(depth.BidList, goex.DepthRecord{goex.ToFloat64(bid[0]), goex.ToFloat64(bid[1])})
	}

	for _, v := range asks {
		ask := v.([]interface{})
		depth.AskList = append(depth.AskList, goex.DepthRecord{goex.ToFloat64(ask[0]), goex.ToFloat64(ask[1])})
	}

	sort.Sort(sort.Reverse(depth.AskList)) //reverse

	return &depth
}

func (bm *Bitstamp) getPairFromChannel(channel string) goex.CurrencyPair {
	if channel == "order_book" {
		return goex.BTC_USD
	}
	metas := strings.Split(channel, "_")
	pairstr := metas[2]
	return goex.NewCurrencyPair2(pairstr[0:3] + "_" + pairstr[3:])
}
