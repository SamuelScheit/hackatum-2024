package memory

import (
	"checkmate/types"
	"sync"
)

// takes: iid,  returns: SearchResultOffer
var OfferSearchResultMap []types.SearchResultOffer
var OfferMap []*types.Offer

// takes: uuid, returns: iid
var IIDMap map[string]int32
var IIDCounter int32

func InitMapStore() {
	OfferMap = make([]*types.Offer, DEFAULT_BITLENGTHSIZE/64)
	OfferSearchResultMap = make([]types.SearchResultOffer, DEFAULT_BITLENGTHSIZE/64)
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
		if IIDCounter >= int32(len(OfferMap)) {
			l := int(IIDCounter) - len(OfferMap) + DEFAULT_BITLENGTHSIZE/64
			OfferMap = append(OfferMap, make([]*types.Offer, l)...)
			OfferSearchResultMap = append(OfferSearchResultMap, make([]types.SearchResultOffer, l)...)
		}
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
