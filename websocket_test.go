package goex

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"
)

func Test_time(t *testing.T) {
	t.Log(time.Now().Unix())
}

func TestNewWsConn(t *testing.T) {
	//os.Setenv("https_proxy" , "socks5://127.0.0.1:1080")
	ws := NewWsConn("wss://api.huobipro.com/ws")
	//ws := NewWsConn("wss://real.okex.com:10441/websocket")
	time.Sleep(time.Second)

	ws.Heartbeat(func() interface{} {
		return map[string]interface{}{"ping": time.Duration(time.Now().Nanosecond())}
	}, 5*time.Second)

	ws.ReConnect()

	ws.Subscribe(map[string]string{
		"sub": "market.btcusdt.detail",
		"id":  "2"})
	ws.Subscribe(map[string]string{
		"sub": "market.btcusdt.depth.step0",
		"id":  "1"})

	//ws.WriteJSON(map[string]string{"event": "addChannel", "channel": "ok_sub_spot_bch_btc_ticker"})
	ws.ReceiveMessage(func(msg []byte) {
		println("receive message...")
		gzipreader, _ := gzip.NewReader(bytes.NewReader(msg))
		data, _ := ioutil.ReadAll(gzipreader)
		var resp map[string]interface{}
		json.Unmarshal(data, &resp)
		if resp["ping"] != nil {
			ws.WriteJSON(map[string]interface{}{"pong": resp["ping"]})
			ws.actived = time.Now()
		}
		println(string(data))
	})

	time.Sleep(2 * time.Minute)

	ws.CloseWs()

	time.Sleep(time.Second)
}
