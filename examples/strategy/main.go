package main

import (
	"gotrader/exchange"
	"gotrader/trader/constant"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化交易所
	okxSwap := exchange.NewExchange(constant.OkxV5Swap)
	okxSpot := exchange.NewExchange(constant.OkxV5Spot)
	log.Infof("init exchange %s", okxSwap.GetName())

	// 创建策略
	symbols := []string{"ACE_USDT", "OP_USDT"}
	for _, symbol := range symbols {
		config := &Config{
			Symbol:        symbol,
			MakerExchange: okxSwap,
			HedgeExchange: okxSpot,
		}
		strategy := NewStrategy("MakerStrategy", config)
		strategy.Start()

		// Sub datafeed
		err := okxSwap.SubscribeBookTicker([]string{symbol, symbol + "_SWAP"}, strategy.OnBookTicker)
		if err != nil {
			log.Errorf("SubscribeBookTicker err %s", err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		log.Infof("Got %s signal. Aborting...\n", sig)
		// Ensure traderEngine has a Stop method
		os.Exit(1)
	}()

	select {}
}
