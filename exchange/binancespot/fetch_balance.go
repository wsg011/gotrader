package binancespot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/types"
)

type BalanceResponse struct {
	Balances []BalanceInfo `json:"balances"`
}

type BalanceInfo struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

func (client *RestClient) FetchBalance() (*types.Assets, error) {
	uri := FetchBalanceUri
	param := make(map[string]interface{})
	body, _, err := client.HttpRequest(http.MethodGet, uri, param)
	if err != nil {
		log.Errorf("binance FetchBalance 网络错误:%v", err)
		return nil, err
	}

	response := new(BalanceResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("binance get /api/v3/exchangeInfo parser err:%v", err)
		return nil, err
	}

	result, err := balanceTransform(response)
	if err != nil {
		err := fmt.Errorf("binance get /api/v3/exchangeInfo transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func balanceTransform(response *BalanceResponse) (*types.Assets, error) {
	result := &types.Assets{}

	assets := make(map[string]types.Asset)
	for _, item := range response.Balances {
		free, err := utils.ParseFloat(item.Free)
		if err != nil {
			log.Errorf("spot binance fetchBalance free参数转换失败:%v", err)
			return nil, err
		}
		locked, err := utils.ParseFloat(item.Locked)
		if err != nil {
			log.Errorf("spot binance fetchBalance locked参数转换失败:%v", err)
			return nil, err
		}
		total := free + locked
		if decimal.NewFromFloat(total).Equal(decimal.NewFromInt(0)) {
			continue
		}

		coin := strings.ToUpper(item.Asset)
		info := types.Asset{
			Coin:   coin,
			Total:  total,
			Frozen: locked,
			Free:   free,
		}
		assets[coin] = info
	}
	result.Assets = assets
	result.TotalUsdEq = 0
	return result, nil
}

// func balanceTransform(response *BalanceRsp) (*types.Assets, error) {
// 	bal := response.Data[0]
// 	assets := make(map[string]types.Asset, len(bal.Details))
// 	for _, a := range bal.Details {
// 		assets[a.Ccy] = a.ToAssets()
// 	}
// 	totalEq, _ := strconv.ParseFloat(bal.TotalEq, 64)
// 	return &types.Assets{
// 		Assets:     assets,
// 		TotalUsdEq: totalEq,
// 	}, nil
// }
