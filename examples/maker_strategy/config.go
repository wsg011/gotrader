package main

import (
	"gotrader/trader"
	"gotrader/trader/types"
)

type Config struct {
	Symbol        string
	MakerExchange trader.Exchange
	HedgeExchange trader.Exchange
	IsMaker       bool
	IsTacker      bool
	Hedge         bool
}

type Vars struct {
	epoch           int64
	BookTicker      *types.BookTicker
	HedgeBookTicker *types.BookTicker

	basisMean   float64
	basisStd    float64
	fundingRate float64
}
