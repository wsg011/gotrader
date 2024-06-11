package binanceportfolio

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	KeepLiveTime = 20 * time.Minute
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

func (client *RestClient) RefreshListenKey(key string) error {
	url := FetchListenKey
	param := make(map[string]interface{})
	param["listenKey"] = key
	_, _, err := client.HttpRequest(http.MethodPut, url, param)
	if err != nil {
		log.Errorf("RefreshListenKey err: %s", err)
		return err
	}
	log.Infof("RefreshListenKey success: %s", key)

	return nil
}

func (client *RestClient) KeepUserStream(key string) {
	timer := time.NewTimer(30 * time.Minute) // 创建计时器一次

	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			client.RefreshListenKey(key)
			timer.Reset(30 * time.Minute)
		case <-client.stopChan: // 假设有一个stopChan用于接收退出信号
			return
		}
	}
}
