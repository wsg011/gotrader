package binanceportfolio

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type PositionInfo struct {
	Symbol           string  `json:"symbol"`
	PositionAmt      float64 `json:"positionAmt,string"`
	EntryPrice       float64 `json:"entryPrice,string"` // Use string option for parsing floats encoded as strings
	MarkPrice        float64 `json:"markPrice,string"`
	UnRealizedProfit float64 `json:"unRealizedProfit,string"`
	LiquidationPrice float64 `json:"liquidationPrice,string"`
	Leverage         float64 `json:"leverage,string"`
	MaxNotionalValue string  `json:"maxNotionalValue"`
	PositionSide     string  `json:"positionSide"`
	Notional         float64 `json:"notional,string"` // Convert string to float64 as well
	UpdateTime       int64   `json:"updateTime"`
}

func (client *RestClient) FetchPositons() ([]*types.Position, error) {
	uri := FetchPositionsUri
	body, res, err := client.HttpRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Errorf("binance FetchPositons 网络错误:%v", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Errorf("binance FetchPositons %v %s", res.StatusCode, body)
		return nil, fmt.Errorf("%s", body)
	}

	var positions []PositionInfo
	if err = json.Unmarshal([]byte(body), &positions); err != nil {
		log.Fatalf("binance FetchPositons Error parsing JSON: %v", err)
	}
	result, err := positionTransform(positions)
	if err != nil {
		err := fmt.Errorf("binance FetchPositons transform err:%s", err)
		return nil, err
	}
	return result, nil
}

func positionTransform(response []PositionInfo) ([]*types.Position, error) {
	result := make([]*types.Position, 0, len(response))
	for _, item := range response {
		log.Infof("item %+v", item)
		info := &types.Position{
			Symbol:        strings.Replace(item.Symbol, "USDT", "_USDT", 1),
			LiquidationPx: item.LiquidationPrice,
			Position:      math.Abs(item.PositionAmt),
			AvgCost:       item.EntryPrice,
			UnrealisedPnl: item.UnRealizedProfit,
			Leverage:      item.Leverage,
		}
		info.Side = getPositionSide(item.PositionSide, item.PositionAmt)
		result = append(result, info)

	}
	return result, nil
}

func getPositionSide(side string, pos float64) string {
	switch side {
	case "BOTH":
		if pos > 0 {
			return constant.Long.Name()
		}
		return constant.Short.Name()
	case "LONG":
		return constant.Long.Name()
	default:
		return constant.Short.Name()
	}
}
