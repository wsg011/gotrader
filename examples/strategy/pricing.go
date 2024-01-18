package main

func (s *MakerStrategy) Pricing() {
	if s.vars.BookTicker == nil || (s.vars.HedgeBookTicker == nil) {
		return
	}
	ask_price, bid_price := s.vars.HedgeBookTicker.AskPrice, s.vars.HedgeBookTicker.BidPrice
	if s.vars.epoch%100 == 0 {
		log.Infof("Pricing [%f:%f]", ask_price, bid_price)
	}
}
