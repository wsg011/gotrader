package main

import (
	"gotrader/event"
	"gotrader/exchange"
	"gotrader/trader"
	"gotrader/trader/constant"
	"gotrader/trader/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("main", "strategy")

func main() {
	eventEngine := event.NewEventEngine()
	onBookTickerCallback := func(bookticker *types.BookTicker) {
		eventEngine.Push(constant.EVENT_BOOKTICKER, bookticker)
	}

	// 创建 TraderEngine 实例
	traderEngine := trader.NewTraderEngine(eventEngine)

	// 初始化交易所
	okxSwap := exchange.NewExchange(constant.OkxV5Swap)
	log.Infof("init exchange %s", okxSwap.GetName())

	// 创建策略
	symbols := []string{"BTC_USDT", "ETH_USDT"}
	for _, symbol := range symbols {
		// Sub datafeed
		err := okxSwap.SubscribeBookTicker([]string{symbol}, onBookTickerCallback)
		if err != nil {
			log.Errorf("SubscribeBookTicker err %s", err)
		}

		mockStrategy := NewStrategy()
		traderEngine.AddStrategy(mockStrategy)
	}

	// 启动
	traderEngine.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		log.Infof("Got %s signal. Aborting...\n", sig)
		// Ensure traderEngine has a Stop method
		traderEngine.Stop()
		os.Exit(1)
	}()

	select {}
}
