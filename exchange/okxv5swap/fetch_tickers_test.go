package okxv5swap

import (
	"testing"
)

func TestFetchTickers(t *testing.T) {
	// Test the HttpRequest function
	client := NewRestClient("", "", "")
	resp, err := client.FetchTickers()
	if err != nil {
		t.Fatalf("HttpRequest failed: %v", err)
	}

	t.Log(resp[0])
}
