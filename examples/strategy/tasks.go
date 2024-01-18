package main

import (
	"gotrader/pkg/utils"
	"time"
)

func (s *MakerStrategy) AddTasks() {
	s.cron.AddFunc("*/1 * * * *", func() {
		r := utils.GenerateRangeNum(5, 30)
		time.Sleep(time.Duration(r) * time.Second)
		log.Infof("1 min healcheck. %s", s.config.Symbol)
	})

	s.cron.Start()
	log.Infof("start cron")
}
