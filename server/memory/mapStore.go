package memory

import (
	"checkmate/types"
	"sync"
)

// takes: iid,  returns: SearchResultOffer
var OfferSearchResultMap map[int32]types.SearchResultOffer
var OfferMap map[int32]*types.Offer

// takes: uuid, returns: iid
var IIDMap map[string]int32
var IIDCounter int32

func InitMapStore() {
	OfferMap = make(map[int32]*types.Offer)
	OfferSearchResultMap = make(map[int32]types.SearchResultOffer)
	IIDMap = make(map[string]int32)

	IIDCounter = 0
}

var mu sync.Mutex

func InsertOffers(offers *[]types.Offer) error {
	mu.Lock()
	defer mu.Unlock()

	for _, offer := range *offers {
		InsertOffer(&offer)
	}
	return nil
}

func InsertOffer(offer *types.Offer) {
	if offer.IID == 0 {
		IIDCounter++
		offer.IID = IIDCounter
	}

	IIDMap[offer.ID] = offer.IID
	OfferSearchResultMap[offer.IID] = types.SearchResultOffer{
		ID:   offer.ID,
		Data: offer.Data,
	}
	OfferMap[offer.IID] = offer
	PriceTree.Add(offer.Price, offer.IID)

	indexKillometer(offer)
	indexVollkasko(offer)
	indexCarType(offer)
	indexNumSeats(offer)
	indexDays(offer)
	indexRegion(offer)
	indexStartDate(offer)
	indexEndDate(offer)
}
