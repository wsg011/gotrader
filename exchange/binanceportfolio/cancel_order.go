package binanceportfolio

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wsg011/gotrader/trader/types"
)

func (client *RestClient) CancelUMOrders(orders []*types.Order) ([]*types.OrderResult, error) {
	result := make([]*types.OrderResult, 0)
	for _, order := range orders {
		param := formCancelRequest(order)
		log.Infof("param %v", param)

		uri := CancelUMOrderUri
		body, res, err := client.HttpRequest(http.MethodDelete, uri, param)
		if err != nil {
			log.Errorf("binance DELETE /papi/v1/um/order err: %v", err)
			continue
		}
		if res.StatusCode != 200 {
			log.Errorf("binance DELETE /papi/v1/um/order err: %v %s", res.StatusCode, body)
			continue
		}

		var orderResponse OrderResponse
		if err = json.Unmarshal(body, &orderResponse); err != nil {
			log.Infof("binance DELETE /papi/v1/um/order parsing JSON err: %v", err)
			continue
		}

		info := orderCancelTransform(order.Symbol, &orderResponse)
		result = append(result, info)
	}

	return result, nil
}

func orderCancelTransform(symbol string, info *OrderResponse) *types.OrderResult {
	var result types.OrderResult
	if info.Status != "NEW" {
		result.IsSuccess = true
	}
	result.OrderId = strconv.FormatInt(info.OrderId, 10)
	result.ClientId = info.ClientOrderId
	return &result
}

func formCancelRequest(order *types.Order) map[string]interface{} {
	result := map[string]interface{}{
		"symbol": Symbol2Binance(order.Symbol),
		// "orderId":           order.OrderID,
		"origClientOrderId": order.ClientID,
	}
	return result
}
