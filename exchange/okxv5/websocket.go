package okxv5

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"

	"github.com/bytedance/sonic"
)

type OkWsData struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	Arg   struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Data json.RawMessage `json:"data"`
}

type OkImp struct {
	accessKey  string
	secretKey  string
	passphrase string
	isPrivate  bool
	pingTimer  *time.Timer
	rspHandle  func(interface{})
}

func NewOkPubWsClient(rspHandle func(interface{})) *ws.WsClient {
	imp := &OkImp{rspHandle: rspHandle}
	client := ws.NewWsClient(PubWsUrl, imp, constant.OkxV5Spot, 20*time.Second, 30*time.Second)
	return client
}

func NewOkPriWsClient(accessKey, secretKey, passphrase string, rspHandle func(interface{})) *ws.WsClient {
	imp := &OkImp{
		accessKey:  accessKey,
		secretKey:  secretKey,
		passphrase: passphrase,
		rspHandle:  rspHandle,
		isPrivate:  true,
	}
	client := ws.NewWsClient(PriWsUrl, imp, constant.OkxV5Spot, 20*time.Second, 30*time.Second)
	return client
}

func (ok *OkImp) Ping(cli *ws.WsClient) {
	// log.Infof("ping")
	cli.WriteBytes([]byte("ping"))
}
func (ok *OkImp) OnConnected(cli *ws.WsClient, typ ws.ConnectType) {
	if !ok.isPrivate {
		log.Info("ok public ws connected")
		return
	}
	log.Info("ok private ws connected")
	ok.Login(cli)
}

func (ok *OkImp) Handle(cli *ws.WsClient, bs []byte) {
	if ok.pingTimer == nil {
		ok.pingTimer = time.AfterFunc(time.Second*20, func() {
			cli.WriteBytes([]byte("ping"))
		})
	}

	if len(bs) == 4 && string(bs) == "pong" {
		// log.Infof("RecvPongTime %s", time.Now())
		cli.SetRecvPongTime(time.Now())
		return
	}

	var dat OkWsData
	if err := sonic.Unmarshal(bs, &dat); err != nil {
		log.WithError(err).Error("unmarshal ok ws data failed")
		return
	}

	if (dat.Code != "" && dat.Code != "0") || dat.Event == "error" {
		err := fmt.Errorf("code:%s, msg:%s", dat.Code, dat.Msg)
		log.WithError(err).Error("ok ws data error")
		return
	}

	if dat.Event == "subscribe" {
		log.WithField("arg", dat.Arg).Info("ok subscribe success")
		return
	}

	if dat.Event == "login" {
		log.WithField("Event", dat.Event).Info("ok login success")
		return
	}

	switch dat.Arg.Channel {
	case "bbo-tbt":
		ok.onBboTbtRecv(dat.Arg.InstId, dat.Data)
	case "orders":
		ok.onOrders(dat.Arg.InstId, dat.Data)
	default:
		log.WithField("dat", string(dat.Data)).Warn("unknown ok message")
	}
}

func (ok *OkImp) Login(cli *ws.WsClient) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := generateOkxSignature(timestamp, ok.secretKey)

	okxReq := map[string]interface{}{
		"op": "login",
		"args": []map[string]string{
			{
				"apiKey":     ok.accessKey,
				"passphrase": ok.passphrase,
				"timestamp":  timestamp,
				"sign":       signature,
			},
		},
	}

	// 发送请求
	cli.Write(okxReq)
	log.Infof("okx login")
}

