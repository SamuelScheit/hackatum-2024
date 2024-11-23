package main

import (
	"encoding/json"
	"fmt"

	"checkmate/database"
	"checkmate/types"

	"github.com/valyala/fasthttp"
)

type OffersRequest struct {
	Offers []types.Offer `json:"offers" validate:"required"`
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	var req OffersRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	if len(req.Offers) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("No offers provided")
		return
	}

	for i, offer := range req.Offers {
		if offer.ID == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString(fmt.Sprintf("Offer at index %d has an empty ID", i))
			return
		}
		if offer.MostSpecificRegionID <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString(fmt.Sprintf("Offer at index %d has an invalid MostSpecificRegionID", i))
			return
		}
		if offer.StartDate <= 0 || offer.EndDate <= 0 || offer.StartDate >= offer.EndDate {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString(fmt.Sprintf("Offer at index %d has invalid date ranges", i))
			return
		}
		if offer.Price <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString(fmt.Sprintf("Offer at index %d has an invalid price", i))
			return
		}
		if offer.CarType == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString(fmt.Sprintf("Offer at index %d has an empty CarType", i))
			return
		}
	}

	// Insert each validated offer into the database
	for _, offer := range req.Offers {
		fmt.Println("Inserting offer", offer)
		err := database.InsertSingleOffer(offer)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(fmt.Sprintf("Failed to insert offer ID %s: %v", offer.ID, err))
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("OK"))
}
