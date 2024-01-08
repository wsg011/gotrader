package main

import (
	"gotrader/trader"
	"time"
)

func main() {
	// 设置市场数据通道
	publicDataChan := make(chan trader.MarketDataInterface, 1)
	privateDataChan := make(chan trader.MarketDataInterface, 1)

	// 创建 TraderEngine 实例
	engine := trader.NewTraderEngine(publicDataChan, privateDataChan)

	// 添加策略
	mockStrategy := NewStrategy()
	mockStrategy.Start()

	engine.AddStrategy(mockStrategy, []string{"BTC_USDT"})

	// 启动 TraderEngine
	engine.Start()

	// 发送模拟市场数据
	publicDataChan <- trader.BookTicker{Symbol: "BTC_USDT", BidPrice: 128, AskPrice: 129}
	privateDataChan <- trader.Trade{Symbol: "BTC_USDT", Price: 128.5, Volume: 10}

	// 等待一段时间以确保数据被处理
	time.Sleep(1 * time.Second)

	// 停止 TraderEngine
	engine.Stop()
}
