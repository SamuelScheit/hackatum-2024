package memory

import (
	"checkmate/types"
	"time"
)

// indices
var DEFAULT_BITLENGTHSIZE = 10000

// vollkasko
var VollkaskoIndex BitArray
var NoVollkaskoIndex BitArray

// cartype
var FamilyCarIndex BitArray
var LuxuryCarIndex BitArray
var SportsCarIndex BitArray
var SmallCarIndex BitArray

// numSeats
var MinSeatIndexMap []BitArray
var ExactSeatIndexMap []BitArray

// days
var DaysIndexMap []BitArray
var StartTree *LinkedBtree
var EndTree *LinkedBtree

var PriceTree *LinkedBtree
var RegionTree *LinkedBtree
var KilometerTree *LinkedBtree
var MaxKilometer = int32(0)
var MinKilometer = int32(1 << 30)
var MinPrice = int32(1 << 30)
var MaxPrice = int32(0)
var MinSeats = int32(1 << 30)
var MaxSeats = int32(0)

func InitIndex() {
	PriceTree = NewLinkedBtree()
	RegionTree = NewLinkedBtree()
	KilometerTree = NewLinkedBtree()

	StartTree = NewLinkedBtree()
	EndTree = NewLinkedBtree()

	DaysIndexMap = make([]BitArray, 100)

	for i := 0; i < len(DaysIndexMap); i++ {
		DaysIndexMap[i] = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	}

	VollkaskoIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	NoVollkaskoIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)

	FamilyCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	LuxuryCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	SportsCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	SmallCarIndex = *NewBitArray(DEFAULT_BITLENGTHSIZE)

	MinSeatIndexMap = make([]BitArray, 10)
	ExactSeatIndexMap = make([]BitArray, 10)

	for i := 0; i < len(MinSeatIndexMap); i++ {
		MinSeatIndexMap[i] = *NewBitArray(DEFAULT_BITLENGTHSIZE)
		ExactSeatIndexMap[i] = *NewBitArray(DEFAULT_BITLENGTHSIZE)
	}
}

func indexKillometer(offer *types.Offer) {
	KilometerTree.Add(offer.FreeKilometers, offer.IID)
}

func indexVollkasko(offer *types.Offer) {
	if offer.HasVollkasko {
		VollkaskoIndex.SetBit(int(offer.IID))
	} else {
		NoVollkaskoIndex.SetBit(int(offer.IID))
	}
}

func indexCarType(offer *types.Offer) {
	index := GetCarTypeIndex(offer.CarType)
	index.SetBit(int(offer.IID))
}

func indexNumSeats(offer *types.Offer) {
	if int32(offer.NumberSeats) >= int32(len(MinSeatIndexMap)) {
		for i := len(MinSeatIndexMap); i <= int(offer.NumberSeats); i++ {
			MinSeatIndexMap = append(MinSeatIndexMap, *NewBitArray(DEFAULT_BITLENGTHSIZE))
			ExactSeatIndexMap = append(ExactSeatIndexMap, *NewBitArray(DEFAULT_BITLENGTHSIZE))
		}
	}

	for i := int(offer.NumberSeats); i < len(MinSeatIndexMap); i++ {
		MinSeatIndexMap[i].SetBit(int(offer.IID))
	}

	ExactSeatIndexMap[offer.NumberSeats].SetBit(int(offer.IID))

}

func indexDays(offer *types.Offer) {
	amountDays := MillisecondsToDays(offer.EndDate - offer.StartDate)

	DaysIndexMap[amountDays].SetBit(int(offer.IID))
}

func indexRegion(offer *types.Offer) {
	RegionTree.Add(int32(offer.MostSpecificRegionID), offer.IID)
}

func indexStartDate(offer *types.Offer) {
	days := MillisecondsToDays(offer.StartDate)
	StartTree.Add(days, offer.IID)
}

func indexEndDate(offer *types.Offer) {
	days := MillisecondsToDays(offer.EndDate)
	EndTree.Add(days, offer.IID)
}

func MillisecondsToDays(milliseconds int64) int32 {
	// Convert milliseconds to seconds
	seconds := milliseconds / 1000

	// Create a time.Time object from the Unix timestamp
	t := time.Unix(seconds, 0)

	// Calculate the number of days since the Unix epoch
	days := t.Unix() / (24 * 60 * 60)

	return int32(days)
}
