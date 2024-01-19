package main

import (
	"gotrader/exchange"
	"gotrader/trader/constant"
	"gotrader/trader/types"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/BurntSushi/toml"
)

type GlobalConfig struct {
	Exchange map[string]ExchangeConfig `toml:"exchange"`
	Symbols  []string                  `toml:"symbols"`
}

type ExchangeConfig struct {
	Name      string `toml:"name"`
	ApiKey    string `toml:"api_key"`
	SecretKey string `toml:"secret_key"`
	Passphase string `toml:"passphase"`
}

func ReadConfig(path string) (*GlobalConfig, error) {
	var config GlobalConfig
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	_, filename, _, _ := runtime.Caller(0)
	configFile := filepath.Join(filepath.Dir(filename), "config.toml")

	globalConfig, err := ReadConfig(configFile)
	if err != nil {
		log.Errorf("load global config errr: %s", err)
		return
	}
	log.Infof("load global config, symbols: %v", globalConfig.Symbols)

	// 初始化交易所
	okxSwapParams := &types.ExchangeParameters{
		AccessKey:  globalConfig.Exchange["okx"].ApiKey,
		SecretKey:  globalConfig.Exchange["okx"].SecretKey,
		Passphrase: globalConfig.Exchange["okx"].Passphase,
	}
	okxSwap := exchange.NewExchange(constant.OkxV5Swap, okxSwapParams)
	okxSpot := exchange.NewExchange(constant.OkxV5Spot, okxSwapParams)
	log.Infof("init exchange %s", okxSwap.GetName())

	balance, err := okxSpot.FetchBalance()
	if err != nil {
		log.Errorf("fetch balance error %s", err)
		return
	}
	log.Infof("balance %v", balance.TotalUsdEq)

	// 创建策略
	symbols := globalConfig.Symbols
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
