package trader

import (
	"fmt"
	"gotrader/event"
	"gotrader/trader/constant"
)

type DataFeed struct {
	eventEngine *event.EventEngine
	exchanges   map[constant.ExchangeType]Exchange
}

func NewDataFeed(eventEngine *event.EventEngine) *DataFeed {
	return &DataFeed{
		eventEngine: eventEngine,
		exchanges:   make(map[constant.ExchangeType]Exchange),
	}
}

// ReceiveData 接收数据并推送到EventEngine
func (feed *DataFeed) ReceiveData(data interface{}) {
	feed.eventEngine.Push("Data", data)
}

// AddExchange 添加交易所到数据源
func (feed *DataFeed) AddExchange(exchangeType constant.ExchangeType, exchange Exchange) {
	feed.exchanges[exchangeType] = exchange
}

// Subscribe 订阅交易所行情
func (feed *DataFeed) Subscribe(exchangeType constant.ExchangeType, params map[string]interface{}) error {
	exchange, exists := feed.exchanges[exchangeType]
	if !exists {
		return fmt.Errorf("exchange type %v not found", exchangeType)
	}

	return exchange.Subscribe(params)
}
