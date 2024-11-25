package types

type Offer struct {
	IID                  int32
	ID                   string `json:"ID"`
	Data                 string `json:"data"`
	MostSpecificRegionID int32  `json:"mostSpecificRegionID"`
	StartDate            int64  `json:"startDate"`
	EndDate              int64  `json:"endDate"`
	NumberSeats          int32  `json:"numberSeats"`
	Price                int32  `json:"price"`
	CarType              string `json:"carType"`
	HasVollkasko         bool   `json:"hasVollkasko"`
	FreeKilometers       int32  `json:"freeKilometers"`
}

type SearchResultOffer struct {
	ID   string `json:"ID"`
	Data string `json:"data"`
}

type OptimizedSearchResultOffer struct {
	Data []byte
}

func (u OptimizedSearchResultOffer) MarshalJSON() ([]byte, error) {
	return u.Data, nil
}

var MockOffers = []Offer{
	{
		ID:                   "1e8400e2-29b4-4d41-716a-446655440000",
		Data:                 `{"description":"Special discount offer"}`,
		MostSpecificRegionID: 100,
		StartDate:            1700000000000,
		EndDate:              1700003600000,
		NumberSeats:          5,
		Price:                1500,
		CarType:              "Sedan",
		HasVollkasko:         true,
		FreeKilometers:       100,
	},
	{
		ID:                   "2e8400e2-29b4-4d41-716a-446655440001",
		Data:                 `{"description":"Weekend getaway special"}`,
		MostSpecificRegionID: 200,
		StartDate:            1700010000000,
		EndDate:              1700020000000,
		NumberSeats:          4,
		Price:                2500,
		CarType:              "SUV",
		HasVollkasko:         true,
		FreeKilometers:       200,
	},
	{
		ID:                   "3e8400e2-29b4-4d41-716a-446655440002",
		Data:                 `{"description":"Luxury drive experience"}`,
		MostSpecificRegionID: 300,
		StartDate:            1700020000000,
		EndDate:              1700030000000,
		NumberSeats:          2,
		Price:                5000,
		CarType:              "Convertible",
		HasVollkasko:         false,
		FreeKilometers:       300,
	},
	{
		ID:                   "4e8400e2-29b4-4d41-716a-446655440003",
		Data:                 `{"description":"Budget-friendly compact"}`,
		MostSpecificRegionID: 400,
		StartDate:            1700030000000,
		EndDate:              1700040000000,
		NumberSeats:          4,
		Price:                1000,
		CarType:              "Hatchback",
		HasVollkasko:         true,
		FreeKilometers:       150,
	},
	{
		ID:                   "5e8400e2-29b4-4d41-716a-446655440004",
		Data:                 `{"description":"Adventure-ready truck"}`,
		MostSpecificRegionID: 500,
		StartDate:            1700040000000,
		EndDate:              1700050000000,
		NumberSeats:          6,
		Price:                4000,
		CarType:              "Truck",
		HasVollkasko:         true,
		FreeKilometers:       500,
	},
}
