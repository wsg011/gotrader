package binanceufutures

import (
	"fmt"
	"strings"
	"time"

	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/types"
)

type TickerResponse []TickerInfo

type TickerInfo struct {
	Symbol      string `json:"symbol"`
	HighPrice   string `json:"highPrice"`
	LowPrice    string `json:"lowPrice"`
	OpenPrice   string `json:"openPrice"`
	LastPrice   string `json:"lastPrice"`
	LastSize    string `json:"lastSize"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quoteVolume"`
	CloseTime   int64  `json:"closeTime"`
}

func (t *TickerInfo) valid() bool {
	return t.Symbol != ""
}

func (client *RestClient) FetchTickers() ([]*types.Ticker, error) {
	url := RestUrl + TickerRest
	body, _, err := client.HttpGet(url)
	if err != nil {
		log.Errorf("futures binance FetchTicker 网络错误:%v", err)
		return nil, err
	}

	response := make(TickerResponse, 0)
	if err := utils.JsonDecode(body, &response); err != nil {
		return nil, err
	}

	result, err := tickersTransform(&response)
	if err != nil {
		err := fmt.Errorf("binance get /fapi/v1/exchangeInfo transform err:%s", err)
		return nil, err
	}

	return result, nil
}

func tickersTransform(response *TickerResponse) ([]*types.Ticker, error) {
	result := make([]*types.Ticker, 0, len(*response))
	for _, item := range *response {
		if !item.valid() {
			continue
		}
		info, err := trans(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

func trans(item *TickerInfo) (*types.Ticker, error) {
	open, err := utils.ParseFloat(item.OpenPrice)
	if err != nil {
		log.Errorf("futures binance fetch_ticker open参数转换失败:%v", err)
		return nil, err
	}
	high, err := utils.ParseFloat(item.HighPrice)
	if err != nil {
		log.Errorf("futures binance fetch_ticker high参数转换失败:%v", err)
		return nil, err
	}
	low, err := utils.ParseFloat(item.LowPrice)
	if err != nil {
		log.Errorf("futures binance fetch_ticker low参数转换失败:%v", err)
		return nil, err
	}
	lastPrice, err := utils.ParseFloat(item.LastPrice)
	if err != nil {
		log.Errorf("futures binance fetch_ticker close参数转换失败:%v", err)
		return nil, err
	}
	lastSize, err := utils.ParseFloat(item.LastSize)
	if err != nil {
		log.Errorf("futures binance fetch_ticker close参数转换失败:%v", err)
		return nil, err
	}
	baseVolume, err := utils.ParseFloat(item.Volume)
	if err != nil {
		log.Errorf("futures binance fetch_ticker baseVolume参数转换失败:%v", err)
		return nil, err
	}

	info := &types.Ticker{
		Symbol:     strings.Replace(item.Symbol, "USDT", "_USDT", 1),
		Open:       open,
		High:       high,
		Low:        low,
		Vol:        baseVolume,
		LastPrice:  lastPrice,
		LastSize:   lastSize,
		Ts:         utils.Millisec(time.Now()),
		ExchangeTs: item.CloseTime,
	}
	return info, nil
}
