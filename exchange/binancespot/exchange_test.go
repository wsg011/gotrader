package binancespot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

func TestNewBinanceSpot(t *testing.T) {
	params := &types.ExchangeParameters{
		AccessKey:  "your_access_key",
		SecretKey:  "your_secret_key",
		Passphrase: "your_passphrase",
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
