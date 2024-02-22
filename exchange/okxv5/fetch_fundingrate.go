package okxv5

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/types"
)

type FundingRateRsp struct {
	BaseOkRsp
	Data []FundingRate `json:"data"`
}

func (t *FundingRateRsp) valid() bool {
	return t.Code == "0" && len(t.Data) > 0
}

type FundingRate struct {
	Method          string `json:"method"`
	FundingRate     string `json:"fundingRate"`
	FundingTime     string `json:"fundingTime"`
	InstID          string `json:"instId"`
	InstType        string `json:"instType"`
	NextFundingRate string `json:"nextFundingRate"`
	NextFundingTime string `json:"nextFundingTime"`
}

func (client *RestClient) FetchFundingRate(symbol string) (*types.FundingRate, error) {
	queryDict := map[string]interface{}{}
	queryDict["instId"] = Symbol2OkInstId(symbol)
	payload := utils.UrlEncodeParams(queryDict)
	url := RestUrl + fmt.Sprintf("%s?%s", FetchFundingRateUri, payload)

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("ok get /api/v5/public/funding-rate err:%v", err)
		return nil, err
	}

	response := new(FundingRateRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/public/funding-rate parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate empty")
		return nil, err
	}

	result, err := fundingRateTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/public/funding-rate transform err:%s", err)
		return nil, err
	}
	return result, nil
}

func fundingRateTransform(response *FundingRateRsp) (*types.FundingRate, error) {
	fr := response.Data[0]
	rate, err := strconv.ParseFloat(fr.FundingRate, 64)
	if err != nil {
		return nil, err
	}
	nextRate, err := parseStringToFloat(fr.NextFundingRate)
	if err != nil {
		return nil, err
	}
	t, err := strconv.ParseInt(fr.FundingTime, 10, 64)
	if err != nil {
		return nil, err
	}
	nt, err := parseStringToInt(fr.NextFundingTime)
	if err != nil {
		return nil, err
	}

	return &types.FundingRate{
		FundingRate:     rate,
		FundingTime:     t,
		Method:          fr.Method,
		Symbol:          fr.InstID,
		NextFundingRate: nextRate,
		NextFundingTime: nt,
	}, nil
}
