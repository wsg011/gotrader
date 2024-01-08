package types

type FundingRate struct {
	FundingRate     float64 `json:"fundingRate"`
	FundingTime     int64   `json:"fundingTime"`
	Symbol          string  `json:"symbol"`
	NextFundingRate float64 `json:"nextFundingRate"`
	NextFundingTime int64   `json:"nextFundingTime"`
}
