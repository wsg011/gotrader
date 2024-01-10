package trader

import (
	"gotrader/event"
	"sync"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("main", "engine")

type TraderEngine struct {
	strategiesMap map[string]Strategy
	eventEngine   *event.EventEngine
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

func NewTraderEngine(eventEngine *event.EventEngine) *TraderEngine {
	return &TraderEngine{
		strategiesMap: make(map[string]Strategy),
		eventEngine:   eventEngine,
		stopChan:      make(chan struct{}),
	}
}

func (trader *TraderEngine) AddStrategy(strategy Strategy, symbols []string) {

}

func (trader *TraderEngine) Start() {
	log.Infof("Start trader")
}

func (trader *TraderEngine) Stop() {
	close(trader.stopChan)

	// 等待子线程执行完毕
	trader.wg.Wait()
}
