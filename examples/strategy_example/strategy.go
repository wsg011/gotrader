package main

import (
	"fmt"
	"gotrader/trader"
)

type MockStrategy struct {
	pricingChan chan struct{}
	stopChan    chan struct{}
}

func NewStrategy() *MockStrategy {
	return &MockStrategy{
		pricingChan: make(chan struct{}),
		stopChan:    make(chan struct{}),
	}
}

func (ms MockStrategy) OnOrderBook(data trader.OrderBook) {
	// 测试策略逻辑，例如打印数据
	fmt.Println("Strategy OrderBook data:", data)

	// 通知定价
	ms.pricingChan <- struct{}{}
}

func (ms MockStrategy) OnBookTicker(data trader.BookTicker) {
	// 测试策略逻辑，例如打印数据
	fmt.Println("Strategy BookTicker data:", data)

	// 通知定价
	ms.pricingChan <- struct{}{}
}

func (ms MockStrategy) OnTrade(data trader.Trade) {
	// 测试策略逻辑，例如打印数据
	fmt.Println("Strategy Trade data:", data)
}

func (ms MockStrategy) Prepare() {

}

func (ms MockStrategy) Start() {
	ms.Prepare()

	go ms.Run()
}

func (ms MockStrategy) Run() {
	fmt.Println("Run Strategy")
	for {
		select {
		case <-ms.pricingChan:
			ms.Pricing()
		case <-ms.stopChan:
			return
		}
	}
}

func (ms MockStrategy) Close() {
	close(ms.stopChan)
}
