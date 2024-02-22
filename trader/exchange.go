package trader

import (
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type Exchange interface {
	GetName() (name string)
	GetType() (typ constant.ExchangeType)

	// rest Public
	FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error)
	FetchFundingRate(symbol string) (*types.FundingRate, error)
	FetchFundingRateHistory(symbol string, limit int64) ([]*types.FundingRate, error)
	FetchSymbols() ([]*types.SymbolInfo, error)

	// rest Private
	FetchBalance() (*types.Assets, error)
	FetchPositons() ([]*types.Position, error)
	CreateBatchOrders([]*types.Order) ([]*types.OrderResult, error)
	CancelBatchOrders(orders []*types.Order) ([]*types.OrderResult, error)

	// ws
	Subscribe(params map[string]interface{}) (err error)
	SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error)
	SubscribeOrders(symbols []string, callback func(orders []*types.Order)) (err error)
}
