package binancespot

var (
	RestUrl  = "https://api.binance.com"
	PubWsUrl = "wss://stream.binance.com:9443/public"
	PriWsUrl = "wss://stream.binance.com:9443/private"
)

const (
	FetchSymbolUri    = "/api/v3/exchangeInfo" // 权重(IP): 20
	FetchOrderBookUri = "/api/v3/depth"        // 权重(IP):1-100,101-500,501-1000,1001-5000 对应权重 5,25,50,250
	FetchTradesUri    = "/api/v3/trades"       // 权重(IP): 10
	FetchHistoryTrade = "/api/v3/aggTrades"    // 权重(IP): 2
	FetchOhlcvUri     = "/api/v3/klines"       // 权重(IP): 2
	FetchTickerUri    = "/api/v3/ticker/24hr"  // 权重(IP): 有 symbol:2,没有 symbol:80,传 symbols 有梯度权重

	FetchBalanceUri         = "/api/v3/account"                        // 权重(IP): 20
	CancelAllOrderUri       = "/api/v3/openOrders"                     // 权重(IP): 1
	CreateOrderUri          = "/api/v3/order"                          // 权重(UID): 1 权重(IP): 1
	CancelOrderUri          = "/api/v3/order"                          // 权重(IP): 1
	FetchOpenOrderUri       = "/api/v3/openOrders"                     // 权重(IP): 有 symbol:6,没有 symbol:80,
	FetchSingleOrder        = "/api/v3/order"                          // 权重(IP): 4
	FetchUserTradeUri       = "/api/v3/myTrades"                       // 权重(IP):20
	FetchWithDrawHistoryUri = "/sapi/v1/capital/withdraw/history"      // 权重(IP): 1
	FetchDepositHistoryUri  = "/sapi/v1/capital/deposit/hisrec"        // 权重(IP): 1
	PrivateWithDrawUri      = "/sapi/v1/capital/withdraw/apply"        // 权重(UID): 600
	PrivateTransfer         = "/sapi/v1/asset/transfer"                // 权重(UID): 900
	PrivateTransferWithType = "/sapi/v1/sub-account/universalTransfer" // 权重(IP): 1
	PrivateCurrenciesUri    = "/sapi/v1/capital/config/getall"         // 权重(IP): 10
	PrivateDepositAddr      = "/sapi/v1/capital/deposit/address"       // 权重(IP): 10
)
