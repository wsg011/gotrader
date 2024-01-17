package trader

import (
	"gotrader/event"
	"sync"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("main", "engine")

type TraderEngine struct {
	strategies  []Strategy
	eventEngine *event.EventEngine
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewTraderEngine(eventEngine *event.EventEngine) *TraderEngine {
	return &TraderEngine{
		strategies:  make([]Strategy, 0),
		eventEngine: eventEngine,
		stopChan:    make(chan struct{}),
	}
}

// AddStrategy 增加策略，并订阅EventEngine事件
func (trader *TraderEngine) AddStrategy(strategy Strategy) {
	trader.strategies = append(trader.strategies, strategy)

	// trader.eventEngine.Register(constant.EVENT_BOOKTICKER, strategy.OnBookTicker)
}

func (trader *TraderEngine) Start() {
	log.Infof("Start trader")
	for _, strategy := range trader.strategies {
		strategy.Start()
	}

	// TODO 定时检查策略状态
}

func (trader *TraderEngine) Stop() {
	close(trader.stopChan)

	// 等待子线程执行完毕
	trader.wg.Wait()
}
