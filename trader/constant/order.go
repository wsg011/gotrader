package constant

// ORDER TYPE
type OrderType int

func (t OrderType) Name() string {
	switch t {
	case OrderTypeLimit:
		return "LIMIT"
	case OrderTypeMarket:
		return "MARKET"
	case OrderTypeIoc:
		return "IOC"
	}
	return "unknown_orderType"
}

const (
	OrderTypeLimit OrderType = iota
	OrderTypeMarket
	OrderTypeIoc
)

// ORDER SIDE
type OrderSide int

func (s OrderSide) Name() string {
	switch s {
	case OrderBuy:
		return "BUY"
	case OrderSell:
		return "SELL"
	}
	return "unknown_orderSide"
}

const (
	OrderBuy OrderSide = iota
	OrderSell
)

// ORDER STATUS
type OrderStatus int

func (s OrderStatus) Name() string {
	switch s {
	case OrderSubmit:
		return "submit"
	case OrderComfirmed:
		return "confirmed"
	case OrderPartialFilled:
		return "partialFilled"
	case OrderFilled:
		return "filled"
	case OrderFailed:
		return "failed"
	case OrderCanceled:
		return "cancelled"
	case OrderClosed:
		return "closed"
	}
	return "unknown_orderStatus"
}

func (s OrderStatus) IsOver() bool {
	return s == OrderFilled ||
		s == OrderFailed ||
		s == OrderCanceled ||
		s == OrderClosed
}

const (
	OrderSubmit OrderStatus = iota
	OrderComfirmed
	OrderPartialFilled
	OrderFilled
	OrderFailed
	OrderCanceled
	OrderClosed
)
