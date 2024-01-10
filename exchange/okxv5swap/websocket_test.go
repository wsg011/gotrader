package okxv5swap

import (
	"gotrader/pkg/ws"
	"testing"
	"time"
)

func TestSwwapWs(t *testing.T) {
	okxWsClient := NewOkPubWsClient()
	if err := okxWsClient.Dial(ws.Connect); err != nil {
		t.Fatalf("Dial error: %s", err)
	}
	t.Log("Dial success")

	okxWsClient.Subscribe("BTC-USDT", "bbo-tbt")
	okxWsClient.Subscribe("ETH-USDT", "bbo-tbt")
	okxWsClient.Subscribe("SOL-USDT", "bbo-tbt")

	time.Sleep(10 * time.Second)
}
