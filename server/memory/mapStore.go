package memory

import (
	"checkmate/types"
)

// takes: iid,  returns: SearchResultOffer
var offerSearchResultMap map[int32]types.SearchResultOffer
var offerMap map[int32]*types.Offer

// takes: uuid, returns: iid
var iidMap map[string]int32
var iidCounter int32

func InitMapStore() {
	offerMap = make(map[int32]*types.Offer)
	offerSearchResultMap = make(map[int32]types.SearchResultOffer)
	iidMap = make(map[string]int32)

	iidCounter = 0
}

func InsertOffers(offers *[]types.Offer) error {
	for _, offer := range *offers {
		InsertOffer(&offer)
	}
	return nil
}

func InsertOffer(offer *types.Offer) {
	if offer.IID == 0 {
		iidCounter++
		offer.IID = iidCounter
	}

	iidMap[offer.ID] = offer.IID
	offerSearchResultMap[offer.IID] = types.SearchResultOffer{
		ID:   offer.ID,
		Data: offer.Data,
	}
	offerMap[offer.IID] = offer
	priceTree.Add(offer.Price, offer.IID)

	indexVollkasko(offer)
	indexCarType(offer)
	indexNumSeats(offer)
	indexDays(offer)
	indexRegion(offer)
	indexStartDate(offer)
	indexEndDate(offer)
}
