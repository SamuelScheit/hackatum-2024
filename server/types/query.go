package types

type QueryResponse struct {
	Offers             OptimizedSearchResultOffer `json:"offers" validate:"required"`
	PriceRanges        []PriceRange               `json:"priceRanges" validate:"required"`
	CarTypeCounts      CarTypeCount               `json:"carTypeCounts" validate:"required"`
	SeatsCount         []SeatsCount               `json:"seatsCount" validate:"required"`
	FreeKilometerRange []FreeKilometerRange       `json:"freeKilometerRange" validate:"required"`
	VollkaskoCount     VollkaskoCount             `json:"vollkaskoCount" validate:"required"`
}

type PriceRange struct {
	Start uint `json:"start"`
	End   uint `json:"end"`
	Count uint `json:"count"`
}

type CarTypeCount struct {
	Small  uint `json:"small"`
	Sports uint `json:"sports"`
	Luxury uint `json:"luxury"`
	Family uint `json:"family"`
}

type SeatsCount struct {
	NumberSeats uint `json:"numberSeats"`
	Count       uint `json:"count"`
}

type FreeKilometerRange struct {
	Start uint `json:"start"`
	End   uint `json:"end"`
	Count uint `json:"count"`
}

type VollkaskoCount struct {
	TrueCount  uint `json:"trueCount"`
	FalseCount uint `json:"falseCount"`
}
