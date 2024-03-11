package okxv5

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

type TickerRsp struct {
	BaseOkRsp
	Data []Ticker `json:"data"`
}

func (t *TickerRsp) valid() bool {
	return t.Code == "0" && len(t.Data) > 0
}

type Ticker struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24H   string `json:"open24h"`
	High24H   string `json:"high24h"`
	Low24H    string `json:"low24h"`
	VolCcy24H string `json:"volCcy24h"`
	Vol24H    string `json:"vol24h"`
	Ts        string `json:"ts"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
}

func (client *RestClient) FetchTickers() ([]*types.Ticker, error) {
	queryDict := map[string]interface{}{}
	queryDict["instType"] = "SPOT"
	if client.exchangeType == constant.OkxV5Swap {
		queryDict["instType"] = "SWAP"
	}
	payload := utils.UrlEncodeParams(queryDict)
	url := RestUrl + fmt.Sprintf(TickersRest, payload)

	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("ok get /api/v5/market/tickers err:%v", err)
		return nil, err
	}

	response := new(TickerRsp)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("ok get /api/v5/market/tickers parser err:%v", err)
		return nil, err
	}

	if !response.valid() {
		err := fmt.Errorf("ok get /api/v5/market/tickers fail, code:%s, msg:%s", response.Code, response.Msg)
		return nil, err
	}

	if len(response.Data) == 0 {
		err := fmt.Errorf("ok get /api/v5/market/tickers empty")
		return nil, err
	}

	result, err := tickersTransform(response)
	if err != nil {
		err := fmt.Errorf("ok get /api/v5/market/tickers transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func tickersTransform(response *TickerRsp) ([]*types.Ticker, error) {
	result := make([]*types.Ticker, 0, len(response.Data))
	for _, tickerData := range response.Data {
		symbol := OkInstId2Symbol(tickerData.InstID)
		open24H, _ := strconv.ParseFloat(tickerData.Open24H, 64)
		high24H, _ := strconv.ParseFloat(tickerData.High24H, 64)
		low24H, _ := strconv.ParseFloat(tickerData.Low24H, 64)
		lastPrice, _ := strconv.ParseFloat(tickerData.Last, 64)
		lastSize, _ := strconv.ParseFloat(tickerData.LastSz, 64)
		vol, _ := strconv.ParseFloat(tickerData.Vol24H, 64)
		askPrice, _ := strconv.ParseFloat(tickerData.AskPx, 64)
		bidPrice, _ := strconv.ParseFloat(tickerData.BidPx, 64)
		askSize, _ := strconv.ParseFloat(tickerData.AskSz, 64)
		bidSize, _ := strconv.ParseFloat(tickerData.BidSz, 64)
		exchangeTs, _ := strconv.ParseInt(tickerData.Ts, 10, 64)

		ticker := &types.Ticker{
			Symbol:     symbol,
			Open:       open24H,
			High:       high24H,
			Low:        low24H,
			Vol:        vol,
			LastPrice:  lastPrice,
			LastSize:   lastSize,
			AskPrice:   askPrice,
			BidPrice:   bidPrice,
			AskSize:    askSize,
			BidSize:    bidSize,
			ExchangeTs: exchangeTs,
			Ts:         utils.Millisec(time.Now()),
		}
		result = append(result, ticker)
	}
	return result, nil
}
