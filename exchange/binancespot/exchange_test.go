package binancespot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

func TestNewBinanceSpot(t *testing.T) {
	params := &types.ExchangeParameters{
		// AccessKey:  "your_access_key",
		// SecretKey:  "your_secret_key",
		// Passphrase: "your_passphrase",
		AccessKey:  "Yo0gUrDtgMCcEQSxK4v6vQg90qIU1O3NZX1VKkUt1PBDu0r9Pu1PsrM1OJnooXZg",
		SecretKey:  "KHCNySha8EpwnIUDn6KyAEw1G7mFrp0MOzjWK6SWStWyMDITqx6xxa1Q6BKTVski",
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
	// assert.Equal(t, constant.OkxV5Swap, typ, "Exchange type mismatch")
}

func TestFetchBalance(t *testing.T) {
	params := &types.ExchangeParameters{
		// AccessKey:  "your_access_key",
		// SecretKey:  "your_secret_key",
		// Passphrase: "your_passphrase",
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
