package constant

const (
	ASK = "ask"
	BID = "bid"

	BUY  = "BUY"
	SELL = "SELL"
)

type RoleType int

func (e RoleType) Name() string {
	switch e {
	case Taker:
		return "taker"
	case Maker:
		return "maker"
	}
	return "unknownRole"
}

const (
	Taker RoleType = iota
	Maker
)

const (
	SPOT = "SPOT"
	PERP = "PERP"
)
