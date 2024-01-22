package main

import (
	"github.com/wsg011/gotrader/trader/types"
)

type DemoStrategy struct {
	name string

	bookTicker *types.BookTicker

	stopChan chan struct{}
}

func NewStrategy(name string) *DemoStrategy {
	return &DemoStrategy{
		name:     name,
		stopChan: make(chan struct{}),
	}
}

func (s *DemoStrategy) Start() {
	go s.Run()
}

func (s *DemoStrategy) OnBookTicker(bookticker *types.BookTicker) {
	s.bookTicker = bookticker
	log.Infof("%s Bookticker [%f:%f]", s.bookTicker.Symbol, s.bookTicker.BidPrice, s.bookTicker.BidPrice)
}

func (s *DemoStrategy) Run() {
	log.Infof("Run Strategy")
	for {
		<-s.stopChan
		return
	}
}

func (s *DemoStrategy) Close() {
	log.Infof("Close Strategy")
	close(s.stopChan)
}
