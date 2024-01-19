package main

import (
	"fmt"
	"gotrader/pkg/utils"
	"gotrader/trader/constant"
	"gotrader/trader/types"
	"strconv"
	"time"
)

func (s *MakerStrategy) Pricing() {
	if s.vars.BookTicker == nil || (s.vars.HedgeBookTicker == nil) {
		return
	}
	ask_price, bid_price := s.vars.HedgeBookTicker.AskPrice, s.vars.HedgeBookTicker.BidPrice

	ask_price = ask_price * (1 + s.vars.basisMean + 2*s.vars.basisStd - s.vars.fundingRate + 0.0004)
	bid_price = bid_price * (1 + s.vars.basisMean + 2*s.vars.basisStd - s.vars.fundingRate - 0.0004)
	if s.vars.epoch%100 == 0 {
		log.Infof("Swap [%f:%f]", s.vars.BookTicker.AskPrice, s.vars.BookTicker.BidPrice)
		log.Infof("Spot [%f:%f]", s.vars.HedgeBookTicker.AskPrice, s.vars.HedgeBookTicker.BidPrice)
		log.Infof("Pricing [%f:%f]", ask_price, bid_price)
	}

	askOrders, bidOrders := s.genOrders(ask_price, bid_price)
	place_orders, cancel_orders := s.OMS.MatchOrders(askOrders, bidOrders)

	if s.vars.epoch%100 == 0 {
		log.Infof("place_orders %v", place_orders)
		log.Infof("cancel_orders %v", cancel_orders)
	}
}

func (s *MakerStrategy) genOrders(askPrice, biPrice float64) ([]types.Order, []types.Order) {
	askPriceStr := strconv.FormatFloat(askPrice, 'f', 2, 64)
	qty := "1.0"
	askOrder := types.Order{
		Symbol:       s.config.Symbol,
		Type:         constant.OrderTypeLimit,
		ClientID:     "ijijijoj",
		Side:         constant.OrderSell,
		Price:        askPriceStr,
		OrigQty:      qty,
		ExchangeType: s.config.MakerExchange.GetType(),
		CreateAt:     utils.Millisec(time.Now()),
		Status:       constant.OrderSubmit,
		HedgingPrice: fmt.Sprint(s.vars.HedgeBookTicker.AskPrice),
	}

	askOrders := make([]types.Order, 0, 2)
	askOrders = append(askOrders, askOrder)
	return askOrders, askOrders
}
