package binancespot

import (
	"fmt"
	"time"

	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type BinanceSpotExchange struct {
	exchangeType constant.ExchangeType

	restClient  *RestClient
	pubWsClient *ws.WsClient
	priWsClient *ws.WsClient

	// callbacks
	onBooktickerCallback func(*types.BookTicker)
	onOrderCallback      func([]*types.Order)
}

func NewBinanceSpot(params *types.ExchangeParameters) *BinanceSpotExchange {
	apiKey := params.AccessKey
	secretKey := params.SecretKey
	passPhrase := params.Passphrase

	// new client
	client := NewRestClient(apiKey, secretKey, passPhrase, constant.BinanceSpot)
	exchange := &BinanceSpotExchange{
		exchangeType: constant.OkxV5Swap,
		restClient:   client,
	}
	// pubWsClient
	pubWsClient := NewBinanceSpotPubWsClient(exchange.OnPubWsHandle)
	if err := pubWsClient.Dial(ws.Connect); err != nil {
		log.Errorf("pubWsClient.Dial err %s", err)
	} else {
		exchange.pubWsClient = pubWsClient
		log.Infof("pubWsClient.Dial success")
	}
	// // priWsClient
	// if len(apiKey) > 0 {
	// 	priWsClient := NewOkPriWsClient(apiKey, secretKey, passPhrase, exchange.OnPriWsHandle)
	// 	if err := priWsClient.Dial(ws.Connect); err != nil {
	// 		log.Errorf("priWsClient.Dial err %s", err)
	// 	} else {
	// 		exchange.priWsClient = priWsClient
	// 		log.Infof("priWsClient.Dial success")
	// 	}
	// }
	return exchange
}

func (binance *BinanceSpotExchange) GetName() (name string) {
	return binance.exchangeType.Name()

}

func (binance *BinanceSpotExchange) GetType() (typ constant.ExchangeType) {
	return binance.exchangeType
}

func (binance *BinanceSpotExchange) FetchSymbols() ([]*types.SymbolInfo, error) {
	return binance.restClient.FetchSymbols()

}
func (binance *BinanceSpotExchange) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {

	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceSpotExchange) FetchBalance() (*types.Assets, error) {
	return binance.restClient.FetchBalance()
}

func (binance *BinanceSpotExchange) CreateBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {

	return binance.restClient.CreateBatchOrders(orders)
}

func (binance *BinanceSpotExchange) Subscribe(params map[string]interface{}) error {
	// if okx.puWsClient == nil {
	// 	return fmt.Errorf("pubWsClient is nil")
	// }
	// if err := okx.puWsClient.Write(params); err != nil {
	// 	return fmt.Errorf("Subscribe err: %s", err)
	// }
	return nil
}

func (binance *BinanceSpotExchange) SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error) {
	for _, symbol := range symbols {
		var args []string
		args = append(args, fmt.Sprintf("%s@bookTicker", Symbol2BinanceWsInstId(symbol)))

		params := map[string]interface{}{
			"method": "SUBSCRIBE",
			"params": args,
			"id":     1,
		}

		if binance.pubWsClient == nil {
			return fmt.Errorf("pubWsClient is nil")
		}

		if err := binance.pubWsClient.Write(params); err != nil {
			return fmt.Errorf("Subscribe err: %s", err)
		}
		binance.onBooktickerCallback = callback

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

func (binance *BinanceSpotExchange) SubscribeOrders(symbols []string, callback func(orders []*types.Order)) (err error) {

	return fmt.Errorf("not impl")
}

func (binance *BinanceSpotExchange) OnPubWsHandle(data interface{}) {
	switch v := data.(type) {
	case *types.BookTicker:
		// callback
		if binance.onBooktickerCallback != nil {
			binance.onBooktickerCallback(v)
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
