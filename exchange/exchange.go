package exchange

import (
	"fmt"
	"gotrader/exchange/okxv5swap"
	"gotrader/trader"
	"gotrader/trader/constant"
)

func NewExchange(exchangeType constant.ExchangeType) trader.Exchange {
	switch exchangeType {
	case constant.OkxV5Swap:
		return okxv5swap.NewOkxV5Swap("", "", "")
	default:
		panic(fmt.Sprintf("new exchange error [%v]", exchangeType))
	}
}
