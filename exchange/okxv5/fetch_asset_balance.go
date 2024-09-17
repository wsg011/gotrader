package okxv5

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wsg011/gotrader/trader/types"
)

type OkxAssetBalance struct {
	AvailBal  string `json:"availBal"`
	Bal       string `json:"bal"`
	Ccy       string `json:"ccy"`
	FrozenBal string `json:"frozenBal"`
}

func (b OkxAssetBalance) ToAssets() types.Asset {
	free, _ := strconv.ParseFloat(b.AvailBal, 64)
	frozen, _ := strconv.ParseFloat(b.FrozenBal, 64)
	total, _ := strconv.ParseFloat(b.Bal, 64)
	return types.Asset{
		Coin:   b.Ccy,
		Free:   free,
		Frozen: frozen,
		Total:  total,
	}
}

type AssetBalanceRsp struct {
	BaseOkRsp
	Data []*OkxAssetBalance `json:"data"`
}

func (client *RestClient) FetchAssetBalance() (*types.Assets, error) {
	url := FetchAssetBalanceUri
	body, _, err := client.HttpRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("ok get /api/v5/asset/balances err:%v", err)
		return nil, err
	}

	response := new(AssetBalanceRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/asset/balances parser err:%v", err)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/asset/balances empty")
		return nil, err
	}

	result, err := assetBalanceTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/account/balance transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func assetBalanceTransform(response *AssetBalanceRsp) (*types.Assets, error) {
	bal := response.Data
	assets := make(map[string]types.Asset, len(bal))
	for _, a := range bal {
		assets[a.Ccy] = a.ToAssets()
	}

	// TODO: 计算totalEq
	totalEq := 0.0
	return &types.Assets{
		Assets:     assets,
		TotalUsdEq: totalEq,
	}, nil
}
