package okxv5

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wsg011/gotrader/pkg/httpx"
	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

var httpClient = httpx.NewClient()

type RestClient struct {
	apiKey       string
	secretKey    string
	passPhrase   string
	exchangeType constant.ExchangeType
}

type BaseOkRsp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewRestClient(apiKey, secretKey, passPhrase string, exchangeType constant.ExchangeType) *RestClient {
	client := &RestClient{
		apiKey:       apiKey,
		secretKey:    secretKey,
		passPhrase:   passPhrase,
		exchangeType: exchangeType,
	}
	return client
}

func (client *RestClient) HttpRequest(method string, uri string, payload []byte) ([]byte, *http.Response, error) {
	var param string
	if payload != nil {
		param = string(payload)
	}
	currentTime := IsoTime()
	toSignStr := currentTime + method + uri + param
	signature := utils.GenBase64Digest(utils.HmacSha256(toSignStr, client.secretKey))
	url := RestUrl + uri
	head := map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-KEY":        client.apiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  currentTime,
		"OK-ACCESS-PASSPHRASE": client.passPhrase,
	}
	args := &httpx.Request{
		Url:    url,
		Head:   head,
		Method: method,
		Body:   payload,
	}
	body, res, err := httpClient.Request(args)
	if err != nil {
		return nil, res, err
	}
	return *body, res, err
}

func (client *RestClient) HttpGet(url string) ([]byte, *http.Response, error) {
	body, res, err := httpClient.Get(url)
	if err != nil {
		return nil, res, err
	}
	return *body, res, nil
}

type KlineRsp struct {
	BaseOkRsp
	Data [][]string `json:"data"`
}

func (t *KlineRsp) valid() bool {
	return t.Code == "0" && len(t.Data) > 0
}

func (client *RestClient) FetchKline(symbol string, interval string, limit int64) ([]types.Kline, error) {
	queryDict := map[string]interface{}{}
	queryDict["instId"] = Symbol2OkInstId(symbol)
	queryDict["bar"] = interval
	queryDict["limit"] = limit

	payload := utils.UrlEncodeParams(queryDict)
	url := RestUrl + fmt.Sprintf(FetchKlineUri, payload)

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("ok get /api/v5/market/candles err:%v", err)
		return nil, err
	}

	response := new(KlineRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/market/candles parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/market/candles fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/market/candles empty")
		return nil, err
	}

	result, err := klineTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/market/candles transform err:%s", err)
		return nil, err
	}
	return result, nil
}

func klineTransform(response *KlineRsp) ([]types.Kline, error) {
	result := make([]types.Kline, 0, len(response.Data))
	for _, dat := range response.Data {
		if len(dat) < 9 {
			log.Errorf("data len less 9 %v", len(dat))
			continue
		}
		ts, _ := strconv.ParseInt(dat[0], 10, 64)
		open, err := strconv.ParseFloat(dat[1], 64)
		if err != nil {
			log.Errorf("parser open err %s", err)
			continue
		}
		high, err := strconv.ParseFloat(dat[2], 64)
		if err != nil {
			log.Errorf("parser high err %s", err)
			continue
		}
		low, err := strconv.ParseFloat(dat[3], 64)
		if err != nil {
			log.Errorf("parser low err %s", err)
			continue
		}
		close, err := strconv.ParseFloat(dat[4], 64)
		if err != nil {
			log.Errorf("parser close err %s", err)
			continue
		}
		vol, err := strconv.ParseFloat(dat[5], 64)
		if err != nil {
			log.Errorf("parser vol err %s", err)
			continue
		}
		confirm, err := strconv.ParseInt(dat[8], 10, 64)
		if err != nil {
			log.Errorf("parser confirm err %s", err)
			continue
		}

		k := types.Kline{
			Ts:      ts,
			Open:    open,
			High:    high,
			Low:     low,
			Close:   close,
			Vol:     vol,
			Confirm: confirm,
		}
		result = append(result, k)
	}
	return result, nil
}

