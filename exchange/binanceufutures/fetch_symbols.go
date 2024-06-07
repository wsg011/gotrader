package binanceufutures

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/types"
)

type SymbolsResponse struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol     string                   `json:"symbol"`
	BaseAsset  string                   `json:"baseAsset"`
	QuoteAsset string                   `json:"quoteAsset"`
	Status     string                   `json:"status"`
	Filters    []map[string]interface{} `json:"filters"`
}

func (client *RestClient) FetchSymbols() ([]*types.SymbolInfo, error) {
	url := RestUrl + SymbolsRest

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("binance get /fapi/v1/exchangeInfo err:%v", err)
		return nil, err
	}

	response := new(SymbolsResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("binance get /fapi/v1/exchangeInfo parser err:%v", err)
		return nil, err
	}

	result, err := symbolTransform(response)
	if err != nil {
		err := fmt.Errorf("binance get /fapi/v1/exchangeInfotransform err:%s", err)
		return nil, err
	}

	return result, nil
}

func symbolTransform(response *SymbolsResponse) ([]*types.SymbolInfo, error) {
	result := make([]*types.SymbolInfo, 0, len(response.Symbols))

	for _, item := range response.Symbols {
		if strings.Contains(item.Symbol, "_") {
			continue
		}
		if item.Status != "TRADING" {
			continue
		}

		prec := make(map[string]interface{})
		for _, ft := range item.Filters {
			if ft["filterType"] == "PRICE_FILTER" {
				prec["px_prec"] = utils.DecimalMath(ft["tickSize"].(string))
			} else if ft["filterType"] == "LOT_SIZE" {
				prec["qty_prec"] = utils.DecimalMath(ft["stepSize"].(string))
				min, err := utils.ParseFloat(ft["minQty"].(string))
				if err != nil {
					log.Errorf("spot binance fetch_symbol minQty参数转换失败:%v", err)
					return nil, err
				} else {
					prec["min_cnt"] = min
				}
				max, err := utils.ParseFloat(ft["maxQty"].(string))
				if err != nil {
					log.Errorf("spot binance fetch_symbol maxQty参数转换失败:%v", err)
					return nil, err
				} else {
					prec["max_cnt"] = max
				}

			}
		}

		baseCoin := strings.ToLower(item.BaseAsset)
		quoteCoin := strings.ToLower(item.QuoteAsset)

		info := &types.SymbolInfo{
			Name:       baseCoin + "_" + quoteCoin,
			Base:       baseCoin,
			Quote:      quoteCoin,
			Symbol:     strings.ToUpper(baseCoin) + "_" + strings.ToUpper(quoteCoin),
			FaceVal:    1,
			Multiplier: 1,
			// todo 精确度转换
			PxPrec:  prec["px_prec"].(int32),
			QtyPrec: prec["qty_prec"].(int32),
			MinCnt:  prec["min_cnt"].(float64),
			MaxCnt:  prec["max_cnt"].(float64),
		}

		result = append(result, info)
	}
	return result, nil
}
