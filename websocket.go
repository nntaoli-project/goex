package goex

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type WsConfig struct {
	WsUrl                 string
	ProxyUrl              string
	ReqHeaders            map[string][]string          //连接的时候加入的头部信息
	HeartbeatIntervalTime time.Duration                //
	HeartbeatData         []byte                       //心跳数据
	HeartbeatFunc         func() interface{}           //心跳数据2
	ReconnectIntervalTime time.Duration                //定时重连时间间隔
	ProtoHandleFunc       func([]byte) error           //协议处理函数
	UnCompressFunc        func([]byte) ([]byte, error) //解压函数
	ErrorHandleFunc       func(err error)
	IsDump                bool
}

type WsConn struct {
	*websocket.Conn
	sync.Mutex
	WsConfig

	activeTime  time.Time
	activeTimeL sync.Mutex

	mu             chan struct{} // lock write data
	closeHeartbeat chan struct{}
	closeReconnect chan struct{}
	closeRecv      chan struct{}
	closeCheck     chan struct{}
	subs           []interface{}
}

type WsBuilder struct {
	wsConfig *WsConfig
}

func NewWsBuilder() *WsBuilder {
	return &WsBuilder{&WsConfig{}}
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

func (b *WsBuilder) Dump() *WsBuilder {
	b.wsConfig.IsDump = true
	return b
}

func (b *WsBuilder) Heartbeat(data []byte, t time.Duration) *WsBuilder {
	b.wsConfig.HeartbeatIntervalTime = t
	b.wsConfig.HeartbeatData = data
	return b
}
func (b *WsBuilder) Heartbeat2(heartbeat func() interface{}, t time.Duration) *WsBuilder {
	b.wsConfig.HeartbeatIntervalTime = t
	b.wsConfig.HeartbeatFunc = heartbeat
	return b
}

func (b *WsBuilder) ReconnectIntervalTime(t time.Duration) *WsBuilder {
	b.wsConfig.ReconnectIntervalTime = t
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
	if b.wsConfig.ErrorHandleFunc == nil {
		b.wsConfig.ErrorHandleFunc = func(err error) {
			log.Println(err)
		}
	}
	wsConn := &WsConn{WsConfig: *b.wsConfig}
	return wsConn.NewWs()
}

func (ws *WsConn) NewWs() *WsConn {
	ws.Lock()
	defer ws.Unlock()

	ws.connect()

	ws.mu = make(chan struct{}, 1)
	ws.closeHeartbeat = make(chan struct{}, 1)
	ws.closeReconnect = make(chan struct{}, 1)
	ws.closeRecv = make(chan struct{}, 1)
	ws.closeCheck = make(chan struct{}, 1)

	ws.HeartbeatTimer()
	ws.ReConnectTimer()
	ws.checkStatusTimer()

	return ws
}

func (ws *WsConn) connect() {
	dialer := websocket.DefaultDialer

	if ws.ProxyUrl != "" {
		proxy, err := url.Parse(ws.ProxyUrl)
		if err == nil {
			log.Println("proxy url :", proxy)
			dialer.Proxy = http.ProxyURL(proxy)
		} else {
			log.Println("proxy url error ? ", err)
		}
	}

	wsConn, resp, err := dialer.Dial(ws.WsUrl, http.Header(ws.ReqHeaders))
	if err != nil {
		panic(err)
	}

	ws.Conn = wsConn

	if ws.IsDump {
		dumpData, _ := httputil.DumpResponse(resp, true)
		log.Println(string(dumpData))
	}

	ws.UpdateActiveTime()
}

func (ws *WsConn) SendJsonMessage(v interface{}) error {
	ws.mu <- struct{}{}
	defer func() {
		<-ws.mu
	}()
	return ws.WriteJSON(v)
}

func (ws *WsConn) SendTextMessage(data []byte) error {
	ws.mu <- struct{}{}
	defer func() {
		<-ws.mu
	}()
	return ws.WriteMessage(websocket.TextMessage, data)
}

func (ws *WsConn) ReConnect() {
	ws.Lock()
	defer ws.Unlock()

	log.Println("close ws  error :", ws.Close())
	time.Sleep(time.Second)

	ws.connect()

	//re subscribe
	for _, sub := range ws.subs {
		log.Println("subscribe:", sub)
		ws.SendJsonMessage(sub)
	}
}

func (ws *WsConn) ReConnectTimer() {
	if ws.ReconnectIntervalTime == 0 {
		return
	}

	timer := time.NewTimer(ws.ReconnectIntervalTime)

	go func() {
		ws.clearChannel(ws.closeReconnect)

		for {
			select {
			case <-timer.C:
				log.Println("reconnect websocket")
				ws.ReConnect()
				timer.Reset(ws.ReconnectIntervalTime)
			case <-ws.closeReconnect:
				timer.Stop()
				log.Println("close websocket connect ,  exiting reconnect timer goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) checkStatusTimer() {
	if ws.HeartbeatIntervalTime == 0 {
		return
	}

	timer := time.NewTimer(ws.HeartbeatIntervalTime)

	go func() {
		ws.clearChannel(ws.closeCheck)

		for {
			select {
			case <-timer.C:
				now := time.Now()
				if now.Sub(ws.activeTime) >= 2*ws.HeartbeatIntervalTime {
					log.Println("active time [ ", ws.activeTime, " ] has expired , begin reconnect ws.")
					ws.ReConnect()
				}
				timer.Reset(ws.HeartbeatIntervalTime)
			case <-ws.closeCheck:
				log.Println("check status timer exiting")
				return
			}
		}
	}()
}

func (ws *WsConn) HeartbeatTimer() {
	log.Println("heartbeat interval time = ", ws.HeartbeatIntervalTime)
	if ws.HeartbeatIntervalTime == 0 {
		return
	}

	timer := time.NewTicker(ws.HeartbeatIntervalTime)
	go func() {
		ws.clearChannel(ws.closeHeartbeat)

		for {
			select {
			case <-timer.C:
				var err error
				if ws.HeartbeatFunc != nil {
					err = ws.SendJsonMessage(ws.HeartbeatFunc())
				} else {
					err = ws.SendTextMessage(ws.HeartbeatData)
				}
				if err != nil {
					log.Println("heartbeat error , ", err)
					time.Sleep(time.Second)
				}
			case <-ws.closeHeartbeat:
				timer.Stop()
				log.Println("close websocket connect , exiting heartbeat goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) Subscribe(subEvent interface{}) error {
	log.Println("Subscribe:", subEvent)
	err := ws.SendJsonMessage(subEvent)
	if err != nil {
		return err
	}
	ws.subs = append(ws.subs, subEvent)
	return nil
}

func (ws *WsConn) ReceiveMessage() {
	ws.clearChannel(ws.closeRecv)

	go func() {
		for {

			if len(ws.closeRecv) > 0 {
				<-ws.closeRecv
				log.Println("close websocket , exiting receive message goroutine.")
				return
			}

			t, msg, err := ws.ReadMessage()
			if err != nil {
				ws.ErrorHandleFunc(err)
				time.Sleep(time.Second)
				continue
			}

			switch t {
			case websocket.TextMessage:
				ws.ProtoHandleFunc(msg)
			case websocket.BinaryMessage:
				if ws.UnCompressFunc == nil {
					ws.ProtoHandleFunc(msg)
				} else {
					msg2, err := ws.UnCompressFunc(msg)
					if err != nil {
						ws.ErrorHandleFunc(fmt.Errorf("%s,%s", "un compress error", err.Error()))
					} else {
						err := ws.ProtoHandleFunc(msg2)
						if err != nil {
							ws.ErrorHandleFunc(err)
						}
					}
				}
			case websocket.CloseMessage:
				ws.CloseWs()
				return
			default:
				log.Println("error websocket message type , content is :\n", string(msg))
			}
		}
	}()
}

func (ws *WsConn) UpdateActiveTime() {
	ws.activeTimeL.Lock()
	defer ws.activeTimeL.Unlock()

	ws.activeTime = time.Now()
}

func (ws *WsConn) CloseWs() {
	ws.clearChannel(ws.closeCheck)
	ws.clearChannel(ws.closeReconnect)
	ws.clearChannel(ws.closeHeartbeat)
	ws.clearChannel(ws.closeRecv)

	ws.closeReconnect <- struct{}{}
	ws.closeHeartbeat <- struct{}{}
	ws.closeRecv <- struct{}{}
	ws.closeCheck <- struct{}{}

	err := ws.Close()
	if err != nil {
		log.Println("close websocket error , ", err)
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
