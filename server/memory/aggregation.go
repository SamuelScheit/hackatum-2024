package memory

import (
	"checkmate/types"
	"math"
	"slices"
)

func min(one int, two int) int {
	if one < two {
		return one
	}
	return two
}

func getPriceRangeAggregation(opts *types.GetParams,
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

	// opts.PriceRangeWidth
	for i := 0; i <= int(IIDCounter); i++ {

		offer = OfferMap[i]

		if bit, _ := kilometerFiltered.GetBit(i); bit == 1 {

			freeKilometers := int32(math.Floor(float64(offer.FreeKilometers)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

			if kilometer, ok = freeKilometerRange[freeKilometers]; !ok {
				kilometer = &types.FreeKilometerRange{
					Start: int(freeKilometers),
					End:   int(freeKilometers + minFreeKilometerWidth),
					Count: 0,
				}
				freeKilometerRange[freeKilometers] = kilometer
				freeKilometerRanges = append(freeKilometerRanges, kilometer)
			}

			kilometer.Count++
		}

		if bit, _ := numberSeats.GetBit(i); bit == 1 {

			seats := int32(offer.NumberSeats)

			if seat, ok = seatsCount[seats]; !ok {
				seat = &types.SeatsCount{
					NumberSeats: int(seats),
					Count:       0,
				}
				seatsCount[seats] = seat
				seatsCounts = append(seatsCounts, seat)
			}

			seat.Count++
		}

		if bit, _ := priceRangeFiltered.GetBit(i); bit == 1 {

			price := int32(math.Floor(float64(offer.Price)/float64(priceRangeWidth))) * (priceRangeWidth)

			if rang, ok = priceRange[price]; !ok {
				rang = &types.PriceRange{
					Start: int(price),
					End:   int(price + priceRangeWidth),
					Count: 0,
				}
				priceRange[price] = rang
				priceRanges = append(priceRanges, rang)
			}

			rang.Count++

		}
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