func (ok *OkImp) onBboTbtRecv(instId string, dat json.RawMessage) {
	type bookTicker struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		Ts        string     `json:"ts"`
		Checksum  int        `json:"checksum"`
		PrevSeqID int        `json:"prevSeqId"`
		SeqID     int        `json:"seqId"`
	}

	tickers := make([]bookTicker, 0, 1)
	if err := sonic.Unmarshal(dat, &tickers); err != nil {
		log.WithError(err).Error("unmarshal ok tbt failed")
		return
	}

	if len(tickers) == 0 {
		log.Warn("empty ok tbt")
		return
	}

	var (
		ticker      = tickers[0]
		ask1        = ticker.Asks[0]
		askPrice, _ = strconv.ParseFloat(ask1[0], 64)
		askQty, _   = strconv.ParseFloat(ask1[1], 64)
		bid1        = ticker.Bids[0]
		bidPrice, _ = strconv.ParseFloat(bid1[0], 64)
		bidQty, _   = strconv.ParseFloat(bid1[1], 64)
		ts, _       = strconv.ParseInt(ticker.Ts, 10, 64)
		exchange    = constant.OkxV5Spot
	)

	if strings.Contains(instId, "-SWAP") {
		exchange = constant.OkxV5Swap
	}

	evt := &types.BookTicker{
		Symbol:     OkInstId2Symbol(instId),
		Exchange:   exchange,
		AskPrice:   askPrice,
		AskQty:     askQty,
		BidPrice:   bidPrice,
		BidQty:     bidQty,
		ExchangeTs: ts * 1000,
		TraceId:    utils.RandomString(8),
		Ts:         utils.Microsec(time.Now()),
	}

	ok.rspHandle(evt)
}

func (ok *OkImp) onOrders(instId string, dat json.RawMessage) {
	type Order struct {
		Symbol         string `json:"instId"`
		OrderType      string `json:"ordType"`
		OrderID        string `json:"ordId"`
		ClientID       string `json:"clOrdId"`
		Side           string `json:"side"`
		Price          string `json:"px"`
		Quantity       string `json:"sz"`
		ExecutedQty    string `json:"accFillSz"`
		ExecutedAmount string `json:"fillNotionalUsd"`
		AvgPrice       string `json:"avgPx"`
		Fee            string `json:"fee"`
		Status         string `json:"state"`
		CreateAt       string `json:"cTime"`
		UpdateAt       string `json:"uTime"`
	}

	var orders []Order
	if err := sonic.Unmarshal(dat, &orders); err != nil {
		log.WithError(err).Error("unmarshal ok tbt failed")
		return
	}

	var result []*types.Order
	for _, ord := range orders {
		orderSide := Okx2Side[ord.OrderType]
		orderType := Okx2Type[ord.OrderType]
		orderSatus := Okex2Status[ord.Status]
		evt := &types.Order{
			Symbol:      OkInstId2Symbol(ord.Symbol), // 假设这是一个转换函数
			Type:        convertOrderType(orderType), // 假设可以直接转换
			OrderID:     ord.OrderID,
			ClientID:    ord.ClientID,
			Side:        convertOrderSide(orderSide), // 假设可以直接转换
			Price:       ord.Price,
			OrigQty:     ord.Quantity,
			ExecutedQty: ord.ExecutedQty,
			AvgPrice:    ord.AvgPrice,
			Fee:         ord.Fee,
			Status:      orderSatus,                                      // 假设可以直接转换
			CreateAt:    time.Now().UnixNano() / int64(time.Millisecond), // 示例：使用当前时间的毫秒表示
			UpdateAt:    time.Now().UnixNano() / int64(time.Millisecond),
		}
		result = append(result, evt)

	}
	ok.rspHandle(result)
}

func convertOrderSide(side string) constant.OrderSide {
	switch side {
	case "BUY":
		return constant.OrderBuy
	case "SELL":
		return constant.OrderSell
	default:
		log.WithField("side", side).Warn("unknown okx order side")
		return constant.OrderSell
	}
}

func convertOrderType(typ string) constant.OrderType {
	switch typ {
	case "MARKET":
		return constant.Market
	case "GTC":
		return constant.GTC
	case "IOC":
		return constant.IOC
	case "FOK":
		return constant.FOK
	case "POST_ONLY":
		return constant.PostOnly
	default:
		return constant.Limit // 假设默认为Limit类型
	}
}

// func convertOrderStatus(std string) constant.OrderStatus {
// 	// 订单状态
// 	// canceled：撤单成功
// 	// live：等待成交
// 	// partially_filled：部分成交
// 	// filled：完全成交
// 	switch std {
// 	case "OPEN":
// 		return constant.OrderOpen
// 	case "CLOSED":
// 		return constant.OrderClosed
// 	default:
// 		return constant.OrderOpen
// 	}
// }
