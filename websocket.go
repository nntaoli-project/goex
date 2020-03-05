package goex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type WsConfig struct {
	WsUrl                 string
	ProxyUrl              string
	ReqHeaders            map[string][]string //连接的时候加入的头部信息
	HeartbeatIntervalTime time.Duration       //
	HeartbeatData         func() []byte       //心跳数据2
	IsAutoReconnect       bool
	ProtoHandleFunc       func([]byte) error           //协议处理函数
	UnCompressFunc        func([]byte) ([]byte, error) //解压函数
	ErrorHandleFunc       func(err error)
	IsDump                bool
	readDeadLineTime      time.Duration
	reconnectInterval     time.Duration
}

type WsConn struct {
	c *websocket.Conn
	sync.Mutex
	WsConfig
	writeBufferChan        chan []byte
	pingMessageBufferChan  chan []byte
	pongMessageBufferChan  chan []byte
	closeMessageBufferChan chan []byte
	subs                   []interface{}
	close                  chan bool
}

type WsBuilder struct {
	wsConfig *WsConfig
}

func NewWsBuilder() *WsBuilder {
	return &WsBuilder{&WsConfig{
		ReqHeaders:        make(map[string][]string, 1),
		reconnectInterval: time.Second * 10,
	}}
}

func (b *WsBuilder) WsUrl(wsUrl string) *WsBuilder {
	b.wsConfig.WsUrl = wsUrl
	return b
}

func (b *WsBuilder) ProxyUrl(proxyUrl string) *WsBuilder {
	b.wsConfig.ProxyUrl = proxyUrl
	return b
}

func (b *WsBuilder) ReqHeader(key, value string) *WsBuilder {
	b.wsConfig.ReqHeaders[key] = append(b.wsConfig.ReqHeaders[key], value)
	return b
}

func (b *WsBuilder) AutoReconnect() *WsBuilder {
	b.wsConfig.IsAutoReconnect = true
	return b
}

func (b *WsBuilder) Dump() *WsBuilder {
	b.wsConfig.IsDump = true
	return b
}

func (b *WsBuilder) Heartbeat(heartbeat func() []byte, t time.Duration) *WsBuilder {
	b.wsConfig.HeartbeatIntervalTime = t
	b.wsConfig.HeartbeatData = heartbeat
	return b
}

func (b *WsBuilder) ReconnectInterval(t time.Duration) *WsBuilder {
	b.wsConfig.reconnectInterval = t
	return b
}

func (b *WsBuilder) ProtoHandleFunc(f func([]byte) error) *WsBuilder {
	b.wsConfig.ProtoHandleFunc = f
	return b
}

func (b *WsBuilder) UnCompressFunc(f func([]byte) ([]byte, error)) *WsBuilder {
	b.wsConfig.UnCompressFunc = f
	return b
}

func (b *WsBuilder) ErrorHandleFunc(f func(err error)) *WsBuilder {
	b.wsConfig.ErrorHandleFunc = f
	return b
}

func (b *WsBuilder) Build() *WsConn {
	wsConn := &WsConn{WsConfig: *b.wsConfig}
	return wsConn.NewWs()
}

func (ws *WsConn) NewWs() *WsConn {
	ws.Lock()
	defer ws.Unlock()

	if ws.HeartbeatIntervalTime == 0 {
		ws.readDeadLineTime = time.Minute
	} else {
		ws.readDeadLineTime = ws.HeartbeatIntervalTime * 2
	}

	if err := ws.connect(); err != nil {
		Log.Panic(fmt.Errorf("[%s] %s", ws.WsUrl, err.Error()))
	}

	ws.close = make(chan bool, 1)
	ws.pingMessageBufferChan = make(chan []byte, 10)
	ws.pongMessageBufferChan = make(chan []byte, 10)
	ws.closeMessageBufferChan = make(chan []byte, 10)
	ws.writeBufferChan = make(chan []byte, 10)

	go ws.writeRequest()
	go ws.receiveMessage()

	return ws
}

func (ws *WsConn) connect() error {
	dialer := websocket.DefaultDialer

	if ws.ProxyUrl != "" {
		proxy, err := url.Parse(ws.ProxyUrl)
		if err == nil {
			Log.Infof("[ws][%s] proxy url:%s", ws.WsUrl, proxy)
			dialer.Proxy = http.ProxyURL(proxy)
		} else {
			Log.Errorf("[ws][%s]parse proxy url [%s] err %s  ", ws.WsUrl, ws.ProxyUrl, err.Error())
		}
	}

	wsConn, resp, err := dialer.Dial(ws.WsUrl, http.Header(ws.ReqHeaders))
	if err != nil {
		Log.Errorf("[ws][%s] %s", ws.WsUrl, err.Error())
		if ws.IsDump && resp != nil {
			dumpData, _ := httputil.DumpResponse(resp, true)
			Log.Debugf("[ws][%s] %s", ws.WsUrl, string(dumpData))
		}
		return err
	}

	ws.c = wsConn

	if ws.HeartbeatIntervalTime > 0 {
		wsConn.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
	}

	if ws.IsDump {
		dumpData, _ := httputil.DumpResponse(resp, true)
		Log.Debugf("[ws][%s] %s", ws.WsUrl, string(dumpData))
	}
	Log.Infof("[ws][%s] connected", ws.WsUrl)

	return nil
}

