package binancespot

import (
	"fmt"

	"github.com/wsg011/gotrader/trader/types"
)

func (binance *BinanceSpotExchange) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {

	return nil, fmt.Errorf("not impl")
}
