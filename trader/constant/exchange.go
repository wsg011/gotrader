package constant

import "fmt"

const (
	Exchange_PionexSpot  = "pionexSpot"
	Exchange_OkxV5Spot   = "okxV5Spot"
	Exchange_OkxV5Future = "okxV5Future"
	Exchange_OkxV5Swap   = "okxV5Swap"
)

type ExchangeType int

func (e ExchangeType) Name() string {
	switch e {
	case PionexSpot:
		return Exchange_PionexSpot
	case OkxV5Spot:
		return Exchange_OkxV5Spot
	case OkxV5Swap:
		return Exchange_OkxV5Swap
	}
	return "unknown"
}

const (
	PionexSpot ExchangeType = iota
	OkxV5Spot
	OkxV5Future
	OkxV5Swap
)

func MustConverToExchangeType(name string) ExchangeType {
	switch name {
	case Exchange_PionexSpot:
		return PionexSpot
	case Exchange_OkxV5Spot:
		return OkxV5Spot
	case Exchange_OkxV5Swap:
		return OkxV5Future
	}
	err := fmt.Errorf("unknonw exchange name:%s", name)
	panic(err)
}
