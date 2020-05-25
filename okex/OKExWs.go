package okex

import (
	"encoding/json"
	"fmt"
	"github.com/nntaoli-project/goex/internal/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/nntaoli-project/goex"
)

type wsResp struct {
	Event     string `json:"event"`
	Channel   string `json:"channel"`
	Table     string `json:"table"`
	Data      json.RawMessage
	Success   bool        `json:"success"`
	ErrorCode interface{} `json:"errorCode"`
}

type OKExV3Ws struct {
	base *OKEx
	*WsBuilder
	once       *sync.Once
	WsConn     *WsConn
	respHandle func(channel string, data json.RawMessage) error
}

func NewOKExV3Ws(base *OKEx, handle func(channel string, data json.RawMessage) error) *OKExV3Ws {
	okV3Ws := &OKExV3Ws{
		once:       new(sync.Once),
		base:       base,
		respHandle: handle,
	}
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:8443/ws/v3").
		ReconnectInterval(time.Second).
		AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).
		DecompressFunc(FlateDecompress).ProtoHandleFunc(okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3Ws) clearChan(c chan wsResp) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}

func (okV3Ws *OKExV3Ws) getTablePrefix(currencyPair CurrencyPair, contractType string) string {
	if contractType == SWAP_CONTRACT {
		return "swap"
	}
	return "futures"
}

func (okV3Ws *OKExV3Ws) ConnectWs() {
	okV3Ws.once.Do(func() {
		okV3Ws.WsConn = okV3Ws.WsBuilder.Build()
	})
}

func (okV3Ws *OKExV3Ws) parseChannel(channel string) (string, error) {
	metas := strings.Split(channel, "/")
	if len(metas) != 2 {
		return "", fmt.Errorf("unknown channel: %s", channel)
	}
	return metas[1], nil
}

func (okV3Ws *OKExV3Ws) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}

func (okV3Ws *OKExV3Ws) handle(msg []byte) error {
	logger.Debug("[ws] [response] ", string(msg))
	if string(msg) == "pong" {
		return nil
	}

	var wsResp wsResp
	err := json.Unmarshal(msg, &wsResp)
	if err != nil {
		logger.Error(err)
		return err
	}

	if wsResp.ErrorCode != nil {
		logger.Error(string(msg))
		return fmt.Errorf("%s", string(msg))
	}

	if wsResp.Event != "" {
		switch wsResp.Event {
		case "subscribe":
			logger.Info("subscribed:", wsResp.Channel)
			return nil
		case "error":
			logger.Errorf(string(msg))
		default:
			logger.Info(string(msg))
		}
		return fmt.Errorf("unknown websocket message: %v", wsResp)
	}

	if wsResp.Table != "" {
		err = okV3Ws.respHandle(wsResp.Table, wsResp.Data)
		if err != nil {
			logger.Error("handle ws data error:", err)
		}
		return err
	}

	return fmt.Errorf("unknown websocket message: %v", wsResp)
}

func (okV3Ws *OKExV3Ws) Subscribe(sub map[string]interface{}) error {
	okV3Ws.ConnectWs()
	return okV3Ws.WsConn.Subscribe(sub)
}
