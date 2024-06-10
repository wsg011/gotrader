package binanceportfolio

import (
	"fmt"

	"github.com/wsg011/gotrader/pkg/ws"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type BinancePortfolioExchange struct {
	exchangeType constant.ExchangeType

	restClient  *RestClient
	pubWsClient *ws.WsClient
	priWsClient *ws.WsClient

	// callbacks
	onBooktickerCallback func(*types.BookTicker)
	onOrderCallback      func([]*types.Order)
}

func NewBinancePortfoli(params *types.ExchangeParameters) *BinancePortfolioExchange {
	apiKey := params.AccessKey
	secretKey := params.SecretKey
	passPhrase := params.Passphrase

	// new client
	client := NewRestClient(apiKey, secretKey, passPhrase, constant.OkxV5Swap)
	exchange := &BinancePortfolioExchange{
		exchangeType: constant.OkxV5Swap,
		restClient:   client,
	}
	// pubWsClient
	// pubWsClient := NewOkPubWsClient(exchange.OnPubWsHandle)
	// if err := pubWsClient.Dial(ws.Connect); err != nil {
	// 	log.Errorf("pubWsClient.Dial err %s", err)
	// } else {
	// 	exchange.pubWsClient = pubWsClient
	// 	log.Infof("pubWsClient.Dial success")
	// }
	// // priWsClient
	if len(apiKey) > 0 {
		listenKey, err := exchange.GetListenKey()
		if err != nil {
			log.Errorf("GetListenKey err %s", err)
		} else {
			priWsClient := NewBinancePriWsClient(apiKey, secretKey, passPhrase, listenKey, exchange.OnPriWsHandle)
			if err := priWsClient.Dial(ws.Connect); err != nil {
				log.Errorf("priWsClient.Dial err %s", err)
			} else {
				exchange.priWsClient = priWsClient
				log.Infof("priWsClient.Dial success")
			}
			exchange.KeepUserStream(listenKey)
		}
	}
	return exchange
}

func (binance *BinancePortfolioExchange) GetName() (name string) {
	return binance.exchangeType.Name()

}

func (binance *BinancePortfolioExchange) GetType() (typ constant.ExchangeType) {
	return binance.exchangeType
}

func (binance *BinancePortfolioExchange) GetListenKey() (string, error) {
	return binance.restClient.GetListenKey()
}

func (binance *BinancePortfolioExchange) KeepUserStream(listenKey string) {
	go binance.restClient.KeepUserStream(listenKey)
}

func (binance *BinancePortfolioExchange) FetchSymbols() ([]*types.SymbolInfo, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) FetchBalance() (*types.Assets, error) {
	return binance.restClient.FetchBalance()
}

func (binance *BinancePortfolioExchange) CreateBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	return binance.restClient.CreateBatchOrders(orders)
}

func (binance *BinancePortfolioExchange) CreateUMOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	return binance.restClient.CreateBatchOrders(orders)
}

func (binance *BinancePortfolioExchange) CreateMMOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	return binance.restClient.CreateMMOrders(orders)
}

func (binance *BinancePortfolioExchange) CancelBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {

	return binance.restClient.CancelBatchOrders(orders)
}

func (binance *BinancePortfolioExchange) FetchTickers() ([]*types.Ticker, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) FetchFundingRate(symbol string) (*types.FundingRate, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) FetchFundingRateHistory(symbol string, limit int64) ([]*types.FundingRate, error) {
	return nil, fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) FetchPositons() ([]*types.Position, error) {
	return binance.restClient.FetchPositons()
}

func (binance *BinancePortfolioExchange) Subscribe(params map[string]interface{}) error {
	// if okx.puWsClient == nil {
	// 	return fmt.Errorf("pubWsClient is nil")
	// }
	// if err := okx.puWsClient.Write(params); err != nil {
	// 	return fmt.Errorf("Subscribe err: %s", err)
	// }
	return nil
}

func (binance *BinancePortfolioExchange) SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error) {

	return fmt.Errorf("not impl")
}

func (binance *BinancePortfolioExchange) SubscribeOrders(symbols []string, callback func(orders []*types.Order)) error {
	binance.onOrderCallback = callback
	return nil
}

func (binance *BinancePortfolioExchange) OnPriWsHandle(data interface{}) {
	switch v := data.(type) {
	case []*types.Order:
		if binance.onOrderCallback != nil {
			binance.onOrderCallback(v)
		} else {
			log.Errorf("onOrder Callback not set")
		}
	default:
		log.Errorf("Unknown type %s", v)
	}
}
