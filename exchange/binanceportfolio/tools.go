package binanceportfolio

import (
	"strings"

	"github.com/wsg011/gotrader/trader/constant"
)

var (
	RestUrl  = "https://papi.binance.com"
	PubWsUrl = "wss://fstream.binance.com/ws//public"
	PriWsUrl = "wss://fstream.binance.com/ws/private"

	BinanceOrderSide = map[string]string{
		constant.OrderBuy.Name():   "BUY",
		constant.OrderSell.Name():  "SELL",
		constant.Long.Name():       "LONG",
		constant.Short.Name():      "SHORT",
		constant.CloseLong.Name():  "CLOSE_LONG",
		constant.CloseShort.Name(): "CLOSE_SHORT",
		constant.All.Name():        "ALL",
	}
	BinanceOrderType = map[string]string{
		constant.Limit.Name():    "LIMIT",
		constant.Market.Name():   "MARKET",
		constant.GTC.Name():      "GTC",
		constant.IOC.Name():      "IOC",
		constant.FOK.Name():      "FOK",
		constant.PostOnly.Name(): "POST_ONLY",
	}
	Side2Binance = map[string]string{
		constant.OrderBuy.Name():  "BUY",
		constant.OrderSell.Name(): "SELL",
	}

	Type2Binance = map[string]string{
		constant.Limit.Name():    "LIMIT",
		constant.Market.Name():   "MARKET",
		constant.GTC.Name():      "GTC",
		constant.FOK.Name():      "IOC",
		constant.IOC.Name():      "IOC",
		constant.PostOnly.Name(): "POST_ONLY",
	}
)

const (
	FetchBalanceUri   = "/papi/v1/balance"
	FetchPositionsUri = "/papi/v1/um/positionRisk"
	CreateOrderUri    = "/papi/v1/um/order"
)

func Symbol2Binance(symbol string) string {
	return strings.Replace(symbol, "_", "", -1)
}
