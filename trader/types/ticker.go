package types

type Ticker struct {
	Symbol     string
	MarketType string
	Open       float64
	High       float64
	Low        float64
	Vol        float64
	LastPrice  float64
	LastSize   float64
	AskPrice   float64
	AskSize    float64
	BidPrice   float64
	BidSize    float64
	ExchangeTs int64
	Ts         int64
}
