package okxv5

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type PositionResponse struct {
	Code    string         `json:"code"`
	Data    []PositionData `json:"data"`
	Message string         `json:"msg"`
}

type PositionData struct {
	Adl         string `json:"adl"`
	AvailPos    string `json:"availPos"`
	AvgPx       string `json:"avgPx"`
	CTime       string `json:"cTime"`
	Ccy         string `json:"ccy"`
	DeltaBS     string `json:"deltaBS"`
	DeltaPA     string `json:"deltaPA"`
	GammaBS     string `json:"gammaBS"`
	GammaPA     string `json:"gammaPA"`
	Imr         string `json:"imr"`
	InstId      string `json:"instId"`
	InstType    string `json:"instType"`
	Interest    string `json:"interest"`
	Last        string `json:"last"`
	Level       string `json:"lever"`
	Liab        string `json:"liab"`
	LiabCcy     string `json:"liabCcy"`
	LiqPx       string `json:"liqPx"`
	MarkPx      string `json:"markPx"`
	Margin      string `json:"margin"`
	MgnMode     string `json:"mgnMode"`
	MgnRatio    string `json:"mgnRatio"`
	Mmr         string `json:"mmr"`
	NotionalUsd string `json:"notionalUsd"`
	OptVal      string `json:"optVal"`
	PTime       string `json:"PTime"`
	Pos         string `json:"pos"`
	PosCcy      string `json:"posCcy"`
	PosId       string `json:"posId"`
	PosSide     string `json:"posSide"`
	ThetaBS     string `json:"thetaBS"`
	ThetaPA     string `json:"thetaPA"`
	TradeId     string `json:"tradeId"`
	UTime       string `json:"UTime"`
	Upl         string `json:"upl"`
	UplRatio    string `json:"uplRatio"`
	VegaBS      string `json:"vegaBS"`
	VegaPA      string `json:"vegaPA"`
}

func (o *PositionResponse) valid() bool {
	return o.Code == "0"
}

func (client *RestClient) FetchPositons() ([]*types.Position, error) {
	queryDict := map[string]interface{}{}
	if client.exchangeType == constant.OkxV5Swap {
		queryDict["instType"] = "SWAP"
	} else {
		queryDict["instType"] = "SPOT"

	}
	payload := utils.UrlEncodeParams(queryDict)
	url := fmt.Sprintf("%s?%s", FetchPositionsUri, payload)
	body, _, err := client.HttpRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("ok get /api/v5/account/positions err:%v", err)
		return nil, err
	}
	response := new(PositionResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/account/positions parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/account/positions fail, code:%s, msg:%s", response.Code, response.Message)
		return nil, err
	}

	// if len(response.Data) == 0 {
	// 	err := fmt.Errorf("ok get /api/v5/account/positions empty")
	// 	return nil, err
	// }

	result, err := positionTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/account/positions transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func positionTransform(response *PositionResponse) ([]*types.Position, error) {
	result := make([]*types.Position, 0, len(response.Data))
	for _, item := range response.Data {
		liquidationPx, err := utils.ParseFloat(item.LiqPx)
		if err != nil {
			liquidationPx = 0
		}
		position, err := utils.ParseFloat(item.Pos)
		if err != nil {
			position = 0
		}
		avgCost, err := utils.ParseFloat(item.AvgPx)
		if err != nil {
			avgCost = 0
		}
		unrealisedPnl, err := utils.ParseFloat(item.Upl)
		if err != nil {
			unrealisedPnl = 0
		}
		last, err := utils.ParseFloat(item.Last)
		if err != nil {
			last = 0
		}
		maintMarginRatio, err := utils.ParseFloat(item.MgnRatio)
		if err != nil {
			maintMarginRatio = 0
		}
		margin, err := utils.ParseFloat(getPositionMargin(item))
		if err != nil {
			margin = 0
		}
		leverage, err := utils.ParseFloat(item.Level)
		if err != nil {
			leverage = 0
		}
		if math.Abs(position) == 0 {
			continue
		}
		info := &types.Position{
			MarginMode:       Okex2MarginMode[item.MgnMode],
			Symbol:           OkInstId2Symbol(item.InstId),
			LiquidationPx:    liquidationPx,
			Position:         math.Abs(position),
			FrozenPosition:   0,
			AvgCost:          avgCost,
			UnrealisedPnl:    unrealisedPnl,
			Last:             last,
			MaintMarginRatio: maintMarginRatio,
			Margin:           margin,
			Leverage:         leverage,
		}
		info.Side = getPositionSide(item.PosSide, position)
		result = append(result, info)
	}
	return result, nil
}

func getPositionMargin(p PositionData) string {
	if p.MgnMode == "isolated" {
		return p.Margin
	}
	return p.Imr
}

func getPositionSide(side string, pos float64) string {
	switch side {
	case "net":
		if pos > 0 {
			return constant.Long.Name()
		}
		return constant.Short.Name()
	case "long":
		return constant.Long.Name()
	default:
		return constant.Short.Name()
	}
}
