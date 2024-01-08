package constant

type RestApi string

var (
	API_DEBUG1               RestApi = "apiDebug1"
	API_DEBUG2               RestApi = "apiDebug2"
	API_FETCH_BALANCE        RestApi = "fetchBalance"
	API_FETCH_POSITION       RestApi = "fetchPosition"
	API_GET_TICKER           RestApi = "getTicker"
	API_GET_TICKERS          RestApi = "getTickers"
	API_GET_FUNDINGRATE      RestApi = "getFundingRate"
	API_GET_KLINE            RestApi = "getKline"
	API_PLACE_SPOT_ORDER     RestApi = "placeSpotOrder"
	API_PLACE_FUTURE_ORDER   RestApi = "placeFutureOrder"
	API_RPC_PLACE_SPOT_ORDER RestApi = "placeRpcSpotOrder"
	API_QUERY_ORDER          RestApi = "placeQueryOrder"
)
