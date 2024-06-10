package binanceportfolio

import (
	"encoding/json"
	"net/http"
)

type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

func (client *RestClient) GetListenKey() (string, error) {
	url := FetchListenKey
	body, _, err := client.HttpRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Errorf("binance FetchListenKey 网络错误:%v", err)
		return "", err
	}

	var response ListenKeyResponse
	if err = json.Unmarshal(body, &response); err != nil {
		log.Errorf("binance FetchListenKey parser err:%v", err)
		return "", err
	}

	return response.ListenKey, nil
}
