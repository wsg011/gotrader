package okxv5

import (
	"fmt"
	"gotrader/pkg/ws"
	"gotrader/trader/types"
	"time"
)

type OkxV5Exchange struct {
	client   *RestClient
	wsClient *ws.WsClient

	// callbacks
	onBooktickerCallback func(*types.BookTicker)
	onOrderCallback      func(*types.Order)
}

func NewOkxV5Swap(apiKey, secretKey, passPhrase string) *OkxV5Exchange {
	client := NewRestClient(apiKey, secretKey, passPhrase)
	exchange := &OkxV5Exchange{
		client: client,
	}
	pubWsClient := NewOkPubWsClient(exchange.OnPubWsHandle)
	if err := pubWsClient.Dial(ws.Connect); err != nil {
		log.Errorf("pubWsClient.Dial err %s", err)
	} else {
		exchange.wsClient = pubWsClient
		log.Infof("pubWsClient.Dial success")
	}
	return exchange
}

func (okx *OkxV5Exchange) GetName() (name string) {
	return "OkxV5Swap"

}

func (okx *OkxV5Exchange) GetTickers() {

}

func (okx *OkxV5Exchange) Subscribe(params map[string]interface{}) error {
	if okx.wsClient == nil {
		return fmt.Errorf("pubWsClient is nil")
	}
	if err := okx.wsClient.Write(params); err != nil {
		return fmt.Errorf("Subscribe err: %s", err)
	}
	return nil
}

func (okx *OkxV5Exchange) SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) error {
	okx.onBooktickerCallback = callback

	for _, symbol := range symbols {
		params := map[string]interface{}{
			"op": "subscribe",
			"args": []map[string]string{
				{
					"channel": "bbo-tbt",
					"instId":  Symbol2OkInstId(symbol),
				},
			},
		}
		if okx.wsClient == nil {
			return fmt.Errorf("pubWsClient is nil")
		}
		if err := okx.wsClient.Write(params); err != nil {
			return fmt.Errorf("Subscribe err: %s", err)
		}
		time.Sleep(200 * time.Microsecond)
	}
	return nil
}

// SubscribeOrder 订阅订单频道
func (okx *OkxV5Exchange) SubscribeOrder(symbols []string, callback func(*types.Order)) error {
	okx.onOrderCallback = callback

	params := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel":  "orders",
				"instType": "FUTURES",
			},
		},
	}
	if err := okx.wsClient.Write(params); err != nil {
		return fmt.Errorf("Subscribe err: %s", err)
	}
	return nil
}

func (okx *OkxV5Exchange) OnPubWsHandle(data interface{}) {
	switch v := data.(type) {
	case *types.BookTicker:
		// callback
		if okx.onBooktickerCallback != nil {
			okx.onBooktickerCallback(v)
		} else {
			log.Errorf("OnBookTicker Callback not set")
		}
	case *types.OrderBook:
		fmt.Println("OrderBook type", v)
	case *types.Trade:
		fmt.Println("Trade type", v)
	default:
		log.Errorf("Unknown type %s", v)
	}

}
