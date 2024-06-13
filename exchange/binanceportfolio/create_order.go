package binanceportfolio

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type OrderResponse struct {
	ClientOrderId           string `json:"clientOrderId"`
	CumQty                  string `json:"cumQty"`
	CumQuote                string `json:"cumQuote"`
	ExecutedQty             string `json:"executedQty"`
	OrderId                 int64  `json:"orderId"`
	AvgPrice                string `json:"avgPrice"`
	OrigQty                 string `json:"origQty"`
	Price                   string `json:"price"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	Status                  string `json:"status"`
	Symbol                  string `json:"symbol"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int64  `json:"goodTillDate"`
	UpdateTime              int64  `json:"updateTime"`
}

func (client *RestClient) CreateUMOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	result := make([]*types.OrderResult, 0)
	for _, order := range orders {
		param := formRequest(order)

		uri := CreateOrderUri
		body, res, err := client.HttpRequest(http.MethodPost, uri, param)
		if err != nil {
			log.Errorf("binance post /papi/v1/um/order err: %v", err)
			continue
		}
		if res.StatusCode != 200 {
			log.Errorf("binance post /papi/v1/um/order err: %v %s", res.StatusCode, body)
			continue
		}

		var orderResponse OrderResponse
		if err = json.Unmarshal(body, &orderResponse); err != nil {
			log.Infof("binance post /papi/v1/um/order parsing JSON err: %v", err)
			continue
		}

		info := orderTransform(order.Symbol, &orderResponse)
		result = append(result, info)
	}

	return result, nil
}

func orderTransform(symbol string, info *OrderResponse) *types.OrderResult {
	var result types.OrderResult
	if info.Status == "NEW" || info.Status == "FILLED" {
		result.IsSuccess = true
	}
	result.OrderId = strconv.FormatInt(info.OrderId, 10)
	result.ClientId = info.ClientOrderId
	return &result
}

func formRequest(order *types.Order) map[string]interface{} {
	oSide := BinanceOrderSide[order.Side.Name()]
	oType := BinanceOrderType[order.Type.Name()]

	// 目前只支持全仓模式
	result := map[string]interface{}{
		"symbol":   Symbol2Binance(order.Symbol),
		"side":     Side2Binance[oSide],
		"type":     Type2Binance[oType],
		"quantity": order.OrigQty,
	}
	if order.Type == constant.Limit {
		result["price"] = order.Price
		result["timeInForce"] = "GTC"
	}
	if order.ClientID != "" {
		result["newClientOrderId"] = order.ClientID
	}

	return result
}

func (client *RestClient) CreateMMOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	result := make([]*types.OrderResult, 0)
	for _, order := range orders {
		param := formRequest(order)

		uri := CreateMMOrderUri
		body, res, err := client.HttpRequest(http.MethodPost, uri, param)
		if err != nil {
			log.Errorf("binance post /papi/v1/margin/order err: %v", err)
			continue
		}
		if res.StatusCode != 200 {
			log.Errorf("binance post /papi/v1/margin/order err: %v %s", res.StatusCode, body)
			continue
		}

		var orderResponse OrderResponse
		if err = json.Unmarshal(body, &orderResponse); err != nil {
			log.Infof("binance post /papi/v1/margin/order parsing JSON err: %v", err)
			continue
		}

		info := orderTransform(order.Symbol, &orderResponse)
		result = append(result, info)
	}

	return result, nil
}
