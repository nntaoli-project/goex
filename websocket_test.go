package goex

import (
	"fmt"
	"testing"
	"time"
)

func Test_time(t *testing.T) {
	t.Log(time.Now().Unix())
}

func ProtoHandle(data []byte) error {
	println(string(data))
	return nil
}

func TestNewWsConn(t *testing.T) {

	ws := NewWsBuilder().Dump().WsUrl("wss://real.okex.com:10440/ws/v1").ProxyUrl("socks5://127.0.0.1:1080").UnCompressFunc(FlateUnCompress).
		Heartbeat([]byte(fmt.Sprintf("{\"event\": \"%s\"}", "ping")), 5*time.Second).ProtoHandleFunc(ProtoHandle).Build()
	t.Log(ws.Subscribe(map[string]string{
		"event":   "addChannel", "channel": "ok_sub_spot_btc_usdt_depth_5"}))
	ws.ReceiveMessage()
	time.Sleep(time.Minute)
	ws.CloseWs()
}
