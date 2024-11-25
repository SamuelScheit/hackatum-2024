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

	freeKilometersStart := int32(math.Floor(float64(MinKilometer)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

	for i := freeKilometersStart; i <= MaxKilometer; i += minFreeKilometerWidth {
		kilometersStart := KilometerTree.BitArrayGreaterEqual(i)
		kilometersEnd := KilometerTree.BitArrayLessThan(i + minFreeKilometerWidth)
		kilometers := LogicalAnd(kilometerFiltered, kilometersStart)
		LogicalAndInPlace(kilometers, kilometersEnd)
		count := int32(kilometers.CountSetBits())
		if count == 0 {
			continue
		}

		kilometer := &types.FreeKilometerRange{
			Start: i,
			End:   i + minFreeKilometerWidth,
			Count: count,
		}
		freeKilometerRanges = append(freeKilometerRanges, kilometer)
	}

	var rang *types.PriceRange
	priceRange := map[int32]*types.PriceRange{}
	priceRanges := []*types.PriceRange{}
	priceRangeWidth := int32(opts.PriceRangeWidth)
	// priceRangeStart := int32(math.Floor(float64(MinPrice)/float64(priceRangeWidth))) * priceRangeWidth

	var offer *types.Offer
	var ok bool
	var kilometer *types.FreeKilometerRange

	_ = kilometer
	_ = freeKilometerRange

	for i := 0; i <= int(IIDCounter); i++ {
		offer = OfferMap[i]

		if bit, _ := priceRangeFiltered.GetBit(i); bit == 1 {

			price := int32(math.Floor(float64(offer.Price)/float64(priceRangeWidth))) * (priceRangeWidth)

			if rang, ok = priceRange[price]; !ok {
				rang = &types.PriceRange{
					Start: (price),
					End:   (price + priceRangeWidth),
					Count: 0,
				}
				priceRange[price] = rang
				priceRanges = append(priceRanges, rang)
			}

			rang.Count++

		}

		// if bit, _ := kilometerFiltered.GetBit(i); bit == 1 {

		// 	freeKilometers := int32(math.Floor(float64(offer.FreeKilometers)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

		// 	if kilometer, ok = freeKilometerRange[freeKilometers]; !ok {
		// 		kilometer = &types.FreeKilometerRange{
		// 			Start: (freeKilometers),
		// 			End:   (freeKilometers + minFreeKilometerWidth),
		// 			Count: 0,
		// 		}
		// 		freeKilometerRange[freeKilometers] = kilometer
		// 		freeKilometerRanges = append(freeKilometerRanges, kilometer)
		// 	}

		// 	kilometer.Count++
		// }
	}

	// Takes too long, because loop is too large if width is small

	// for i := priceRangeStart; i <= MaxPrice; i += priceRangeWidth {
	// 	priceStart := PriceTree.BitArrayGreaterEqual(i)
	// 	priceEnd := PriceTree.BitArrayLessThan(i + priceRangeWidth)
	// 	priceRange := LogicalAnd(priceRangeFiltered, priceStart)
	// 	LogicalAndInPlace(priceRange, priceEnd)
	// 	count := int32(priceRange.CountSetBits())
	// 	if count == 0 {
	// 		continue
	// 	}

	// 	rang := &types.PriceRange{
	// 		Start: i,
	// 		End:   i + priceRangeWidth,
	// 		Count: count,
	// 	}
	// 	priceRanges = append(priceRanges, rang)
	// }

	seatsCounts := []*types.SeatsCount{}

	for i := MinSeats; i <= MaxSeats; i++ {
		seats, err := GetExactNumberOfSeatsIndex(int(i))
		if err != nil {
			fmt.Println("not enough seats indexes:", i, err)
			panic(err)
		}
		seats = LogicalAnd(numberSeats, seats)
		count := seats.CountSetBits()
		if count == 0 {
			continue
		}

		seat := &types.SeatsCount{
			NumberSeats: int(i),
			Count:       count,
		}
		seatsCounts = append(seatsCounts, seat)
	}

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
