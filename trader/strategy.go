package trader

import "gotrader/trader/types"

type Strategy interface {
	GetName() string
	OnBookTicker(bookticker types.BookTicker)
	OnOrderBook(orderbook types.OrderBook)
	OnTrade(trade types.Trade)
	OnOrder(order types.Order)
	Run()
	Start()
	Close()
}
