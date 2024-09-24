package binancespot

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type BinanceWsData struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

type BinanceImp struct {
	accessKey  string
	secretKey  string
	passphrase string
	isPrivate  bool
	pingTimer  *time.Timer
	rspHandle  func(interface{})
}

func NewBinanceSpotPubWsClient(rspHandle func(interface{})) *ws.WsClient {
	imp := &BinanceImp{rspHandle: rspHandle}
	client := ws.NewWsClient(PubWsUrl, imp, constant.BinanceSpot, 20*time.Second, 30*time.Second)

	return client
}

func (binance *BinanceImp) Ping(cli *ws.WsClient) {
	deadline := time.Now().Add(10 * time.Second)
	err := cli.Conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
	if err != nil {
		log.Errorf("ping error %s", err)
		return
	}

	log.Infof("ping %s", deadline)
}

func (binance *BinanceImp) OnConnected(cli *ws.WsClient, typ ws.ConnectType) {
	if !binance.isPrivate {
		log.Info("binance spot public ws connected")
		return
	}
	log.Info("binance spot private ws connected")
	// ok.Login(cli)
	// keepAlive(cli.Conn, 20*time.Second)
}

func (binance *BinanceImp) Subscribe(symbol string, topic string) map[string]interface{} {
	params := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "bbo-tbt",
				"instId":  Symbol2BinanceWsInstId(symbol),
			},
		},
	}
	return params
}

func (binance *BinanceImp) Handle(cli *ws.WsClient, bs []byte) {
	var dat BinanceWsData
	if err := sonic.Unmarshal(bs, &dat); err != nil {
		log.WithError(err).Error("unmarshal ok ws data failed, bs", bs)
		return
	}

	parts := strings.Split(dat.Stream, "@")
	if len(parts) != 2 {
		log.Errorf("Stream format is incorrect %s", dat.Stream)
		return
	}

	channel := parts[1]
	switch channel {
	case "bookTicker":
		binance.onBboTbtRecv(dat.Data)
	}
}

func (binance *BinanceImp) onBboTbtRecv(data json.RawMessage) {
	/***
	{
	  "u":400900217,     // order book updateId
	  "s":"BNBUSDT",     // 交易对
	  "b":"25.35190000", // 买单最优挂单价格
	  "B":"31.21000000", // 买单最优挂单数量
	  "a":"25.36520000", // 卖单最优挂单价格
	  "A":"40.66000000"  // 卖单最优挂单数量
	}
	***/

	type bookTicker struct {
		UpdateID int64  `json:"u"`
		Symbol   string `json:"s"`
		BidPrice string `json:"b"`
		BidSize  string `json:"B"`
		AskPrice string `json:"a"`
		AskSize  string `json:"A"`
	}

	var ticker bookTicker
	if err := sonic.Unmarshal(data, &ticker); err != nil {
		log.WithError(err).Error("unmarshal binance tbt failed")
		return
	}

	askPirce, _ := utils.ParseFloat(ticker.AskPrice)
	askSize, _ := utils.ParseFloat(ticker.AskSize)
	bidPirce, _ := utils.ParseFloat(ticker.BidPrice)
	bidSize, _ := utils.ParseFloat(ticker.BidSize)

	evt := &types.BookTicker{
		Symbol:     strings.Replace(ticker.Symbol, "USDT", "_USDT", 1),
		Exchange:   constant.BinanceSpot,
		AskPrice:   askPirce,
		AskQty:     askSize,
		BidPrice:   bidPirce,
		BidQty:     bidSize,
		ExchangeTs: utils.Microsec(time.Now()),
		TraceId:    utils.RandomString(8),
		Ts:         utils.Microsec(time.Now()),
	}

	binance.rspHandle(evt)
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				return
			}
		}
	}()
}
