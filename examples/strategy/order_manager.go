package main

import (
	"gotrader/trader/types"
)

type OrderManager struct {
	askOrders []types.Order
	bidOrders []types.Order
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		askOrders: make([]types.Order, 0),
		bidOrders: make([]types.Order, 0),
	}
}

func (oms *OrderManager) MatchOrders(askOrders, bidOrders []types.Order) ([]types.Order, []types.Order) {
	// create orders
	createOrders := make([]types.Order, 0)
	createOrders = append(createOrders, askOrders...)
	createOrders = append(createOrders, bidOrders...)

	// cancel order
	cancelOrders := make([]types.Order, 0)
	if oms.askOrders != nil {
		cancelOrders = append(cancelOrders, oms.askOrders...)
		cancelOrders = append(cancelOrders, oms.bidOrders...)
	}

	return createOrders, cancelOrders
}
