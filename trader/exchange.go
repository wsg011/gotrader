package trader

import (
	"gotrader/trader/constant"
	"gotrader/trader/types"
)

type Exchange interface {
	GetName() (name string)
	GetType() (typ constant.ExchangeType)
	Subscribe(params map[string]interface{}) (err error)
	SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error)
	FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error)
	FetchFundingRate(symbol string) (*types.FundingRate, error)
	FetchBalance() (*types.Assets, error)
}
