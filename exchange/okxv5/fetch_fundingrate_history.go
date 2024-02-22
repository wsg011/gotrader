package okxv5

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/types"
)

type FundingRateHistoryRsp struct {
	BaseOkRsp
	Data []FundingRateHistory `json:"data"`
}

func (t *FundingRateHistoryRsp) valid() bool {
	return t.Code == "0" && len(t.Data) > 0
}

type FundingRateHistory struct {
	Method      string `json:"method"`
	FundingRate string `json:"fundingRate"`
	FundingTime string `json:"fundingTime"`
	InstID      string `json:"instId"`
	InstType    string `json:"instType"`
}

func (client *RestClient) FetchFundingRateHistory(symbol string, limit int64) ([]*types.FundingRate, error) {
	queryDict := map[string]interface{}{}
	queryDict["instId"] = Symbol2OkInstId(symbol)
	queryDict["limit"] = limit
	payload := utils.UrlEncodeParams(queryDict)
	url := RestUrl + fmt.Sprintf("%s?%s", FetchFundingRateHistoryUri, payload)

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("ok get /api/v5/public/funding-rate-history err:%v", err)
		return nil, err
	}

	response := new(FundingRateHistoryRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/public/funding-rate-history parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate-history fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate empty")
		return nil, err
	}

	result, err := fundingRateHistoryTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate transform err:%s", err)
		return nil, err
	}
	return result, nil
}

func fundingRateHistoryTransform(response *FundingRateHistoryRsp) ([]*types.FundingRate, error) {
	result := make([]*types.FundingRate, 0)
	for _, fr := range response.Data {
		rate, err := strconv.ParseFloat(fr.FundingRate, 64)
		if err != nil {
			log.Errorf("parser FundingRateHistoryRsp FundingRate err %s", err)
			continue
		}
		t, err := strconv.ParseInt(fr.FundingTime, 10, 64)
		if err != nil {
			log.Errorf("parser FundingRateHistoryRsp FundingTime err %s", err)
			continue
		}
		fundingRate := &types.FundingRate{
			FundingRate: rate,
			FundingTime: t,
			Method:      fr.Method,
			Symbol:      fr.InstID,
		}
		result = append(result, fundingRate)
	}
	return result, nil
}

// func fundingRateTransform(response *FundingRateRsp) (*types.FundingRate, error) {
// 	fr := response.Data[0]
// 	rate, err := strconv.ParseFloat(fr.FundingRate, 64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nextRate, err := parseStringToFloat(fr.NextFundingRate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	t, err := strconv.ParseInt(fr.FundingTime, 10, 64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nt, err := parseStringToInt(fr.NextFundingTime)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &types.FundingRate{
// 		FundingRate:     rate,
// 		FundingTime:     t,
// 		Method:          fr.Method,
// 		Symbol:          fr.InstID,
// 		NextFundingRate: nextRate,
// 		NextFundingTime: nt,
// 	}, nil
// }
