package types

type FundingRate struct {
	Symbol          string
	Method          string
	FundingRate     float64
	FundingTime     int64
	NextFundingRate float64
	NextFundingTime int64
}
