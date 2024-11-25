package types

type QueryResponse struct {
	Offers             []SearchResultOffer   `json:"offers" validate:"required"`
	PriceRanges        []*PriceRange         `json:"priceRanges" validate:"required"`
	CarTypeCounts      CarTypeCount          `json:"carTypeCounts" validate:"required"`
	SeatsCount         []*SeatsCount         `json:"seatsCount" validate:"required"`
	FreeKilometerRange []*FreeKilometerRange `json:"freeKilometerRange" validate:"required"`
	VollkaskoCount     VollkaskoCount        `json:"vollkaskoCount" validate:"required"`
}

type PriceRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type CarTypeCount struct {
	Small  int `json:"small"`
	Sports int `json:"sports"`
	Luxury int `json:"luxury"`
	Family int `json:"family"`
}

type SeatsCount struct {
	NumberSeats int `json:"numberSeats"`
	Count       int `json:"count"`
}

type FreeKilometerRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type VollkaskoCount struct {
	TrueCount  int `json:"trueCount"`
	FalseCount int `json:"falseCount"`
}
