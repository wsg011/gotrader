package binancespot

import (
	"net/http"

	"github.com/wsg011/gotrader/pkg/httpx"
	"github.com/wsg011/gotrader/pkg/utils"
	"github.com/wsg011/gotrader/trader/constant"
)

var httpClient = httpx.NewClient()

type RestClient struct {
	apiKey       string
	secretKey    string
	passPhrase   string
	exchangeType constant.ExchangeType
}

type BaseOkRsp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewRestClient(apiKey, secretKey, passPhrase string, exchangeType constant.ExchangeType) *RestClient {
	client := &RestClient{
		apiKey:       apiKey,
		secretKey:    secretKey,
		passPhrase:   passPhrase,
		exchangeType: exchangeType,
	}
	return client
}
func (client *RestClient) HttpRequest(method string, uri string, payload []byte) ([]byte, *http.Response, error) {
	var param string
	if payload != nil {
		param = string(payload)
	}
	currentTime := utils.IsoTime()
	toSignStr := currentTime + method + uri + param
	signature := utils.GenBase64Digest(utils.HmacSha256(toSignStr, client.secretKey))
	url := RestUrl + uri
	head := map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-KEY":        client.apiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  currentTime,
		"OK-ACCESS-PASSPHRASE": client.passPhrase,
	}
	args := &httpx.Request{
		Url:    url,
		Head:   head,
		Method: method,
		Body:   payload,
	}
	body, res, err := httpClient.Request(args)
	if err != nil {
		return nil, res, err
	}
	return *body, res, err
}

func (client *RestClient) HttpGet(url string) ([]byte, *http.Response, error) {
	body, res, err := httpClient.Get(url)
	if err != nil {
		return nil, res, err
	}
	return *body, res, nil
}