type OkBalance struct {
	AdjEq       string             `json:"adjEq"`
	BorrowFroz  string             `json:"borrowFroz"`
	Details     []*OkBalanceDetail `json:"details"`
	Imr         string             `json:"imr"`
	IsoEq       string             `json:"isoEq"`
	MgnRatio    string             `json:"mgnRatio"`
	Mmr         string             `json:"mmr"`
	NotionalUsd string             `json:"notionalUsd"`
	OrdFroz     string             `json:"ordFroz"`
	TotalEq     string             `json:"totalEq"`
	UTime       string             `json:"uTime"`
}

// 去掉了一些不太可能使用的字段
type OkBalanceDetail struct {
	AvailBal      string `json:"availBal"`
	AvailEq       string `json:"availEq"`
	CashBal       string `json:"cashBal"`
	Ccy           string `json:"ccy"`
	CrossLiab     string `json:"crossLiab"`
	Eq            string `json:"eq"`
	EqUsd         string `json:"eqUsd"`
	FixedBal      string `json:"fixedBal"`
	FrozenBal     string `json:"frozenBal"`
	IsoEq         string `json:"isoEq"`
	IsoLiab       string `json:"isoLiab"`
	IsoUpl        string `json:"isoUpl"`
	MgnRatio      string `json:"mgnRatio"`
	NotionalLever string `json:"notionalLever"`
	OrdFrozen     string `json:"ordFrozen"`
	UTime         string `json:"uTime"`
	Upl           string `json:"upl"`
	UplLiab       string `json:"uplLiab"`
}

func (b OkBalanceDetail) ToAssets() types.Asset {
	free, _ := strconv.ParseFloat(b.CashBal, 64)
	frozen, _ := strconv.ParseFloat(b.FrozenBal, 64)
	total, _ := strconv.ParseFloat(b.Eq, 64)
	return types.Asset{
		Coin:   b.Ccy,
		Free:   free,
		Frozen: frozen,
		Total:  total,
	}
}

type BalanceRsp struct {
	BaseOkRsp
	Data []*OkBalance `json:"data"`
}

func (t *BalanceRsp) valid() bool {
	return t.Code == "0" && len(t.Data) > 0
}

func (client *RestClient) FetchBalance() (*types.Assets, error) {
	url := FetchBalanceUri
	body, _, err := client.HttpRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("ok get /api/v5/account/balance err:%v", err)
		return nil, err
	}

	response := new(BalanceRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/account/balance parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/account/balance fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/account/balance empty")
		return nil, err
	}

	result, err := balanceTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/account/balance transform err:%s", err)
		return nil, err
	}
	return result, nil
}

func balanceTransform(response *BalanceRsp) (*types.Assets, error) {
	bal := response.Data[0]
	assets := make(map[string]types.Asset, len(bal.Details))
	for _, a := range bal.Details {
		assets[a.Ccy] = a.ToAssets()
	}
	totalEq, _ := strconv.ParseFloat(bal.TotalEq, 64)
	return &types.Assets{
		Assets:     assets,
		TotalUsdEq: totalEq,
	}, nil
}

