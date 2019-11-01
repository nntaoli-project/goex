package goex

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	TargetName            string
	logger                *log.Logger
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
	return &WsBuilder{&WsConfig{
		ReqHeaders: make(map[string][]string, 1)}}
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

func (b *WsBuilder) TargetName(exName string) *WsBuilder {
	b.wsConfig.TargetName = exName

	return b
}

func (b *WsBuilder) Build() *WsConn {
	if b.wsConfig.TargetName != "" {
		b.wsConfig.logger = log.New(os.Stderr, "["+b.wsConfig.TargetName+"]", log.LstdFlags)
	} else {
		b.wsConfig.logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	if b.wsConfig.ErrorHandleFunc == nil {
		b.wsConfig.ErrorHandleFunc = func(err error) {
			b.wsConfig.logger.Println(err)
		}
	}
	wsConn := &WsConn{WsConfig: *b.wsConfig}
	return wsConn.NewWs()
}

func (ws *WsConn) NewWs() *WsConn {
	ws.Lock()
	defer ws.Unlock()

	if err := ws.connect(); err != nil {
		panic(err)
	}

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

func (ws *WsConn) connect() error {
	dialer := websocket.DefaultDialer

	if ws.ProxyUrl != "" {
		proxy, err := url.Parse(ws.ProxyUrl)
		if err == nil {
			ws.logger.Println("proxy url :", proxy)
			dialer.Proxy = http.ProxyURL(proxy)
		} else {
			ws.logger.Println("proxy url error ? ", err)
		}
	}

	wsConn, resp, err := dialer.Dial(ws.WsUrl, http.Header(ws.ReqHeaders))
	if err != nil {
		if ws.IsDump && resp != nil {
			dumpData, _ := httputil.DumpResponse(resp, true)
			ws.logger.Println(string(dumpData))
		}
		return err
	}

	ws.Conn = wsConn

	if ws.IsDump {
		dumpData, _ := httputil.DumpResponse(resp, true)
		ws.logger.Println(string(dumpData))
	}

	ws.UpdateActiveTime()

	return nil
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

	ws.logger.Println("close ws  error :", ws.Close())
	time.Sleep(time.Second)

	if err := ws.connect(); err != nil {
		ws.logger.Println(ws.WsUrl, "ws connect error ", err)
		return
	}

	//re subscribe
	for _, sub := range ws.subs {
		ws.logger.Println("subscribe:", sub)
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
				ws.logger.Println("reconnect websocket")
				ws.ReConnect()
				timer.Reset(ws.ReconnectIntervalTime)
			case <-ws.closeReconnect:
				timer.Stop()
				ws.logger.Println("close websocket connect ,  exiting reconnect timer goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) checkStatusTimer() {
	var checkStatusTimer *time.Ticker
	if ws.HeartbeatIntervalTime == 0 {
		checkStatusTimer = time.NewTicker(10 * time.Second)
	} else {
		checkStatusTimer = time.NewTicker(ws.HeartbeatIntervalTime)
	}

	go func() {
		ws.clearChannel(ws.closeCheck)

		for {
			select {
			case <-checkStatusTimer.C:
				now := time.Now()
				if now.Sub(ws.activeTime) >= 2*ws.HeartbeatIntervalTime {
					ws.logger.Println("active time [ ", ws.activeTime, " ] has expired , begin reconnect ws.")
					ws.ReConnect()
				}
			case <-ws.closeCheck:
				checkStatusTimer.Stop()
				ws.logger.Println("check status timer exiting")
				return
			}
		}
	}()
}

func (ws *WsConn) HeartbeatTimer() {
	ws.logger.Println("heartbeat interval time = ", ws.HeartbeatIntervalTime)
	if ws.HeartbeatIntervalTime == 0 || (ws.HeartbeatFunc == nil && ws.HeartbeatData == nil) {
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
					ws.logger.Println("heartbeat error , ", err)
					time.Sleep(time.Second)
				}
			case <-ws.closeHeartbeat:
				timer.Stop()
				ws.logger.Println("close websocket connect , exiting heartbeat goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) Subscribe(subEvent interface{}) error {
	ws.logger.Println("Subscribe:", subEvent)
	err := ws.SendJsonMessage(subEvent)
	if err != nil {
		return err
	}
	ws.subs = append(ws.subs, subEvent)
	return nil
}
func (ws *WsConn) messageHandler() {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			ws.logger.Printf("websocket ReceiveMessage err:%v, recover it\n", err)
			ws.messageHandler()
		}
	}()
	for {

		if len(ws.closeRecv) > 0 {
			<-ws.closeRecv
			ws.logger.Println("close websocket , exiting receive message goroutine.")
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
			ws.logger.Println("error websocket message type , content is :\n", string(msg))
		}
	}
}
func (ws *WsConn) ReceiveMessage() {
	ws.clearChannel(ws.closeRecv)
	//exit
	ws.SetCloseHandler(func(code int, text string) error {
		ws.logger.Println("websocket exiting ,code=", code, ",text=", text)
		ws.CloseWs()
		return nil
	})

	go ws.messageHandler()
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
		ws.logger.Println("close websocket error , ", err)
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
