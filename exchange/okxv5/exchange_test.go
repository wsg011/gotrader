package okxv5

import (
	"testing"
	"time"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

var (
	apiKey      = ""
	secretKey   = ""
	passphrase  = ""
	symbol      = "APE_USDT_SWAP"
	hedgeSymbol = "APE_USDT"
	askPrice    = 0.0
	bidPrice    = 0.0
)

var exchange *OkxV5Exchange

func TestMain(t *testing.T) {
	params := &types.ExchangeParameters{
		AccessKey:  apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
	}
	exchange = NewOkxV5Swap(params)
	name := exchange.GetName()
	t.Logf("init Exchang name %s", name)
}

func TestFetchTickers(t *testing.T) {
	resp, err := exchange.FetchTickers()
	if err != nil {
		t.Fatalf("HttpRequest failed: %v", err)
	}

	for _, ticker := range resp {
		if ticker.Symbol == symbol {
			bidPrice = ticker.BidPrice
			askPrice = ticker.AskPrice
			t.Logf("%s ticker [%f | %f]", symbol, ticker.AskPrice, ticker.BidPrice)
		}
	}

	t.Logf("Tickers %v", resp[0])
}

func TestCreateBatchOrders(t *testing.T) {
	now := time.Now()
	// 5 合约 order
	orders := make([]*types.Order, 0)
	for i := 0; i < 5; i++ {
		price := askPrice + 0.01*float64(i+1)
		orderinfo := &types.Order{
			Symbol:       symbol,
			Side:         constant.OrderSell,
			Type:         constant.Limit,
			ClientID:     "",
			Price:        utils.FormatFloat(price, 3),
			OrigQty:      "1",
			ExchangeType: exchange.exchangeType,
			CreateAt:     utils.Millisec(now),
			Status:       constant.OrderSubmit,
		}
		orders = append(orders, orderinfo)
		// t.Logf("create order %v", orderinfo)
	}
	// 5 现货 order
	for i := 0; i < 5; i++ {
		price := bidPrice - 0.01*float64(i+1)
		orderinfo := &types.Order{
			Symbol:       hedgeSymbol,
			Side:         constant.OrderBuy,
			Type:         constant.Limit,
			ClientID:     "",
			Price:        utils.FormatFloat(price, 3),
			OrigQty:      "1",
			ExchangeType: exchange.exchangeType,
			CreateAt:     utils.Millisec(now),
			Status:       constant.OrderSubmit,
		}
		orders = append(orders, orderinfo)
	}

	t.Logf("create batch orders %v", orders)
	result, err := exchange.CreateBatchOrders(orders)
	if err != nil {
		t.Errorf("CreateBatchOrders error: %s", err)
	}
	for _, orderRes := range result {
		t.Logf("orderRes %v", orderRes)
	}

}

// func TestPriWsClient(t *testing.T) {
// 	testRspHandle := func(data interface{}) {
// 		// t.Logf("testRspHandle data: %s", data)
// 	}
// 	okxWsClient := NewOkPriWsClient(apiKey, secretKey, passphrase, testRspHandle)
// 	if err := okxWsClient.Dial(ws.Connect); err != nil {
// 		t.Fatalf("Dial error: %s", err)
// 	}
// 	t.Log("Dial success")
// 	time.Sleep(3 * time.Second)

// 	param := map[string]interface{}{
// 		"op": "subscribe",
// 		"args": []map[string]string{
// 			{
// 				"channel": "bbo-tbt",
// 				"instId":  "BTC-USDT",
// 			},
// 		},
// 	}
// 	okxWsClient.Write(param)

// 	time.Sleep(3 * time.Second)
// 	okxWsClient.Close()
// }

// func TestExchange(t *testing.T) {
// 	params := &types.ExchangeParameters{
// 		AccessKey:  "",
// 		SecretKey:  "",
// 		Passphrase: "",
// 	}
// 	exchange := NewOkxV5Swap(params)
// 	name := exchange.GetName()
// 	t.Logf("init exchange %s", name)

// 	param := map[string]interface{}{
// 		"op": "subscribe",
// 		"args": []map[string]string{
// 			{
// 				"channel": "bbo-tbt",
// 				"instId":  "BTC-USDT",
// 			},
// 		},
// 	}

// 	exchange.Subscribe(param)
// }

// func TestSwapWs(t *testing.T) {
// 	testRspHandle := func(data interface{}) {
// 		// t.Logf("testRspHandle data: %s", data)
// 	}
// 	okxWsClient := NewOkPubWsClient(testRspHandle)
// 	if err := okxWsClient.Dial(ws.Connect); err != nil {
// 		t.Fatalf("Dial error: %s", err)
// 	}
// 	t.Log("Dial success")
// 	param := map[string]interface{}{
// 		"op": "subscribe",
// 		"args": []map[string]string{
// 			{
// 				"channel": "bbo-tbt",
// 				"instId":  "BTC-USDT",
// 			},
// 		},
// 	}
// 	okxWsClient.Write(param)

// 	time.Sleep(3 * time.Second)
// 	okxWsClient.Close()
// }
