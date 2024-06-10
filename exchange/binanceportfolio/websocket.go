package binanceportfolio

import (
	"encoding/json"
	"time"

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
	client := ws.NewWsClient(url, imp, constant.OkxV5Spot, 20*time.Second, 30*time.Second)
	return client
}

func (binance *BinanceImp) Ping(cli *ws.WsClient) {
	log.Infof("ping")
	// cli.WriteBytes([]byte("ping"))
}
func (binance *BinanceImp) OnConnected(cli *ws.WsClient, typ ws.ConnectType) {
	if !binance.isPrivate {
		log.Info("binance portfolio public ws connected")
		return
	}
	log.Info("binance portfolio private ws connected")
	// ok.Login(cli)
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
				evt := &types.Order{
					Symbol:      l["s"].(string),               // 假设这是一个转换函数
					Type:        Binance2Type[l["o"].(string)], // 假设可以直接转换
					OrderID:     string(l["i"].(json.Number)),
					ClientID:    l["c"].(string),
					Side:        Binance2Side[l["S"].(string)], // 假设可以直接转换
					Price:       l["L"].(string),
					OrigQty:     l["q"].(string),
					ExecutedQty: l["z"].(string),
					AvgPrice:    l["ap"].(string),
					Fee:         l["n"].(string),
					Status:      Binance2Status[l["X"].(string)], // 假设可以直接转换
					CreateAt:    l["T"].(int64),                  // 示例：使用当前时间的毫秒表示
					UpdateAt:    time.Now().UnixNano() / int64(time.Millisecond),
				}
				binance.rspHandle([]*types.Order{evt})
			}
		}
	}
}
