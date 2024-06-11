package binanceufutures

import (
	"fmt"
	"strings"
)

var (
	RestUrl  = "https://fapi.binance.com"
	PubWsUrl = "wss://fstream.binance.com/stream"
	PriWsUrl = "wss://fstream.binance.com/ws/"
)

// binanceF rest 接口url
const (
	OhlcvRest         = "/fapi/v1/klines?%s"    // limit: 大于1000,500-1000,100-499,小于100时 权重：10,5,2,1
	OrderBookRest     = "/fapi/v1/depth?%s"     // limit：1000,500,100,小于100时 权重：20,10,5,2
	SymbolsRest       = "/fapi/v1/exchangeInfo" // 权重 1
	TickerRest        = "/fapi/v1/ticker/24hr"  // 带 symbol：权重 1，不带 symbol：权重 40
	FetchHistoryTrade = "/fapi/v1/aggTrades"    // 权重 20
	TradeRest         = "/fapi/v1/trades?%s"    // 权重 5
	MarkPriceRest     = "/fapi/v1/premiumIndex" // 权重 1
	//私有接口
	CancelAllOrderRest  = "/fapi/v1/allOpenOrders" // 撤所有单接口 权重 1
	CancelOneOrderRest  = "/fapi/v1/order"         // 撤单接口 权重 1
	CancelMoreOrderRest = "/fapi/v1/batchOrders"   // 批量撤单接口 权重 1
	CreatOneOrderRest   = "/fapi/v1/order"         // 挂单接口 权重 0
	CreatMoreOrderRest  = "/fapi/v1/batchOrders"   // 批量挂单接口 权重 5
	BalanceRest         = "/fapi/v2/account"       // 权重 5
	OpenOrderRest       = "/fapi/v1/openOrders"    // 带 symbol：权重 1，不带 symbol：权重 40
	OrderRest           = "/fapi/v1/order"         // 权重 1
	PositionRest        = "/fapi/v2/positionRisk"  // 权重 5
	UserTread           = "/fapi/v1/userTrades"    // 权重 5
	RiskLimitRest       = "/fapi/v2/positionRisk"  // 权重 5
	FetchLeverage       = "/fapi/v2/account"
	SetLeverage         = "/fapi/v1/leverage"
	FetchFundingFeeRest = "/fapi/v1/income" // 权重 30
)

func Symbol2BinanceWsInstId(symbol string) string {
	tmp := strings.Split(symbol, "_")
	if len(tmp) == 2 {
		return fmt.Sprintf("%s%s", strings.ToLower(tmp[0]), strings.ToLower(tmp[1]))
	}
	panic("bad symbol:" + symbol)
}
