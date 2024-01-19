package okxv5

import (
	"gotrader/pkg/ws"
	"gotrader/trader/types"
	"testing"
	"time"
)

func TestPriWsClient(t *testing.T) {
	testRspHandle := func(data interface{}) {
		// t.Logf("testRspHandle data: %s", data)
	}
	okxWsClient := NewOkPriWsClient("ae02b9cd-6834-462c-90fe-c183dbd4e12b",
		"C629361E13A91F7DAE3E20D4DE208093",
		"zhang2024Q.", testRspHandle)
	if err := okxWsClient.Dial(ws.Connect); err != nil {
		t.Fatalf("Dial error: %s", err)
	}
	t.Log("Dial success")
	time.Sleep(3 * time.Second)

	param := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "bbo-tbt",
				"instId":  "BTC-USDT",
			},
		},
	}
	okxWsClient.Write(param)

	time.Sleep(3 * time.Second)
	okxWsClient.Close()
}

func TestExchange(t *testing.T) {
	params := &types.ExchangeParameters{
		AccessKey:  "",
		SecretKey:  "",
		Passphrase: "",
	}
	exchange := NewOkxV5Swap(params)
	name := exchange.GetName()
	t.Logf("init exchange %s", name)

	param := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "bbo-tbt",
				"instId":  "BTC-USDT",
			},
		},
	}

	exchange.Subscribe(param)
}

func TestSwapWs(t *testing.T) {
	testRspHandle := func(data interface{}) {
		// t.Logf("testRspHandle data: %s", data)
	}
	okxWsClient := NewOkPubWsClient(testRspHandle)
	if err := okxWsClient.Dial(ws.Connect); err != nil {
		t.Fatalf("Dial error: %s", err)
	}
	t.Log("Dial success")
	param := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "bbo-tbt",
				"instId":  "BTC-USDT",
			},
		},
	}
	okxWsClient.Write(param)

	time.Sleep(3 * time.Second)
	okxWsClient.Close()
}

func TestFetchTickers(t *testing.T) {
	// Test the HttpRequest function
	client := NewRestClient("", "", "")
	resp, err := client.FetchTickers()
	if err != nil {
		t.Fatalf("HttpRequest failed: %v", err)
	}

	t.Log(resp[0])
}
