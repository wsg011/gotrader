package types

import "gotrader/trader/constant"

type Position struct {
	Symbol    string // BTC_USDT_PERP
	Exchange  constant.ExchangeType
	Size      float64
	AvgPrice  float64
	MarkPrice float64
	UnPnl     float64
	UpdateAt  int64
}
