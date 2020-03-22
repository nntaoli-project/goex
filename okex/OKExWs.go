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
	once          *sync.Once
	wsConn        *WsConn
	respHandle    func(channel string, data json.RawMessage) error
	loginCh       chan wsResp
	isLogin       bool
	loginLock     *sync.Mutex
	authoriedSubs []map[string]interface{}
}

func NewOKExV3Ws(base *OKEx, handle func(channel string, data json.RawMessage) error) *OKExV3Ws {
	okV3Ws := &OKExV3Ws{
		once:       new(sync.Once),
		base:       base,
		respHandle: handle,
	}
	okV3Ws.loginCh = make(chan wsResp)
	okV3Ws.loginLock = &sync.Mutex{}
	okV3Ws.authoriedSubs = make([]map[string]interface{}, 0)
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:8443/ws/v3").
		ReconnectInterval(2*time.Second).
		AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).
		UnCompressFunc(FlateUnCompress).ProtoHandleFunc(okV3Ws.handle)
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

func (okV3Ws *OKExV3Ws) authoriedSubscribe(data map[string]interface{}) error {
	okV3Ws.authoriedSubs = append(okV3Ws.authoriedSubs, data)
	return okV3Ws.Subscribe(data)
}

func (okV3Ws *OKExV3Ws) reSubscribeAuthoriedChannel() {
	for _, d := range okV3Ws.authoriedSubs {
		okV3Ws.wsConn.SendJsonMessage(d)
	}
}

func (okV3Ws *OKExV3Ws) connectWs() {
	okV3Ws.once.Do(func() {
		okV3Ws.wsConn = okV3Ws.WsBuilder.Build()
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
		case "login":
			select {
			case okV3Ws.loginCh <- wsResp:
				return nil
			default:
				return nil
			}
		case "error":
			var errorCode int
			switch v := wsResp.ErrorCode.(type) {
			case int:
				errorCode = v
			case float64:
				errorCode = int(v) // float64 okex牛逼嗷
			case string:
				i, _ := strconv.ParseInt(v, 10, 64)
				errorCode = int(i)
			}

			switch errorCode {
			// event:error message:Already logged in errorCode:30042
			case 30041:
				//TODO:
				return nil
			case 30042:
				return nil
			}

			// TODO: clearfy which errors should be treated as login result.
			select {
			case okV3Ws.loginCh <- wsResp:
				return nil
			default:
				return fmt.Errorf("error in websocket: %v", wsResp)
			}
		}
		return fmt.Errorf("unknown websocket message: %v", wsResp)
	}

	if wsResp.Table != "" {
		channel, err := okV3Ws.parseChannel(wsResp.Table)
		if err != nil {
			logger.Error("parse ws channel error:", err)
			return err
		}
		err = okV3Ws.respHandle(channel, wsResp.Data)
		if err != nil {
			logger.Error("handle ws data error:", err)
		}
		return err
	}

	return fmt.Errorf("unknown websocket message: %v", wsResp)
}

func (okV3Ws *OKExV3Ws) Login() error {
	// already logined
	if okV3Ws.isLogin {
		return nil
	}

	okV3Ws.connectWs()

	okV3Ws.loginLock.Lock()
	defer okV3Ws.loginLock.Unlock()

	if okV3Ws.isLogin { //double check
		return nil
	}

	okV3Ws.clearChan(okV3Ws.loginCh)

	sign, tm := okV3Ws.base.doParamSign("GET", "/users/self/verify", "")
	op := map[string]interface{}{
		"op": "login", "args": []string{okV3Ws.base.config.ApiKey, okV3Ws.base.config.ApiPassphrase, tm, sign}}
	err := okV3Ws.wsConn.SendJsonMessage(op)
	if err != nil {
		logger.Error("ws login error:", err)
		return err
	}

	//wait login response
	re := <-okV3Ws.loginCh
	if !re.Success {
		return fmt.Errorf("login failed: %v", re)
	}
	logger.Info("ws login success")
	okV3Ws.isLogin = true
	return nil
}

func (okV3Ws *OKExV3Ws) Subscribe(sub map[string]interface{}) error {
	okV3Ws.connectWs()
	return okV3Ws.wsConn.Subscribe(sub)
}
