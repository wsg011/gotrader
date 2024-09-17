package okxv5

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wsg011/gotrader/trader/types"
)

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
