package okxv5

import (
	"fmt"
	"time"

	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type OkxV5Exchange struct {
	exchangeType constant.ExchangeType

	restClient  *RestClient
	pubWsClient *ws.WsClient
	priWsClient *ws.WsClient

	// callbacks
	onBooktickerCallback func(*types.BookTicker)
	onOrderCallback      func([]*types.Order)
}

func NewOkxV5Swap(params *types.ExchangeParameters) *OkxV5Exchange {
	apiKey := params.AccessKey
	secretKey := params.SecretKey
	passPhrase := params.Passphrase

	// new client
	client := NewRestClient(apiKey, secretKey, passPhrase, constant.OkxV5Swap)
	exchange := &OkxV5Exchange{
		exchangeType: constant.OkxV5Swap,
		restClient:   client,
	}
	// pubWsClient
	pubWsClient := NewOkPubWsClient(exchange.OnPubWsHandle)
	if err := pubWsClient.Dial(ws.Connect); err != nil {
		log.Errorf("pubWsClient.Dial err %s", err)
	} else {
		exchange.pubWsClient = pubWsClient
		log.Infof("pubWsClient.Dial success")
	}
	// priWsClient
	if len(apiKey) > 0 {
		priWsClient := NewOkPriWsClient(apiKey, secretKey, passPhrase, exchange.OnPriWsHandle)
		if err := priWsClient.Dial(ws.Connect); err != nil {
			log.Errorf("priWsClient.Dial err %s", err)
		} else {
			exchange.priWsClient = priWsClient
			log.Infof("priWsClient.Dial success")
		}
	}
	return exchange
}

func NewOkxV5Spot(params *types.ExchangeParameters) *OkxV5Exchange {
	apiKey := params.AccessKey
	secretKey := params.SecretKey
	passPhrase := params.Passphrase

	// new client
	client := NewRestClient(apiKey, secretKey, passPhrase, constant.OkxV5Spot)
	exchange := &OkxV5Exchange{
		exchangeType: constant.OkxV5Spot,
		restClient:   client,
	}
	pubWsClient := NewOkPubWsClient(exchange.OnPubWsHandle)
	if err := pubWsClient.Dial(ws.Connect); err != nil {
		log.Errorf("pubWsClient.Dial err %s", err)
	} else {
		exchange.pubWsClient = pubWsClient
		log.Infof("pubWsClient.Dial success")
	}

	if len(apiKey) > 0 {
		priWsClient := NewOkPriWsClient(apiKey, secretKey, passPhrase, exchange.OnPriWsHandle)
		if err := priWsClient.Dial(ws.Connect); err != nil {
			log.Errorf("priWsClient.Dial err %s", err)

		} else {
			exchange.priWsClient = priWsClient
			log.Infof("priWsClient.Dial success")
		}
	}
	return exchange
}

func (okx *OkxV5Exchange) GetName() (name string) {
	return okx.exchangeType.Name()

}

func (okx *OkxV5Exchange) GetType() (typ constant.ExchangeType) {
	return okx.exchangeType
}

func (okx *OkxV5Exchange) FetchTickers() ([]types.BookTicker, error) {
	return okx.restClient.FetchTickers()
}

func (okx *OkxV5Exchange) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {
	return okx.restClient.FetchKline(symbol, interval, limit)
}

func (okx *OkxV5Exchange) FetchFundingRate(symbol string) (*types.FundingRate, error) {
	return okx.restClient.FetchFundingRate(symbol)
}

func (okx *OkxV5Exchange) FetchFundingRateHistory(symbol string, limit int64) ([]*types.FundingRate, error) {
	return okx.restClient.FetchFundingRateHistory(symbol, limit)
}

func (okx *OkxV5Exchange) FetchSymbols() ([]*types.SymbolInfo, error) {
	return okx.restClient.FetchSymbols()
}

func (okx *OkxV5Exchange) FetchBalance() (*types.Assets, error) {
	return okx.restClient.FetchBalance()
}

func (okx *OkxV5Exchange) FetchPositons() ([]*types.Position, error) {
	return okx.restClient.FetchPositons()
}

func (okx *OkxV5Exchange) CreateBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	return okx.restClient.CreateBatchOrders(orders)
}

func (okx *OkxV5Exchange) CancelBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	return okx.restClient.CancelBatchOrders(orders)
}

func (okx *OkxV5Exchange) Subscribe(params map[string]interface{}) error {
	// if okx.puWsClient == nil {
	// 	return fmt.Errorf("pubWsClient is nil")
	// }
	// if err := okx.puWsClient.Write(params); err != nil {
	// 	return fmt.Errorf("Subscribe err: %s", err)
	// }
	return nil
}

func (okx *OkxV5Exchange) SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) error {
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
		if okx.pubWsClient == nil {
			return fmt.Errorf("pubWsClient is nil")
		}
		if err := okx.pubWsClient.Write(params); err != nil {
			return fmt.Errorf("Subscribe err: %s", err)
		}
		time.Sleep(200 * time.Microsecond)
	}

	okx.onBooktickerCallback = callback
	return nil
}

// SubscribeOrder 订阅订单频道
func (okx *OkxV5Exchange) SubscribeOrders(symbols []string, callback func(orders []*types.Order)) error {
	/***
	{
		"op": "subscribe",
		"args": [{
			"channel": "orders",
			"instType": "FUTURES",
			"instFamily": "BTC-USD"
		}]
	}
	***/

	// 构建订阅请求的参数
	args := make([]map[string]string, 0)
	for _, symbol := range symbols {
		arg := map[string]string{
			"channel":  "orders",
			"instType": "SWAP",                  // 这里假设所有的symbol都是SWAP类型，根据需要调整
			"instId":   Symbol2OkInstId(symbol), // 为每个symbol设置instFamily
		}
		args = append(args, arg)
	}

	params := map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	}
	if err := okx.priWsClient.Write(params); err != nil {
		return fmt.Errorf("Subscribe err: %s", err)
	}

	okx.onOrderCallback = callback
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

func (okx *OkxV5Exchange) OnPriWsHandle(data interface{}) {
	switch v := data.(type) {
	case []*types.Order:
		if okx.onOrderCallback != nil {
			okx.onOrderCallback(v)
		} else {
			log.Errorf("onOrder Callback not set")
		}
	default:
		log.Errorf("Unknown type %s", v)
	}
}
