package types

import (
	"gotrader/trader/constant"
	"strconv"
)

type Order struct {
	Symbol        string             `json:"symbol"`
	Type          constant.OrderType `json:"type"`
	OrderID       string             `json:"orderId"`
	ClientID      string             `json:"clientID"`
	Side          constant.OrderSide `json:"side"`
	Price         string             `json:"price"`
	OrigQty       string             `json:"origQty"`
	Amount        string             `json:"amount"` // 订单quote额度
	ExecutedQty   string             `json:"executedQty"`
	ExecutedAmt   string
	AvgPrice      string                `json:"avgPrice"`
	Fee           string                `json:"fee"`
	Status        constant.OrderStatus  `json:"status"` // 自定义的订单状态，统一各交易所订单状态
	ExchangeType  constant.ExchangeType `json:"exchangeType"`
	HedgeClientId string
	HedgingPrice  string
	ExpectSpred   float64
	CreateAt      int64 `gorm:"column:createAt;" json:"createAt"`
	UpdateAt      int64 `gorm:"column:updateAt;" json:"updateAt"`
}

func (order *Order) IsOver() bool {
	return order.Status == constant.OrderFilled ||
		order.Status == constant.OrderFailed ||
		order.Status == constant.OrderCanceled ||
		order.Status == constant.OrderClosed
}

func (order *Order) UpdateByTradeEvent(trade *Trade) (float64, float64) {
	executedQty := float64(0)
	if order.ExecutedQty != "" && order.ExecutedQty != "0" {
		executedQty, _ = strconv.ParseFloat(order.ExecutedQty, 64)
	}

	totalFilled := float64(0)
	if trade.FilledSize != "" && trade.FilledSize != "0" {
		totalFilled, _ = strconv.ParseFloat(trade.FilledSize, 64)
	}

	executedAmt := float64(0)
	if order.ExecutedAmt != "" && order.ExecutedAmt != "0" {
		executedAmt, _ = strconv.ParseFloat(order.ExecutedAmt, 64)
	}

	totalFilledAmt := float64(0)
	if trade.FilledAmount != "" && trade.FilledAmount != "0" {
		totalFilledAmt, _ = strconv.ParseFloat(trade.FilledAmount, 64)
	}

	order.ExecutedQty = trade.FilledSize
	order.ExecutedAmt = trade.FilledAmount
	order.Fee = trade.Fee
	if trade.FilledSize == trade.Size {
		order.Status = constant.OrderFilled
		log.WithField("trade", *trade).Info("order filled")
	} else if trade.Status == constant.OrderClosed {
		order.Status = constant.OrderClosed
		log.WithField("trade", *trade).Info("order partial filled")
	} else {
		order.Status = constant.OrderPartialFilled
		log.WithField("trade", *trade).Debug("order partial filling")
	}
	return totalFilled - executedQty, totalFilledAmt - executedAmt
}
