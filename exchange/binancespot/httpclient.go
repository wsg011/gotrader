package binancespot

import (
	"fmt"
	"net/http"
	"time"

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
func (client *RestClient) HttpRequest(method string, uri string, param map[string]interface{}) ([]byte, *http.Response, error) {
	header := map[string]string{
		"X-MBX-APIKEY": client.apiKey,
	}
	param["timestamp"] = time.Now().UnixMilli() - 1000
	toSignStr := utils.UrlEncodeParams(param)
	signature := utils.GenHexDigest(utils.HmacSha256(toSignStr, client.secretKey))
	url := fmt.Sprintf("%s%s?%s&signature=%s", RestUrl, uri, toSignStr, signature)
	args := &httpx.Request{
		Url:    url,
		Head:   header,
		Method: method,
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
