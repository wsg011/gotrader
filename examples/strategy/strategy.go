package main

import (
	"fmt"
	"gotrader/pkg/utils"
	"gotrader/trader/types"
	"time"
)

type MakerStrategy struct {
	name string

	config *Config
	vars   *Vars

	OMS *OrderManager

	pricingChan chan struct{}
	stopChan    chan struct{}
}

func NewStrategy(name string, config *Config) *MakerStrategy {
	return &MakerStrategy{
		name:        name,
		config:      config,
		OMS:         &OrderManager{},
		pricingChan: make(chan struct{}),
		stopChan:    make(chan struct{}),
	}
}

func (s *MakerStrategy) GetName() string {
	return "MakerStrategy"
}

func (s *MakerStrategy) OnOrderBook(data *types.OrderBook) {
	// 通知定价
	s.pricingChan <- struct{}{}
}

func (s *MakerStrategy) OnBookTicker(bookTicker *types.BookTicker) {
	s.vars.epoch++

	switch bookTicker.Exchange {
	case s.config.MakerExchange.GetType():
		s.vars.BookTicker = bookTicker
	case s.config.HedgeExchange.GetType():
		s.vars.HedgeBookTicker = bookTicker
	default:
		log.Warnf("unknow exhcange BookTicker data: %s", bookTicker.Exchange.Name())
	}

	if s.vars.epoch%100 == 0 {
		log.Infof("curr bookTicker %s %s, feedDelay %v processDelay %v us",
			bookTicker.Exchange.Name(),
			bookTicker.Symbol,
			bookTicker.Ts-bookTicker.ExchangeTs,
			utils.Microsec(time.Now())-bookTicker.Ts,
		)

	}

	s.pricingChan <- struct{}{}
}

func (s *MakerStrategy) OnTrade(data *types.Trade) {
	fmt.Println("Strategy Trade data:", data)
}

func (s *MakerStrategy) Prepare() {
	// Load Config
	log.Infof("config %v", s.config)

	// init vars
	vars := &Vars{
		epoch: 0,
	}
	s.vars = vars
}

func (s *MakerStrategy) OnOrder(order *types.Order) {

}

func (s *MakerStrategy) Start() {
	s.Prepare()

	go s.Run()
}

func (s *MakerStrategy) Run() {
	log.Infof("Run Strategy")
	for {
		select {
		case <-s.stopChan:
			return
		case <-s.pricingChan:
			s.Pricing()
		}
	}
}

func (s *MakerStrategy) Close() {
	close(s.stopChan)
}
