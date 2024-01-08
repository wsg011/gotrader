package okxv5swap

import (
	"encoding/json"
	"fmt"
	"gotrader/pkg/utils"
	"gotrader/trader/constant"
	"gotrader/trader/types"
	"strconv"
	"time"
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

func (client *RestClient) FetchTickers() ([]*types.BookTickerEvent, error) {
	queryDict := map[string]interface{}{}
	queryDict["instType"] = "SWAP"
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

func tickersTransform(response *TickerRsp) ([]*types.BookTickerEvent, error) {
	result := make([]*types.BookTickerEvent, 0, len(response.Data))
	for _, ticker := range response.Data {
		symbol := OkInstId2Symbol(ticker.InstID)
		ex := constant.OkxV5Future
		askPrice, _ := strconv.ParseFloat(ticker.AskPx, 64)
		bidPrice, _ := strconv.ParseFloat(ticker.BidPx, 64)
		askSize, _ := strconv.ParseFloat(ticker.AskSz, 64)
		bidSize, _ := strconv.ParseFloat(ticker.BidSz, 64)
		exchangeTs, _ := strconv.ParseInt(ticker.Ts, 10, 64)
		bookTicker := &types.BookTickerEvent{
			Symbol:     symbol,
			Exchange:   ex,
			AskPrice:   askPrice,
			BidPrice:   bidPrice,
			AskQty:     askSize,
			BidQty:     bidSize,
			ExchangeTs: exchangeTs,
			Ts:         utils.Millisec(time.Now()),
		}
		result = append(result, bookTicker)
	}
	return result, nil
}
