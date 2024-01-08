package trader

import "time"

type MarketDataInterface interface {
	DataType() string
	GetSymbol() string
}

type BookTicker struct {
	Symbol   string
	BidPrice float64
	AskPrice float64
	Time     time.Time
}

func (bt BookTicker) DataType() string {
	return "bookticker"
}

func (bt BookTicker) GetSymbol() string {
	return bt.Symbol
}

type OrderBook struct {
	Symbol string
	Bids   []OrderBookEntry
	Asks   []OrderBookEntry
	Time   time.Time
}

type OrderBookEntry struct {
	Price  float64
	Volume float64
}

func (ob OrderBook) DataType() string {
	return "orderbook"
}

func (ob OrderBook) GetSymbol() string {
	return ob.Symbol
}

type Trade struct {
	Symbol string
	Price  float64
	Volume float64
	Time   time.Time
}

func (t Trade) DataType() string {
	return "trade"
}

func (t Trade) GetSymbol() string {
	return t.Symbol
}
