package goex

import (
	"encoding/json"
	. "github.com/nntaoli-project/goex/internal/logger"
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
	Log.SetLevel(DEBUG)

	clientId := "a"
	args := make([]interface{}, 0)

	var heartbeatFunc = func() []byte {
		ts := time.Now().Unix()*1000 + 42029
		args = append(args, ts)
		//ping := fmt.Sprintf("{\"cmd\":\"ping\",\"args\":[%d],\"id\":\"%s\"}", ts, clientId)
		ping2 := map[string]interface{}{
			"cmd":  "ping",
			"id":   clientId,
			"args": args}

		ping3, _ := json.Marshal(ping2)
		return ping3
	}

	//fmt.Println(ping)
	//fmt.Println(ping2)
	//fmt.Println(err, string(ping3))

	ws := NewWsBuilder().Dump().WsUrl("wss://api.fcoin.com/v2/ws").
		ProxyUrl("socks5://127.0.0.1:1080").AutoReconnect().
		Heartbeat(heartbeatFunc, 5*time.Second).ProtoHandleFunc(ProtoHandle).Build()
	t.Log(ws.Subscribe(map[string]string{
		//"cmd":"sub", "args":"[\"ticker.btcusdt\"]", "id": clientId}))
		"cmd":"sub", "args":"ticker.btcusdt", "id": clientId}))
	time.Sleep(time.Second * 20)
	ws.c.Close()
	time.Sleep(time.Second*120)
}
