package binanceufutures

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

func NewBinanceUFuturesPubWsClient(rspHandle func(interface{})) *ws.WsClient {
	imp := &BinanceImp{rspHandle: rspHandle}
	client := ws.NewWsClient(PubWsUrl, imp, constant.BinanceUFutures, 20*time.Second, 30*time.Second)
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
		log.Info("binance ufutures public ws connected")
		return
	}
	log.Info("binance ufutures private ws connected")
	// ok.Login(cli)
}

func (binance *BinanceImp) Handle(cli *ws.WsClient, bs []byte) {
	var dat BinanceWsData
	if err := sonic.Unmarshal(bs, &dat); err != nil {
		log.WithError(err).Error("unmarshal ok ws data failed, bs", bs)
		return
	}

	parts := strings.Split(dat.Stream, "@")
	if len(parts) != 2 {
		log.Errorf("Stream format is incorrect %s", dat)
		return
	}

	channel := parts[1]
	switch channel {
	case "bookTicker":
		binance.onBboTbtRecv(dat.Data)
	}
}

func (binance *BinanceImp) onBboTbtRecv(data json.RawMessage) {
	type bookTicker struct {
		Event     string `json:"e"`
		UpdateID  int64  `json:"u"`
		Symbol    string `json:"s"`
		BidPrice  string `json:"b"`
		BidSize   string `json:"B"`
		AskPrice  string `json:"a"`
		AskSize   string `json:"A"`
		Ts        int64  `json:"T"`
		EventTime int64  `json:"E"`
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
		Exchange:   constant.BinanceUFutures,
		AskPrice:   askPirce,
		AskQty:     askSize,
		BidPrice:   bidPirce,
		BidQty:     bidSize,
		ExchangeTs: ticker.Ts,
		TraceId:    utils.RandomString(8),
		Ts:         utils.Microsec(time.Now()),
	}

	binance.rspHandle(evt)
}
