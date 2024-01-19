package main

import (
	"gotrader/pkg/utils"
	"time"

	"github.com/montanaflynn/stats"
)

func (s *MakerStrategy) AddTasks() {
	// Run task
	s.UpdateBasis()
	s.UpdateFundingRate()

	s.cron.AddFunc("*/1 * * * *", func() {
		r := utils.GenerateRangeNum(5, 30)
		time.Sleep(time.Duration(r) * time.Second)
		s.UpdateBasis()
	})
	s.cron.AddFunc("@hourly", func() {
		r := utils.GenerateRangeNum(5, 30)
		time.Sleep(time.Duration(r) * time.Second)
		s.UpdateFundingRate()
	})

	s.cron.Start()
	log.Infof("start cron")
}

func (s *MakerStrategy) UpdateBasis() {
	var limit int64 = 100
	spotKline, err := s.config.MakerExchange.FetchKline(s.config.Symbol, "1m", limit)
	if err != nil {
		log.Errorf("spot FetchKline error %s", err)
		return
	}
	swapKline, err := s.config.MakerExchange.FetchKline(s.config.Symbol+"-SWAP", "1m", limit)
	if err != nil {
		log.Errorf("swap FetchKline error %s", err)
		return
	}
	if len(spotKline) == 0 || len(swapKline) == 0 {
		log.Errorf("Kline empty")
		return
	}
	if len(spotKline) != len(swapKline) {
		log.Errorf("Kline line not equ")
		return
	}

	spread := make([]float64, 0, len(spotKline))
	for i := 1; i < len(spotKline); i++ {
		s := (swapKline[i].Close - spotKline[i].Close) / spotKline[i].Close
		spread = append(spread, s)
	}

	mean, _ := stats.Mean(spread)
	std, _ := stats.StandardDeviation(spread)
	log.Infof("update basis mean %f std %f", mean, std)
	s.vars.basisMean = mean
	s.vars.basisStd = std
}

func (s *MakerStrategy) UpdateFundingRate() {
	fundingRate, err := s.config.MakerExchange.FetchFundingRate(s.config.Symbol + "-SWAP")
	if err != nil {
		log.Errorf("FetchFundingRate err %s", err)
	}

	s.vars.fundingRate = fundingRate.FundingRate
	log.Infof("update funding rate %v", fundingRate)
}
