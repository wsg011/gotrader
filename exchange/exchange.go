package exchange

import (
	"fmt"
	"gotrader/exchange/okxv5"
	"gotrader/trader"
	"gotrader/trader/constant"
)

func NewExchange(exchangeType constant.ExchangeType) trader.Exchange {
	switch exchangeType {
	case constant.OkxV5Swap:
		return okxv5.NewOkxV5Swap("", "", "")
	case constant.OkxV5Spot:
		return okxv5.NewOkxV5Spot("", "", "")
	default:
		panic(fmt.Sprintf("new exchange error [%v]", exchangeType))
	}
}
