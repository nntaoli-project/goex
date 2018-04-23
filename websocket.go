package goex

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsConn struct {
	*websocket.Conn
	url                      string
	heartbeatIntervalTime    time.Duration
	checkConnectIntervalTime time.Duration
	actived                  time.Time
	close                    chan int
	isClose                  bool
	subs                     []interface{}
}

const (
	SUB_TICKER      = 1 + iota
	SUB_ORDERBOOK
	SUB_KLINE_1M
	SUB_KLINE_15M
	SUB_KLINE_30M
	SUB_KLINE_1D
	UNSUB_TICKER
	UNSUB_ORDERBOOK
)

func NewWsConn(wsurl string) *WsConn {
	wsConn, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		panic(err)
	}
	return &WsConn{Conn: wsConn, url: wsurl, actived: time.Now(), checkConnectIntervalTime: 30 * time.Second, close: make(chan int, 1)}
}

func (ws *WsConn) ReConnect() {

	timer := time.NewTimer(ws.checkConnectIntervalTime)
	go func() {
		for {
			select {
			case <-timer.C:
				if time.Now().Sub(ws.actived) >= 2*ws.checkConnectIntervalTime {
					ws.Close()
					log.Println("start reconnect websocket:", ws.url)
					wsConn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
					if err != nil {
						log.Println("reconnect fail ???")
					} else {
						ws.Conn = wsConn
						ws.actived = time.Now()
						//re subscribe
						for _, sub := range ws.subs {
							log.Println("subscribe:", sub)
							ws.WriteJSON(sub)
						}
					}
				}
				timer.Reset(ws.checkConnectIntervalTime)
			case <-ws.close:
				timer.Stop()
				log.Println("close websocket connect, exiting reconnect goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) Heartbeat(heartbeat func() interface{}, interval time.Duration) {
	ws.heartbeatIntervalTime = interval
	ws.checkConnectIntervalTime = 2 * ws.heartbeatIntervalTime

	timer := time.NewTimer(interval)
	go func() {
		for {
			select {
			case <-timer.C:
				err := ws.WriteJSON(heartbeat())
				if err != nil {
					log.Println("heartbeat error , ", err)
					time.Sleep(time.Second)
				}
				timer.Reset(interval)
			case <-ws.close:
				timer.Stop()
				log.Println("close websocket connect , exiting heartbeat goroutine.")
				return
			}
		}
	}()
}

func (ws *WsConn) Subscribe(subEvent interface{}) error {
	err := ws.WriteJSON(subEvent)
	if err != nil {
		return err
	}
	ws.subs = append(ws.subs, subEvent)
	return nil
}

func (ws *WsConn) ReceiveMessage(handle func(msg []byte)) {
	go func() {
		for {
			t, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				if ws.isClose {
					log.Println("exiting receive message goroutine.")
					break
				}
				time.Sleep(time.Second)
				continue
			}
			switch t {
			case websocket.TextMessage, websocket.BinaryMessage:
				handle(msg)
			case websocket.PongMessage:
				ws.actived = time.Now()
			case websocket.CloseMessage:
				ws.CloseWs()
				return
			default:
				log.Println("error websocket message type , content is :\n", string(msg))
			}
		}
	}()
}

func (ws *WsConn) UpdateActivedTime() {
	ws.actived = time.Now()
}

func (ws *WsConn) CloseWs() {
	ws.close <- 1 //exit reconnect goroutine
	if ws.heartbeatIntervalTime > 0 {
		ws.close <- 1 //exit heartbeat goroutine
	}

	err := ws.Close()
	if err != nil {
		log.Println("close websocket connect error , ", err)
	}

	ws.isClose = true
}
