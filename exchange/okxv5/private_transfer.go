package okxv5

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wsg011/gotrader/exchange/base"
)

type PrivateTransferResult struct {
	TransId     string `json:"transId"`
	Ccy         string `json:"ccy"`
	ClientId    string `json:"clientId"`
	FromAccount string `json:"from"`
	Amount      string `json:"amt"`
	ToAccount   string `json:"to"`
}

type PrivateTransferResponse struct {
	BaseOkRsp
	Data []*PrivateTransferResult `json:"data"`
}

func (client *RestClient) PrivateTransfer(transfer base.TransferParam) (string, error) {
	param := map[string]interface{}{
		"from": TypeMap[transfer.FromType],
		"to":   TypeMap[transfer.ToType],
		"ccy":  strings.ToUpper(transfer.Assert),
		"amt":  transfer.Amount,
	}
	if transfer.TransferType != "0" {
		param["subAcct"] = transfer.ToAccount
		param["type"] = TransferMap[transfer.TransferType]
	}
	payload, _ := json.Marshal(param)
	uri := PrivateTransferUri
	body, _, err := client.HttpRequest(http.MethodPost, uri, payload)
	if err != nil {
		log.Errorf("okx post /api/v5/asset/transfer 请求失败: %v", err)
		return "", err
	}

	response := new(PrivateTransferResponse)
	if err = json.Unmarshal(body, response); err != nil {
		log.Errorf("okx post /api/v5/asset/transfer err 数据解析失败: %v", err)
		return "", err
	}

	if response.Code != "0" {
		err := fmt.Errorf("%s", response.Msg)
		log.Errorf("okx post /api/v5/asset/transfer 划转失败: %s", err)
		return "nil", err
	}

	transferId := response.Data[0].TransId
	return transferId, nil
}
