package binanceportfolio

import (
	"encoding/json"
	"net/http"
)

type AutoCollectionResponse struct {
	Message string `json:"msg"`
}

// AutoCllection 自动归集资金
func (client *RestClient) AutoCollection() (string, error) {
	uri := AutoCollectionUri
	body, _, err := client.HttpRequest(http.MethodPost, uri, nil)
	if err != nil {
		log.Errorf("binance post /papi/v1/auto-collection err: %v", err)
		return "false", err
	}

	var response AutoCollectionResponse
	if err = json.Unmarshal(body, &response); err != nil {
		log.Errorf("binance FetchListenKey parser err:%v", err)
		return "false", err
	}

	return response.Message, nil
}
