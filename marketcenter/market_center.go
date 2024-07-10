package marketcenter

import (
	"sync"

	"github.com/wsg011/gotrader/trader/types"
)

type SubscriptionKey struct {
	Exchange string
	Symbol   string
	Topic    string
}

type MarketCenter struct {
	subscribers map[SubscriptionKey][]func(interface{}) // 订阅者映射，存储回调函数
	mu          sync.Mutex
}

func NewMarketCenter() *MarketCenter {
	return &MarketCenter{
		subscribers: make(map[SubscriptionKey][]func(interface{})),
	}
}

// 订阅行情，使用回调函数接收数据
func (center *MarketCenter) Subscribe(exchange, symbol, topic string, callbackFunc func(data interface{})) {
	center.mu.Lock()
	defer center.mu.Unlock()

	key := SubscriptionKey{Exchange: exchange, Symbol: symbol, Topic: topic}
	center.subscribers[key] = append(center.subscribers[key], callbackFunc)
}

// 发布行情
func (center *MarketCenter) Publish(data interface{}) {
	center.mu.Lock()
	defer center.mu.Unlock()

	// 这里假设数据类型包含Exchange, Symbol, 和 Topic字段
	exchange, symbol, topic := "", "", ""
	switch v := data.(type) {
	case types.BookTicker:
		exchange = v.Exchange.Name()
		symbol = v.Symbol
		topic = "BookTicker"
	case types.Order:
		exchange = v.Exchange.Name()
		symbol = v.Symbol
		topic = "Order"
	}

	key := SubscriptionKey{Exchange: exchange, Symbol: symbol, Topic: topic}
	if callbacks, found := center.subscribers[key]; found {
		for _, callback := range callbacks {
			go callback(data) // 异步调用回调函数
		}
	}
}
