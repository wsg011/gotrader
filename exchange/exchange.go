package exchange

import (
	"fmt"

	"github.com/wsg011/gotrader/exchange/okxv5"
	"github.com/wsg011/gotrader/trader"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

func NewExchange(exchangeType constant.ExchangeType, params *types.ExchangeParameters) trader.Exchange {
	switch exchangeType {
	case constant.OkxV5Swap:
		return okxv5.NewOkxV5Swap(params)
	case constant.OkxV5Spot:
		return okxv5.NewOkxV5Spot(params)
	default:
		panic(fmt.Sprintf("new exchange error [%v]", exchangeType))
	}
}