func (ws *WsConn) reconnect() {
	ws.c.Close() //主动关闭一次
	var err error
	for retry := 1; retry <= 100; retry++ {
		err = ws.connect()
		if err != nil {
			Log.Errorf("[ws] [%s] websocket reconnect fail , %s", ws.WsUrl, err.Error())
		} else {
			break
		}
		time.Sleep(ws.WsConfig.reconnectInterval * time.Duration(retry))
	}

	if err != nil {
		Log.Errorf("[ws] [%s] retry reconnect fail , begin exiting. ", ws.WsUrl)
		ws.CloseWs()
		if ws.ErrorHandleFunc != nil {
			ws.ErrorHandleFunc(errors.New("retry reconnect fail"))
		}
	} else {
		//re subscribe
		var tmp []interface{}
		copy(tmp, ws.subs)
		ws.subs = ws.subs[:0]
		for _, sub := range tmp {
			ws.Subscribe(sub)
		}
	}
}

func (ws *WsConn) writeRequest() {
	var (
		heartTimer *time.Timer
		err        error
	)

	if ws.HeartbeatIntervalTime == 0 {
		heartTimer = time.NewTimer(time.Hour)
	} else {
		heartTimer = time.NewTimer(ws.HeartbeatIntervalTime)
	}

	for {
		select {
		case <-ws.close:
			Log.Infof("[ws][%s] close websocket , exiting write message goroutine.", ws.WsUrl)
			return
		case d := <-ws.writeBufferChan:
			err = ws.c.WriteMessage(websocket.TextMessage, d)
		case d := <-ws.pingMessageBufferChan:
			err = ws.c.WriteMessage(websocket.PingMessage, d)
		case d := <-ws.pongMessageBufferChan:
			err = ws.c.WriteMessage(websocket.PongMessage, d)
		case d := <-ws.closeMessageBufferChan:
			err = ws.c.WriteMessage(websocket.CloseMessage, d)
		case <-heartTimer.C:
			if ws.HeartbeatIntervalTime > 0 {
				//Log.Debug("send heartbeat data")
				err = ws.c.WriteMessage(websocket.TextMessage, ws.HeartbeatData())
				heartTimer.Reset(ws.HeartbeatIntervalTime)
			}
		}

		if err != nil {
			Log.Errorf("[ws][%s] %s", ws.WsUrl, err.Error())
			time.Sleep(time.Second)
		}
	}
}

func (ws *WsConn) Subscribe(subEvent interface{}) error {
	data, err := json.Marshal(subEvent)
	if err != nil {
		Log.Errorf("[ws][%s] json encode error , %s", ws.WsUrl, err)
		return err
	}
	ws.writeBufferChan <- data
	ws.subs = append(ws.subs, subEvent)
	return nil
}

func (ws *WsConn) SendMessage(msg []byte) {
	ws.writeBufferChan <- msg
}

func (ws *WsConn) SendPingMessage(msg []byte) {
	ws.pingMessageBufferChan <- msg
}

func (ws *WsConn) SendPongMessage(msg []byte) {
	ws.pongMessageBufferChan <- msg
}

func (ws *WsConn) SendCloseMessage(msg []byte) {
	ws.closeMessageBufferChan <- msg
}

func (ws *WsConn) SendJsonMessage(m interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	ws.writeBufferChan <- data
	return nil
}

func (ws *WsConn) receiveMessage() {
	//exit
	ws.c.SetCloseHandler(func(code int, text string) error {
		Log.Infof("[ws][%s] websocket exiting [code=%d , text=%s]", ws.WsUrl, code, text)
		ws.CloseWs()
		return nil
	})

	ws.c.SetPongHandler(func(pong string) error {
		Log.Debugf("[%s] received [pong] %s", ws.WsUrl, pong)
		ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})

	ws.c.SetPingHandler(func(ping string) error {
		Log.Debugf("[%s] received [ping] %s", ws.WsUrl, ping)
		ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})

	for {
		select {
		case <-ws.close:
			Log.Infof("[ws][%s] close websocket , exiting receive message goroutine.", ws.WsUrl)
			return
		default:
			t, msg, err := ws.c.ReadMessage()

			if err != nil {
				Log.Errorf("[ws][%s] %s", ws.WsUrl, err.Error())
				if ws.IsAutoReconnect {
					//	if _, ok := err.(*websocket.CloseError); ok {
					Log.Infof("[ws][%s] Unexpected Closed , Begin Retry Connect.", ws.WsUrl)
					ws.reconnect()
					//	}
					continue
				}

				if ws.ErrorHandleFunc != nil {
					ws.ErrorHandleFunc(err)
				}

				return
			}

			ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))

			switch t {
			case websocket.TextMessage:
				ws.ProtoHandleFunc(msg)
			case websocket.BinaryMessage:
				if ws.UnCompressFunc == nil {
					ws.ProtoHandleFunc(msg)
				} else {
					msg2, err := ws.UnCompressFunc(msg)
					if err != nil {
						Log.Errorf("[ws][%s] uncompress error %s", ws.WsUrl, err.Error())
					} else {
						ws.ProtoHandleFunc(msg2)
					}
				}
			case websocket.CloseMessage:
				ws.CloseWs()
				return
			default:
				Log.Errorf("[ws][%s] error websocket message type , content is :\n %s \n", ws.WsUrl, string(msg))
			}
		}
	}
}

func (ws *WsConn) CloseWs() {
	ws.close <- true
	close(ws.close)
	err := ws.c.Close()
	if err != nil {
		Log.Error("[ws][", ws.WsUrl, "]close websocket error ,", err)
	}
}

func (ws *WsConn) clearChannel(c chan struct{}) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}
