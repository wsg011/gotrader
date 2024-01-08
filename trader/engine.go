package trader

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type TraderEngine struct {
	strategiesMap   map[string]Strategy
	publicDataChan  <-chan MarketDataInterface
	privateDataChan <-chan MarketDataInterface
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

var logger = logrus.WithField("main", "engine")

func NewTraderEngine(publicDataChan <-chan MarketDataInterface, privateDataChan <-chan MarketDataInterface) *TraderEngine {
	return &TraderEngine{
		strategiesMap:   make(map[string]Strategy),
		publicDataChan:  publicDataChan,
		privateDataChan: privateDataChan,
		stopChan:        make(chan struct{}),
	}
}

// 添加策略及其对应的交易品种
func (engine *TraderEngine) AddStrategy(strategy Strategy, symbols []string) {
	for _, symbol := range symbols {
		engine.strategiesMap[symbol] = strategy
	}
}

func (engine *TraderEngine) Start() {
	engine.wg.Add(2)
	go engine.processData(engine.publicDataChan)
	go engine.processData(engine.privateDataChan)
}

func (engine *TraderEngine) processData(dataChan <-chan MarketDataInterface) {
	defer engine.wg.Done()
	for {
		select {
		case data := <-dataChan:
			symbol := data.GetSymbol()
			if strategy, exists := engine.strategiesMap[symbol]; exists {
				switch d := data.(type) {
				case BookTicker:
					strategy.OnBookTicker(d)
				case OrderBook:
					strategy.OnOrderBook(d)
				case Trade:
					strategy.OnTrade(d)
				}
			}
		case <-engine.stopChan:
			logger.Infof("stopChan")
			return
		}
	}
}

func (engine *TraderEngine) Stop() {
	close(engine.stopChan)

	// 等待子线程执行完毕
	engine.wg.Wait()
}
