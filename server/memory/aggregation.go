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

	// numberSeats := LogicalAnd(carTypeInital, priceRangeInital)
	// LogicalAndInPlace(numberSeats, freeKilometersInital)
	// LogicalAndInPlace(numberSeats, vollkaskoInital)
	numberSeats := priceRangeInital

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

	freeKilometerRanges := []*types.FreeKilometerRange{}
	minFreeKilometerWidth := int32(opts.MinFreeKilometerWidth)

	freeKilometersStart := MinKilometer
	if opts.MinFreeKilometer.Valid {
		freeKilometersStart = opts.MinFreeKilometer.Int32
	}
	freeKilometersStart = int32(math.Floor(float64(freeKilometersStart)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

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

	priceRanges := []*types.PriceRange{}
	priceRangeWidth := int32(opts.PriceRangeWidth)
	priceRangeStart := MinPrice
	if opts.MinPrice.Valid {
		priceRangeStart = opts.MinPrice.Int32
	}
	priceRangeStart = int32(math.Floor(float64(priceRangeStart)/float64(priceRangeWidth))) * priceRangeWidth

	for i := priceRangeStart; i <= MaxPrice; i += priceRangeWidth {
		priceStart := PriceTree.BitArrayGreaterEqual(i)
		priceEnd := PriceTree.BitArrayLessThan(i + priceRangeWidth)
		priceRange := LogicalAnd(priceRangeFiltered, priceStart)
		LogicalAndInPlace(priceRange, priceEnd)
		count := int32(priceRange.CountSetBits())
		if count == 0 {
			continue
		}

		rang := &types.PriceRange{
			Start: i,
			End:   i + priceRangeWidth,
			Count: count,
		}
		priceRanges = append(priceRanges, rang)
	}

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
