package trader

type Strategy interface {
	OnBookTicker(bookticker BookTicker)
	OnOrderBook(orderbook OrderBook)
	OnTrade(trade Trade)
	Pricing()
	Hedge()
}
