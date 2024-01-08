package okxv5swap

import (
	"fmt"
	"strings"
	"time"
)

var (
	RestUrl = "https://www.okex.com"
)

const (
	OhlcvRest               = "/api/v5/market/candles?%s"
	OrderBookRest           = "/api/v5/market/books?%s"
	SymbolsRest             = "/api/v5/public/instruments?instType=SWAP"
	TickerRest              = "/api/v5/market/ticker?%s"
	TickersRest             = "/api/v5/market/tickers?%s"
	TradeRest               = "/api/v5/market/trades?%s"
	FetchBalanceUri         = "/api/v5/account/balance"
	CreateSingleOrderUri    = "/api/v5/trade/order"
	CreateBatchOrderUri     = "/api/v5/trade/batch-orders"
	CancelSingleOrderUri    = "/api/v5/trade/cancel-order"
	CancelBatchOrderUri     = "/api/v5/trade/cancel-batch-orders"
	FetchOpenOrderUri       = "/api/v5/trade/orders-pending"
	FetchOrderWithIdUri     = "/api/v5/trade/order"
	FetchOrderDefault       = "/api/v5/trade/orders-history-archive"
	FetchUserTradesUri      = "/api/v5/trade/fills-history"
	PrivateTransferUri      = "/api/v5/asset/transfer"
	PrivateCurrenciesUri    = "/api/v5/asset/currencies"
	PrivateWithDrawUri      = "/api/v5/asset/withdrawal"
	FetchDepositHistoryUri  = "/api/v5/asset/deposit-history"
	FetchWithDrawHistoryUri = "/api/v5/asset/withdrawal-history"
	PrivateDepositAddrUri   = "/api/v5/asset/deposit-address"
	FetchTransferStateUri   = "/api/v5/asset/transfer-state"
	TransferProcessing      = 58124 //提币处理中返回此code
)

// IsoTime eg: 2018-03-16T18:02:48.284Z
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func OkInstId2Symbol(instId string) string {
	tmp := strings.Split(instId, "-")
	if len(tmp) == 2 {
		return fmt.Sprintf("%s_%s", tmp[0], tmp[1])
	} else if len(tmp) == 3 {
		return fmt.Sprintf("%s_%s", tmp[0], tmp[1])
	}
	panic("bad instId:" + instId)
}

func Symbol2OkInstId(symbol string) string {
	tmp := strings.Split(symbol, "_")
	if len(tmp) == 2 {
		return fmt.Sprintf("%s-%s", tmp[0], tmp[1])
	} else if len(tmp) == 3 {
		return fmt.Sprintf("%s-%s-SWAP", tmp[0], tmp[1])
	}
	panic("bad symbol:" + symbol)
}

func BaseQuote(symbol string) (string, string) {
	tmp := strings.Split(symbol, "_")
	return tmp[0], tmp[1]
}

func IsPerpSymbol(symbol string) bool {
	return strings.Contains(symbol, "_PERP")
}
