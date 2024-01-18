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
}
