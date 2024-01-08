package ws

import (
	"sync"
	"time"

	"gotrader/trader/constant"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	imp      WsImp
	conn     *websocket.Conn
	url      string
	wch      chan []byte
	exchange constant.ExchangeType
	priv     bool

	recvPingTime time.Time
	recvPongTime time.Time
	pingInterval time.Duration
	pongTimeout  time.Duration

	mutex  sync.Mutex
	quit   chan struct{}
	closed bool
	epoch  int64
}
