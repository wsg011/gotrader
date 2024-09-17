package binanceufutures

import (
	"fmt"
	"time"

	"github.com/wsg011/gotrader/exchange/base"
	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type BinanceUFuturesExchange struct {
	exchangeType constant.ExchangeType

	restClient  *RestClient
	pubWsClient *ws.WsClient
	priWsClient *ws.WsClient

	// callbacks
	onBooktickerCallback func(*types.BookTicker)
	onOrderCallback      func([]*types.Order)
}

func NewBinanceUFutures(params *types.ExchangeParameters) *BinanceUFuturesExchange {
	apiKey := params.AccessKey
	secretKey := params.SecretKey
	passPhrase := params.Passphrase

	// new client
	client := NewRestClient(apiKey, secretKey, passPhrase, constant.BinanceUFutures)
	exchange := &BinanceUFuturesExchange{
		exchangeType: constant.BinanceUFutures,
		restClient:   client,
	}
	// pubWsClient
	pubWsClient := NewBinanceUFuturesPubWsClient(exchange.OnPubWsHandle)
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

func (binance *BinanceUFuturesExchange) GetName() (name string) {
	return binance.exchangeType.Name()

}

func (binance *BinanceUFuturesExchange) GetType() (typ constant.ExchangeType) {
	return binance.exchangeType
}

func (binance *BinanceUFuturesExchange) FetchSymbols() ([]*types.SymbolInfo, error) {
	return binance.restClient.FetchSymbols()
}

func (binance *BinanceUFuturesExchange) FetchBalance() (*types.Assets, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) FetchAssetBalance() (*types.Assets, error) {
	return nil, fmt.Errorf("FetchAssetBalance not imp")
}

func (binance *BinanceUFuturesExchange) CreateBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {

	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) CancelBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {

	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) FetchTickers() ([]*types.Ticker, error) {
	return binance.restClient.FetchTickers()
}

func (binance *BinanceUFuturesExchange) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) FetchFundingRate(symbol string) (*types.FundingRate, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) FetchFundingRateHistory(symbol string, limit int64) ([]*types.FundingRate, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) FetchPositons() ([]*types.Position, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) PrivateTransfer(transfer base.TransferParam) (string, error) {
	return "", fmt.Errorf("PrivateTransfer not imp")
}

func (binance *BinanceUFuturesExchange) Subscribe(params map[string]interface{}) error {
	// if okx.puWsClient == nil {
	// 	return fmt.Errorf("pubWsClient is nil")
	// }
	// if err := okx.puWsClient.Write(params); err != nil {
	// 	return fmt.Errorf("Subscribe err: %s", err)
	// }
	return nil
}

func (binance *BinanceUFuturesExchange) SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error) {
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

func (binance *BinanceUFuturesExchange) SubscribeOrders(symbols []string, callback func(orders []*types.Order)) (err error) {

	return fmt.Errorf("not impl")
}

func (binance *BinanceUFuturesExchange) OnPubWsHandle(data interface{}) {
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
