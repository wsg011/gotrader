package types

type Asset struct {
	Coin   string
	Free   float64
	Frozen float64
	Total  float64
}

type Assets struct {
	Assets     map[string]Asset
	TotalUsdEq float64
}
