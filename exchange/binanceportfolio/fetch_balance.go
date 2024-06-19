package binanceportfolio

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
	Asset               string `json:"asset"`
	TotalWalletBalance  string `json:"totalWalletBalance"`
	CrossMarginAsset    string `json:"crossMarginAsset"`
	CrossMarginBorrowed string `json:"crossMarginBorrowed"`
	CrossMarginFree     string `json:"crossMarginFree"`
	CrossMarginInterest string `json:"crossMarginInterest"`
	CrossMarginLocked   string `json:"crossMarginLocked"`
	UmWalletBalance     string `json:"umWalletBalance"`
	UmUnrealizedPNL     string `json:"umUnrealizedPNL"`
	CmWalletBalance     string `json:"cmWalletBalance"`
	CmUnrealizedPNL     string `json:"cmUnrealizedPNL"`
	UpdateTime          int64  `json:"updateTime"`
	NegativeBalance     string `json:"negativeBalance"`
}

func (client *RestClient) FetchBalance() (*types.Assets, error) {
	uri := FetchBalanceUri
	param := make(map[string]interface{})
	body, _, err := client.HttpRequest(http.MethodGet, uri, param)
	if err != nil {
		log.Errorf("binance FetchBalance 网络错误:%v", err)
		return nil, err
	}

	/**
		[
	    {
	        "asset": "USDT",    // 资产
	        "totalWalletBalance": "122607.35137903", // 钱包余额 =  全仓杠杆未锁定 + 全仓杠杆锁定 + u本位合约钱包余额 + 币本位合约钱包余额
	        "crossMarginAsset": "92.27530794", // 全仓资产 = 全仓杠杆未锁定 + 全仓杠杆锁定
	        "crossMarginBorrowed": "10.00000000", // 全仓杠杆借贷
	        "crossMarginFree": "100.00000000", // 全仓杠杆未锁定
	        "crossMarginInterest": "0.72469206", // 全仓杠杆利息
	        "crossMarginLocked": "3.00000000", //全仓杠杆锁定
	        "umWalletBalance": "0.00000000",  // u本位合约钱包余额
	        "umUnrealizedPNL": "23.72469206",     // u本位未实现盈亏
	        "cmWalletBalance": "23.72469206",       // 币本位合约钱包余额
	        "cmUnrealizedPNL": "",    // 币本位未实现盈亏
	        "updateTime": 1617939110373
	    }
	]
	**/
	var balances []BalanceInfo
	if err = json.Unmarshal(body, &balances); err != nil {
		log.Errorf("binance get /papi/v1/balance parser err:%v", err)
		return nil, err
	}

	log.Infof("balance %s", body)

	result, err := balanceTransform(balances)
	if err != nil {
		err := fmt.Errorf("binance get /papi/v1/balance transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func balanceTransform(response []BalanceInfo) (*types.Assets, error) {
	result := &types.Assets{}

	totalUsdEq := 0.0

	assets := make(map[string]types.Asset)
	for _, item := range response {
		free, err := utils.ParseFloat(item.TotalWalletBalance)
		if err != nil {
			log.Errorf("spot binance fetchBalance TotalWalletBalance 参数转换失败:%v", err)
			return nil, err
		}
		locked, err := utils.ParseFloat(item.CrossMarginLocked)
		if err != nil {
			log.Errorf("spot binance fetchBalance CrossMarginLocked 参数转换失败:%v", err)
			return nil, err
		}
		umUnrealizedPNL, err := utils.ParseFloat(item.UmUnrealizedPNL)
		if err != nil {
			log.Errorf("portfolio binance fetchBalance umUnrealizedPNL 参数转换失败:%v", err)
			return nil, err
		}
		total := free + locked
		if decimal.NewFromFloat(total).Equal(decimal.NewFromInt(0)) {
			continue
		}

		coin := strings.ToUpper(item.Asset)
		if coin == "USDT" {
			total = total + umUnrealizedPNL
			totalUsdEq = total
		}
		info := types.Asset{
			Coin:   coin,
			Total:  total,
			Frozen: locked,
			Free:   free,
		}
		assets[coin] = info

	}
	result.Assets = assets
	result.TotalUsdEq = totalUsdEq
	return result, nil
}
