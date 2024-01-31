package okxv5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type SymbolsResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []Symbol `json:"data"`
}
type Symbol struct {
	Alias        string `json:"alias"`
	BaseCcy      string `json:"baseCcy"`
	Category     string `json:"category"`
	CtMult       string `json:"ctMult"`
	CtType       string `json:"ctType"`
	CtVal        string `json:"ctVal"`
	CtValCcy     string `json:"ctValCcy"`
	ExpTime      string `json:"expTime"`
	InstID       string `json:"instId"`
	InstType     string `json:"instType"`
	Lever        string `json:"lever"`
	ListTime     string `json:"listTime"`
	LotSz        string `json:"lotSz"`
	MaxIcebergSz string `json:"maxIcebergSz"`
	MaxLmtSz     string `json:"maxLmtSz"`
	MaxMktSz     string `json:"maxMktSz"`
	MaxStopSz    string `json:"maxStopSz"`
	MaxTriggerSz string `json:"maxTriggerSz"`
	MaxTwapSz    string `json:"maxTwapSz"`
	MinSz        string `json:"minSz"`
	OptType      string `json:"optType"`
	QuoteCcy     string `json:"quoteCcy"`
	SettleCcy    string `json:"settleCcy"`
	State        string `json:"state"`
	Stk          string `json:"stk"`
	TickSz       string `json:"tickSz"`
	Uly          string `json:"uly"`
}

func (o *SymbolsResponse) valid() bool {
	return o.Code == "0" && len(o.Data) > 0
}

func (client *RestClient) FetchSymbols() ([]*types.SymbolInfo, error) {
	queryDict := map[string]interface{}{}
	if client.exchangeType == constant.OkxV5Swap {
		queryDict["instType"] = "SWAP"
	} else {
		queryDict["instType"] = "SPOT"

	}
	payload := utils.UrlEncodeParams(queryDict)
	url := RestUrl + fmt.Sprintf("%s?%s", SymbolsRest, payload)

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("ok get /api/v5/public/instruments err:%v", err)
		return nil, err
	}

	response := new(SymbolsResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/public/instruments parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/public/instruments fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/public/instruments empty")
		return nil, err
	}

	result, err := symbolTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func symbolTransform(response *SymbolsResponse) ([]*types.SymbolInfo, error) {
	result := make([]*types.SymbolInfo, 0, len(response.Data))
	for _, item := range response.Data {
		if item.State != "live" {
			continue
		}
		minCnt, err := utils.ParseFloat(item.MinSz)
		if err != nil {
			log.Errorf("okx fetch_symbol minCnt参数转换失败:%v", err)
			return nil, err
		}
		maxCnt, err := utils.ParseFloat(item.MaxLmtSz)
		if err != nil {
			log.Errorf("okx fetch_symbol maxCnt参数转换失败:%v", err)
			return nil, err
		}
		FaceVal, err := utils.ParseFloat(item.CtVal)
		if err != nil {
			log.Errorf("okx fetch_symbols FaceVal参数转换失败:%v", err)
			return nil, err
		}
		Multiplier, err := utils.ParseInt(item.CtMult)
		if err != nil {
			log.Errorf("okx fetch_symbols Multiplier参数转换失败:%v", err)
			return nil, err
		}
		baseCoin := strings.ToLower(item.BaseCcy)
		QuoteCoin := strings.ToLower(item.QuoteCcy)
		info := &types.SymbolInfo{
			Base:       baseCoin,
			Quote:      QuoteCoin,
			Symbol:     OkInstId2Symbol(item.InstID),
			FaceVal:    FaceVal,
			Multiplier: Multiplier,
			PxPrec:     int32(strings.Count(item.TickSz, "0")),
			QtyPrec:    int32(strings.Count(item.LotSz, "0")),
			MinCnt:     minCnt,
			MaxCnt:     maxCnt,
			Name:       baseCoin + "_" + QuoteCoin,
		}

		result = append(result, info)
	}
	return result, nil
}
