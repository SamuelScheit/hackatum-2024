package database

import (
	"checkmate/types"
	"encoding/json"
	"os"
	"testing"
	"time"
)

/*
{
    "OfferID": "3c0b5b47-5373-4071-a335-ed4c9035e236",
    "RegionID": 104,
    "CarType": "luxury",
    "NumberDays": 2,
    "NumberSeats": 2,
    "StartTimestamp": "2022-12-06T00:00:00Z",
    "EndTimestamp": "2022-12-08T00:00:00Z",
    "Price": 5576,
    "HasVollkasko": false,
    "FreeKilometers": 414
  }*/

type JsonOffer struct {
	OfferID              string `json:"OfferID"`
	RegionID             int    `json:"RegionID"`
	CarType              string `json:"CarType"`
	NumberDays           int    `json:"NumberDays"`
	NumberSeats          int    `json:"NumberSeats"`
	StartTimestamp       string `json:"StartTimestamp"`
	EndTimestamp         string `json:"EndTimestamp"`
	Price                int    `json:"Price"`
	HasVollkasko         bool   `json:"HasVollkasko"`
	FreeKilometers       int    `json:"FreeKilometers"`
	MostSpecificRegionID int
}

func TestSeedDatabase(t *testing.T) {
	Init()

	db.Exec("DELETE FROM offers")

	data, err := os.ReadFile("../../testdata/testdata.json")
	if err != nil {
		panic(err)
	}

	var offers []JsonOffer

	err = json.Unmarshal(data, &offers)

	if err != nil {
		panic(err)
	}

	parsedOffers := make([]types.Offer, 0)

	for _, offer := range offers {

		start, err := time.Parse(time.RFC3339, offer.StartTimestamp)

		if err != nil {
			panic(err)
		}

		end, err := time.Parse(time.RFC3339, offer.EndTimestamp)

		if err != nil {
			panic(err)
		}

		parsedOffers = append(parsedOffers, types.Offer{
			ID:                   offer.OfferID,
			Data:                 "",
			MostSpecificRegionID: offer.MostSpecificRegionID,
			StartDate:            start.Unix(),
			EndDate:              end.Unix(),
			NumberSeats:          offer.NumberSeats,
			Price:                int32(offer.Price),
			CarType:              offer.CarType,
			HasVollkasko:         offer.HasVollkasko,
			FreeKilometers:       offer.FreeKilometers,
		})

	}

	err = InsertOffers(parsedOffers)

	if err != nil {
		t.Error(err)
	}

	db.Close()
}
