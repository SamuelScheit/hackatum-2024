package memory

import (
	"checkmate/types"
	"time"
)

// indices
var DEFAULT_BITLENGTHSIZE = 10000

// vollkasko
var vollkaskoIndex BitArray

// cartype
var familyCarIndex BitArray
var luxuryCarIndex BitArray
var sportsCarIndex BitArray
var smallCarIndex BitArray

// numSeats
var exactlyOneSeatCarIndex BitArray
var exactlyTwoSeatCarIndex BitArray
var exactlyThreeSeatCarIndex BitArray
var exactlyFourSeatCarIndex BitArray
var exactlyFiveSeatCarIndex BitArray
var exactlySixSeatCarIndex BitArray
var exactlySevenSeatCarIndex BitArray
var exactlyEightSeatCarIndex BitArray

// days
var daysIndexMap []BitArray
var startTree *LinkedBtree
var endTree *LinkedBtree

var priceTree *LinkedBtree
var regionTree *LinkedBtree
var kilometerTree *LinkedBtree

func InitIndex() {
	priceTree = NewLinkedBtree()
	regionTree = NewLinkedBtree()
	kilometerTree = NewLinkedBtree()

	startTree = NewLinkedBtree()
	endTree = NewLinkedBtree()

	daysIndexMap = make([]BitArray, 100)

	for i := 0; i < 100; i++ {
		daysIndexMap[i] = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	}

	vollkaskoIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)

	familyCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	luxuryCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	sportsCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	smallCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)

	exactlyOneSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlyTwoSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlyThreeSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlyFourSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlyFiveSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlySixSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlySevenSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	exactlyEightSeatCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)

}

func indexVollkasko(offer *types.Offer) {
	if offer.HasVollkasko {
		vollkaskoIndex.SetBit(int(offer.IID))
	}
}

func indexCarType(offer *types.Offer) {
	switch offer.CarType {
	case "family":
		familyCarIndex.SetBit(int(offer.IID))
	case "sports":
		sportsCarIndex.SetBit(int(offer.IID))
	case "luxury":
		luxuryCarIndex.SetBit(int(offer.IID))
	case "small":
		smallCarIndex.SetBit(int(offer.IID))
	}
}

func indexNumSeats(offer *types.Offer) {
	switch offer.NumberSeats {
	case 1:
		exactlyOneSeatCarIndex.SetBit(int(offer.IID))
	case 2:
		exactlyTwoSeatCarIndex.SetBit(int(offer.IID))
	case 3:
		exactlyThreeSeatCarIndex.SetBit(int(offer.IID))
	case 4:
		exactlyFourSeatCarIndex.SetBit(int(offer.IID))
	case 5:
		exactlyFiveSeatCarIndex.SetBit(int(offer.IID))
	case 6:
		exactlySixSeatCarIndex.SetBit(int(offer.IID))
	case 7:
		exactlySevenSeatCarIndex.SetBit(int(offer.IID))
	case 8:
		exactlyEightSeatCarIndex.SetBit(int(offer.IID))
	}
}

func indexDays(offer *types.Offer) {
	amountDays := millisecondsToDays(offer.EndDate - offer.StartDate)

	daysIndexMap[amountDays].SetBit(int(offer.IID))

}

func indexRegion(offer *types.Offer) {
	// regionTree.Add(offer.IID, int32(offer.MostSpecificRegionID))
	regionTree.Add(int32(offer.MostSpecificRegionID), offer.IID)
}

func indexStartDate(offer *types.Offer) {
	days := millisecondsToDays(offer.StartDate)
	startTree.Add(days, offer.IID)
}

func indexEndDate(offer *types.Offer) {
	days := millisecondsToDays(offer.StartDate)
	endTree.Add(days, offer.IID)
}

func millisecondsToDays(milliseconds int64) int32 {
	// Convert milliseconds to seconds
	seconds := milliseconds / 1000

	// Create a time.Time object from the Unix timestamp
	t := time.Unix(seconds, 0)

	// Calculate the number of days since the Unix epoch
	days := t.Unix() / (24 * 60 * 60)

	return int32(days)
}
