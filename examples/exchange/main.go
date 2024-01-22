package main

import (
	"github.com/wsg011/gotrader/exchange/okxv5"
	"github.com/wsg011/gotrader/pkg/utils"
	"time"

	"github.com/wsg011/gotrader/trader/types"

	"github.com/sirupsen/logrus"
)

func main() {
	timeFormat := "20060102 15:04:05.999"
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: timeFormat})
	var log = logrus.WithField("package", "exchange")

	// rspHandle := func(data interface{}) {
	// 	// t.Logf("testRspHandle data: %s", data)
	// }
	// wsClient := okxv5swap.NewOkPubWsClient(rspHandle)
	// if err := wsClient.Dial(1); err != nil {
	// 	fmt.Printf("Dial err %s", err)
	// 	return
	// }

	// param := map[string]interface{}{
	// 	"op": "subscribe",
	// 	"args": []map[string]string{
	// 		{
	// 			"channel": "bbo-tbt",
	// 			"instId":  "BTC-USDT",
	// 		},
	// 	},
	// }
	// wsClient.Subscribe(param)

	epoch := 0
	onBookTickerHandle := func(bookticker *types.BookTicker) {
		epoch += 1
		// log.Infof("onBookTickerHandle %v", bookticker)
		if epoch%100 == 0 {
			amount := bookticker.AskPrice * bookticker.AskQty
			amount += 1

			processDelay := utils.Microsec(time.Now()) - bookticker.Ts
			feedDelay := bookticker.Ts - bookticker.ExchangeTs
			log.Infof("%-16s processDelay %v feedDelay %v us", bookticker.Symbol, processDelay, feedDelay)
		}

	}
	params := &types.ExchangeParameters{}
	exchange := okxv5.NewOkxV5Swap(params)
	symbols := []string{"BTC_USDT", "BTC_USDT_SWAP", "ETH_USDT", "ETH_USDT_SWAP", "SOL_USDT", "SOL_USDT_SWAP"}
	err := exchange.SubscribeBookTicker(symbols, onBookTickerHandle)
	if err != nil {
		log.Errorf("SubscribeBookticker err %s", err)
		return
	}
	select {}
}
