package main

import (
	"github.com/wsg011/gotrader/exchange"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("main", "strategy")

func main() {
	// 初始化交易所
	okxSwapParams := &types.ExchangeParameters{
		AccessKey:  "",
		SecretKey:  "",
		Passphrase: "",
	}
	okxSwap := exchange.NewExchange(constant.OkxV5Swap, okxSwapParams)

	// 策略初始化
	strategy := NewStrategy("demo")
	strategy.Start()

	// 订阅行情
	err := okxSwap.SubscribeBookTicker([]string{"BTC_USDT"}, strategy.OnBookTicker)
	if err != nil {
		log.Errorf("SubscribeBookTicker err: %s", err)
		return
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
