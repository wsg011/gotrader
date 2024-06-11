package binancespot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

func TestNewBinanceSpot(t *testing.T) {
	params := &types.ExchangeParameters{
		AccessKey:  "",
		SecretKey:  "",
		Passphrase: "",
	}

	exchange := NewBinanceSpot(params)

	assert.NotNil(t, exchange, "Expected exchange to not be nil")
	assert.Equal(t, constant.OkxV5Swap, exchange.exchangeType, "Exchange type mismatch")
	assert.NotNil(t, exchange.restClient, "Expected non-nil RestClient")
}

func TestFetchSymbols(t *testing.T) {
	exchange := &BinanceSpotExchange{
		exchangeType: constant.OkxV5Swap,
	}
	symbolinfos, err := exchange.FetchSymbols()
	if err != nil {
		assert.Error(t, err, "Expected an error to be returned")
	}
	assert.NotNil(t, symbolinfos, "Expected no symbols to be returned")
	for index, symbolinfo := range symbolinfos {
		t.Logf("symbol info %v", symbolinfo)
		if index > 1 {
			break
		}
	}
}

func TestFetchBalance(t *testing.T) {
	params := &types.ExchangeParameters{
		// AccessKey:  "",
		// SecretKey:  "",
		AccessKey:  "Yo0gUrDtgMCcEQSxK4v6vQg90qIU1O3NZX1VKkUt1PBDu0r9Pu1PsrM1OJnooXZg",
		SecretKey:  "KHCNySha8EpwnIUDn6KyAEw1G7mFrp0MOzjWK6SWStWyMDITqx6xxa1Q6BKTVski",
		Passphrase: "",
	}

	exchange := NewBinanceSpot(params)

	balance, err := exchange.FetchBalance()
	if err != nil {
		t.Errorf("FetchBalance error %s", err)
	} else {
		t.Logf("FetchBalance result %v", balance)
	}

}

func TestSubscribeBookTicker(t *testing.T) {
	params := &types.ExchangeParameters{
		AccessKey:  "",
		SecretKey:  "",
		Passphrase: "",
	}

	exchange := NewBinanceSpot(params)
	callback := func(b *types.BookTicker) {
		t.Logf("Received: %+v", b)
	}

	err := exchange.SubscribeBookTicker([]string{"BTC_USDT"}, callback)
	if err != nil {
		t.Errorf("SubscribeBookTicker err: %s", err)
	}
	time.Sleep(20 * time.Second) // 等待5秒

	tickers, err := exchange.FetchTickers()
	if err != nil {
		t.Errorf("FetchTickers err %s", err)
	} else {
		t.Logf("tickers len %v", len(tickers))
		t.Logf("tickers %+v", tickers[0])
	}
}
