package binanceportfolio

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type BinanceImp struct {
	accessKey  string
	secretKey  string
	passphrase string
	isPrivate  bool
	pingTimer  *time.Timer
	rspHandle  func(interface{})
}

func NewBinancePriWsClient(accessKey, secretKey, passphrase, listenKey string, rspHandle func(interface{})) *ws.WsClient {
	imp := &BinanceImp{
		accessKey:  accessKey,
		secretKey:  secretKey,
		passphrase: passphrase,
		rspHandle:  rspHandle,
		isPrivate:  true,
	}
	url := PriWsUrl + listenKey
	client := ws.NewWsClient(url, imp, constant.BinancePortfolio, 20*time.Second, 30*time.Second)
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
		log.Info("binance portfolio public ws connected")
		return
	}
	log.Info("binance portfolio private ws connected")
	// ok.Login(cli)
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
	dict, err := utils.ByteToMap(bs)
	if err != nil {
		log.Warnf("Binance Portfolio 的私有ws消息解析为map失败：%v，原始消息：%v。", err, bs)
		return
	}
	event, ok := dict["e"]
	if ok {
		switch event {
		case "ORDER_TRADE_UPDATE":
			l := dict["o"].(map[string]interface{})

			if l["x"] == "TRADE" {
				log.Infof("ORDER_TRADE_UPDATE %s", bs)

				CreateAt, err := json.Number(l["T"].(json.Number)).Int64()
				if err != nil {
					// handle error, perhaps set T to 0 or log an error message
					CreateAt = 0
				}

				evt := &types.Order{
					Symbol:      strings.Replace(l["s"].(string), "USDT", "_USDT", 1), // 假设这是一个转换函数
					Type:        Binance2Type[l["o"].(string)],                        // 假设可以直接转换
					OrderID:     string(l["i"].(json.Number)),
					ClientID:    l["c"].(string),
					Side:        Binance2Side[l["S"].(string)], // 假设可以直接转换
					Price:       l["L"].(string),
					OrigQty:     l["q"].(string),
					ExecutedQty: l["z"].(string),
					AvgPrice:    l["ap"].(string),
					Fee:         l["n"].(string),
					MarketType:  dict["fs"].(string),
					Status:      Binance2Status[l["X"].(string)], // 假设可以直接转换
					CreateAt:    CreateAt,                        // 示例：使用当前时间的毫秒表示
					UpdateAt:    time.Now().UnixNano() / int64(time.Millisecond),
				}
				binance.rspHandle([]*types.Order{evt})
			}
		}
	}
}
