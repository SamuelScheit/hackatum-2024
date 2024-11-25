package memory

import (
	"checkmate/types"
	"fmt"
	"math"
	"slices"
)

func min(one int32, two int32) int32 {
	if one < two {
		return one
	}
	return two
}

func max(one int32, two int32) int32 {
	if one > two {
		return one
	}
	return two
}

func getAggregation(opts *types.GetParams,
	priceRangeInital *BitArray,
	carTypeInital *BitArray,
	numberSeatsInital *BitArray,
	freeKilometersInital *BitArray,
	vollkaskoInital *BitArray,
) *types.QueryResponse {

	var rang *types.PriceRange
	priceRange := map[int32]*types.PriceRange{}
	priceRanges := []*types.PriceRange{}
	priceRangeWidth := int32(opts.PriceRangeWidth)

	var offer *types.Offer
	var ok bool

	priceRangeFiltered := LogicalAnd(numberSeatsInital, carTypeInital)
	LogicalAndInPlace(priceRangeFiltered, freeKilometersInital)
	LogicalAndInPlace(priceRangeFiltered, vollkaskoInital)

	carTypeFiltered := LogicalAnd(numberSeatsInital, priceRangeInital)
	LogicalAndInPlace(carTypeFiltered, freeKilometersInital)
	LogicalAndInPlace(carTypeFiltered, vollkaskoInital)

	numberSeats := LogicalAnd(carTypeInital, priceRangeInital)
	LogicalAndInPlace(numberSeats, freeKilometersInital)
	LogicalAndInPlace(numberSeats, vollkaskoInital)

	kilometerFiltered := LogicalAnd(carTypeInital, priceRangeInital)
	LogicalAndInPlace(kilometerFiltered, numberSeatsInital)
	LogicalAndInPlace(kilometerFiltered, vollkaskoInital)

	vollkaskoFiltered := LogicalAnd(carTypeInital, priceRangeInital)
	LogicalAndInPlace(vollkaskoFiltered, numberSeatsInital)
	LogicalAndInPlace(vollkaskoFiltered, freeKilometersInital)

	carTypeCounts := types.CarTypeCount{
		Small:  LogicalAnd(carTypeFiltered, GetCarTypeIndex("small")).CountSetBits(),
		Sports: LogicalAnd(carTypeFiltered, GetCarTypeIndex("sports")).CountSetBits(),
		Luxury: LogicalAnd(carTypeFiltered, GetCarTypeIndex("luxury")).CountSetBits(),
		Family: LogicalAnd(carTypeFiltered, GetCarTypeIndex("family")).CountSetBits(),
	}

	vollkaskoCounts := types.VollkaskoCount{
		TrueCount:  LogicalAnd(vollkaskoFiltered, &VollkaskoIndex).CountSetBits(),
		FalseCount: LogicalAnd(vollkaskoFiltered, &NoVollkaskoIndex).CountSetBits(),
	}

	freeKilometerRange := map[int32]*types.FreeKilometerRange{}
	freeKilometerRanges := []*types.FreeKilometerRange{}
	minFreeKilometerWidth := int32(opts.MinFreeKilometerWidth)

	var kilometer *types.FreeKilometerRange

	seatsCount := map[int32]*types.SeatsCount{}
	seatsCounts := []*types.SeatsCount{}
	var seat *types.SeatsCount

	freeKilometersStart := MinKilometer

	if opts.MinFreeKilometer.Valid {
		freeKilometersStart = opts.MinFreeKilometer.Int32
	}

	freeKilometersStart = int32(math.Floor(float64(freeKilometersStart)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

	for i := freeKilometersStart; i <= MaxKilometer; i += minFreeKilometerWidth {
		kilometersStart := KilometerTree.BitArrayGreaterEqual(i, nil)
		kilometersEnd := KilometerTree.BitArrayLessEqual(i+minFreeKilometerWidth, nil)
		LogicalAndInPlace(kilometerFiltered, kilometersStart)
		LogicalAndInPlace(kilometerFiltered, kilometersEnd)
		count := int32(kilometerFiltered.CountSetBits())
		if count == 0 {
			continue
		}

		kilometer = &types.FreeKilometerRange{
			Start: i,
			End:   i + minFreeKilometerWidth,
			Count: count,
		}
		freeKilometerRanges = append(freeKilometerRanges, kilometer)
	}

	priceRangeStart := MinPrice

	if opts.MinPrice.Valid {
		priceRangeStart = opts.MinPrice.Int32
	}

	priceRangeStart = int32(math.Floor(float64(priceRangeStart)/float64(priceRangeWidth))) * priceRangeWidth

	for i := priceRangeStart; i <= MaxPrice; i += priceRangeWidth {
		priceRangeStart := PriceTree.BitArrayGreaterEqual(i, nil)
		priceRangeEnd := PriceTree.BitArrayLessEqual(i+priceRangeWidth, nil)
		LogicalAndInPlace(priceRangeFiltered, priceRangeStart)
		LogicalAndInPlace(priceRangeFiltered, priceRangeEnd)
		count := int32(priceRangeFiltered.CountSetBits())
		if count == 0 {
			continue
		}

		rang = &types.PriceRange{
			Start: i,
			End:   i + priceRangeWidth,
			Count: count,
		}
		priceRanges = append(priceRanges, rang)
	}

	for i := MinSeats; i <= MaxSeats; i++ {
		seats, err := GetMinNumberOfSeatsIndex(int(i))
		if err != nil {
			fmt.Println("not enough seats indexes:", i, err)
			panic(err)
		}
		LogicalAndInPlace(numberSeats, seats)
		count := numberSeats.CountSetBits()
		if count == 0 {
			continue
		}

		seat = &types.SeatsCount{
			NumberSeats: int(i),
			Count:       count,
		}
		seatsCounts = append(seatsCounts, seat)
	}

	_ = seatsCount
	_ = carTypeCounts
	_ = freeKilometerRange
	_ = vollkaskoCounts
	_ = priceRange
	_ = offer
	_ = kilometerFiltered
	_ = priceRangeFiltered
	_ = carTypeFiltered
	_ = numberSeats
	_ = vollkaskoFiltered
	_ = minFreeKilometerWidth
	_ = priceRangeWidth
	_ = carTypeInital
	_ = freeKilometersInital
	_ = numberSeatsInital
	_ = vollkaskoInital
	_ = kilometer
	_ = rang
	_ = seat
	_ = ok

	slices.SortFunc(priceRanges, func(a, b *types.PriceRange) int {
		return int(a.Start - b.Start)
	})

	slices.SortFunc(seatsCounts, func(a, b *types.SeatsCount) int {
		return int(a.NumberSeats - b.NumberSeats)
	})

	slices.SortFunc(freeKilometerRanges, func(a, b *types.FreeKilometerRange) int {
		return int(a.Start - b.Start)
	})

	response := types.QueryResponse{
		Offers:             []types.SearchResultOffer{},
		PriceRanges:        priceRanges,
		CarTypeCounts:      carTypeCounts,
		SeatsCount:         seatsCounts,
		FreeKilometerRange: freeKilometerRanges,
		VollkaskoCount:     vollkaskoCounts,
	}

	return &response
}
