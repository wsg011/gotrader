package types

import (
	"fmt"
	"github.com/wsg011/gotrader/trader/constant"
	"hash/fnv"
)

type BookTicker struct {
	Symbol     string
	Exchange   constant.ExchangeType
	AskPrice   float64
	AskQty     float64
	BidPrice   float64
	BidQty     float64
	ExchangeTs int64
	Ts         int64
	TraceId    string
}

func (e *BookTicker) Hash() uint32 {
	h := fnv.New32a()
	s := fmt.Sprintf("%f%f%f%f%d", e.AskPrice, e.AskQty, e.BidPrice, e.BidQty, e.ExchangeTs)
	h.Write([]byte(s))
	return h.Sum32()
}

type OrderBookItem struct {
	Price float64
	Qty   float64
}

type OrderBook struct {
	Symbol     string
	Exchange   constant.ExchangeType
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	ExchangeTs int64
	Ts         int64
	TraceId    string
}

type Trade struct {
	Symbol       string
	ClientID     string
	Exchange     constant.ExchangeType
	Size         string
	FilledSize   string
	FilledAmount string
	Fee          string
	Side         constant.OrderSide
	Status       constant.OrderStatus
	Type         constant.OrderType
	ExchangeTs   int64
	Ts           int64
}