type CreateOrderResult struct {
	ClOrdID string `json:"clOrdId"`
	OrdID   string `json:"ordId"`
	Tag     string `json:"tag"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

type CreateOrderResponse struct {
	BaseOkRsp
	Data []*CreateOrderResult `json:"data"`
}

func (client *RestClient) CreateBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	startTime := time.Now()

	param := make([]map[string]interface{}, 0)
	for _, item := range orders {
		param = append(param, formRequest(item))
	}
	payload, _ := json.Marshal(param)
	uri := CreateBatchOrderUri
	body, _, err := client.HttpRequest(http.MethodPost, uri, payload)
	httpRequestTime := time.Now()

	if err != nil {
		log.Errorf("okx post /api/v5/trade/batch-orders err: %v", err)
		return nil, err
	}
	response := new(CreateOrderResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("okx post /api/v5/trade/batch-orders err 数据解析失败:%v", err)
		return nil, err
	}
	if len(response.Data) == 0 {
		err := fmt.Errorf("ok post /api/v5/trade/batch-orders err: %v", response)
		return nil, err
	}

	result := make([]*types.OrderResult, 0)
	var symbol string
	if len(orders) > 0 {
		symbol = orders[0].Symbol
	}
	for _, item := range response.Data {
		info := orderTransform(symbol, item)
		result = append(result, info)
	}
	log.Infof("HTTP cost time: %v", httpRequestTime.Sub(startTime))
	return result, nil
}

func formRequest(order *types.Order) map[string]interface{} {
	oSide := OkxOrderSide[order.Side.Name()]
	oType := OkxOrderType[order.Type.Name()]
	// tdModel := "cash" // 现货
	// tmp := strings.Split(order.Symbol, "_")
	// if len(tmp) == 3 {
	// 	tdModel = "cross"
	// }

	// 目前只支持全仓模式
	tdModel := "cross"
	result := map[string]interface{}{
		"instId":  Symbol2OkInstId(order.Symbol),
		"tdMode":  tdModel,
		"side":    Side2Okx[oSide],
		"ordType": Type2Okx[oType],
		"px":      order.Price,
		"sz":      order.OrigQty,
	}
	if order.ClientID != "" {
		result["clOrdId"] = order.ClientID
	}

	return result
}

func orderTransform(symbol string, info *CreateOrderResult) *types.OrderResult {
	var result types.OrderResult
	if info.SCode == "0" {
		result.IsSuccess = true
	}
	result.OrderId = info.OrdID
	result.ClientId = info.ClOrdID
	result.ErrMsg = info.SMsg
	return &result
}

type CancelOrderResponse struct {
	Code string               `json:"code"`
	Msg  string               `json:"msg"`
	Data []*CancelOrderResult `json:"data"`
}

type CancelOrderResult struct {
	ClOrdID string `json:"clOrdId"`
	OrdID   string `json:"ordId"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

func (client *RestClient) CancelBatchOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	startTime := time.Now()
	param := make([]map[string]interface{}, 0)
	for _, item := range orders {
		param = append(param, formCancelRequest(item))
	}
	payload, _ := json.Marshal(param)
	uri := CancelBatchOrderUri
	body, _, err := client.HttpRequest(http.MethodPost, uri, payload)
	httpRequestTime := time.Now()
	if err != nil {
		log.Errorf("okx post /api/v5/trade/cancel-batch-orders err: %v", err)
		return nil, err
	}
	response := new(CancelOrderResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("okx post /api/v5/trade/cancel-batch-orders err 数据解析失败:%v", err)
		return nil, err
	}
	if len(response.Data) == 0 {
		err := fmt.Errorf("ok post /api/v5/trade/cancel-batch-orders err: %v", response)
		return nil, err
	}

	result := make([]*types.OrderResult, 0)
	var symbol string
	if len(orders) > 0 {
		symbol = orders[0].Symbol
	}
	for _, item := range response.Data {
		info := cancelOrderTransform(symbol, item)
		result = append(result, info)
	}

	log.Infof("HTTP cost time: %v", httpRequestTime.Sub(startTime))
	return result, nil
}

func formCancelRequest(order *types.Order) map[string]interface{} {
	result := map[string]interface{}{
		"instId":  Symbol2OkInstId(order.Symbol),
		"ordId":   order.OrderID,
		"clOrdId": order.ClientID,
	}
	return result
}

func cancelOrderTransform(symbol string, info *CancelOrderResult) *types.OrderResult {
	var result types.OrderResult
	if info.SCode == "0" {
		result.IsSuccess = true
	}
	result.OrderId = info.OrdID
	result.ClientId = info.ClOrdID
	result.ErrMsg = info.SMsg
	return &result
}
