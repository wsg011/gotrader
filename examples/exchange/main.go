package main

import (
	"fmt"
	"gotrader/exchange/okxv5swap"
)

func main() {
	wsClient := okxv5swap.NewOkPubWsClient()
	if err := wsClient.Dial(1); err != nil {
		fmt.Printf("Dial err %s", err)
		return
	}
	wsClient.Subscribe("BTC-USDT", "books")

	select {}
}
