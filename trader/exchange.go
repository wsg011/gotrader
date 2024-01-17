package trader

import "gotrader/trader/types"

type Exchange interface {
	GetName() (name string)
	Subscribe(params map[string]interface{}) (err error)
	SubscribeBookTicker(symbols []string, callback func(*types.BookTicker)) (err error)
}
