package constant

// ORDER TYPE
type OrderType int

func (t OrderType) Name() string {
	switch t {
	case Limit:
		return "LIMIT"
	case Market:
		return "MARKET"
	case IOC:
		return "IOC"
	case GTC:
		return "GTC"
	case FOK:
		return "FOK"
	case PostOnly:
		return "POST_ONLY"
	}
	return "unknown_orderType"
}

const (
	Limit OrderType = iota
	Market
	IOC
	GTC
	FOK
	PostOnly
)

// ORDER SIDE
type OrderSide int

func (s OrderSide) Name() string {
	switch s {
	case OrderBuy:
		return "BUY"
	case OrderSell:
		return "SELL"
	case Long:
		return "LONG"
	case Short:
		return "SHORT"
	case CloseLong:
		return "CLOSELONG"
	case CloseShort:
		return "SHORT"
	case All:
		return "All"
	}
	return "unknown_orderSide"
}

const (
	OrderBuy OrderSide = iota
	OrderSell
	Long
	Short
	CloseLong
	CloseShort
	All
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
