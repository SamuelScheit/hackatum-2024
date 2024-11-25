package memory

import (
	"checkmate/types"
	"slices"
)

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
		Small:  0,
		Sports: 0,
		Luxury: 0,
		Family: 0,
	}

	freeKilometerRange := map[int32]*types.FreeKilometerRange{}
	freeKilometerRanges := []*types.FreeKilometerRange{}
	minFreeKilometerWidth := int32(opts.MinFreeKilometerWidth)

	var kilometer *types.FreeKilometerRange

	vollkaskoCounts := types.VollkaskoCount{
		TrueCount:  0,
		FalseCount: 0,
	}

	seatsCount := map[int32]*types.SeatsCount{}
	seatsCounts := []*types.SeatsCount{}
	var seat *types.SeatsCount

	vollkaskoCounts.TrueCount = vollkaskoFiltered.CountSetBits()
	vollkaskoCounts.FalseCount = priceRangeFiltered.size - vollkaskoCounts.TrueCount

	// opts.PriceRangeWidth
	// for i := 0; i < priceRangeFiltered.size; i++ {

	// 	if offer, ok = OfferMap[int32(i)]; !ok {
	// 		continue
	// 	}

	// 	if bit, _ := vollkaskoFiltered.GetBit(i); bit == 1 {

	// 		if offer.HasVollkasko {
	// 			vollkaskoCounts.TrueCount++
	// 		} else {
	// 			vollkaskoCounts.FalseCount++
	// 		}
	// 	}

	// 	if bit, _ := kilometerFiltered.GetBit(i); bit == 1 {

	// 		freeKilometers := int32(math.Floor(float64(offer.FreeKilometers)/float64(minFreeKilometerWidth))) * minFreeKilometerWidth

	// 		if kilometer, ok = freeKilometerRange[freeKilometers]; !ok {
	// 			kilometer = &types.FreeKilometerRange{
	// 				Start: int(freeKilometers),
	// 				End:   int(freeKilometers + minFreeKilometerWidth),
	// 				Count: 0,
	// 			}
	// 			freeKilometerRange[freeKilometers] = kilometer
	// 			freeKilometerRanges = append(freeKilometerRanges, kilometer)
	// 		}

	// 		kilometer.Count++
	// 	}

	// 	if bit, _ := carTypeFiltered.GetBit(i); bit == 1 {

	// 		// TODO: speedup
	// 		switch offer.CarType {
	// 		case "small":
	// 			carTypeCounts.Small++
	// 		case "sports":
	// 			carTypeCounts.Sports++
	// 		case "luxury":
	// 			carTypeCounts.Luxury++
	// 		case "family":
	// 			carTypeCounts.Family++
	// 		}
	// 	}

	// 	if bit, _ := numberSeats.GetBit(i); bit == 1 {

	// 		seats := int32(offer.NumberSeats)

	// 		if seat, ok = seatsCount[seats]; !ok {
	// 			seat = &types.SeatsCount{
	// 				NumberSeats: int(seats),
	// 				Count:       0,
	// 			}
	// 			seatsCount[seats] = seat
	// 			seatsCounts = append(seatsCounts, seat)
	// 		}

	// 		seat.Count++
	// 	}

	// 	if bit, _ := priceRangeFiltered.GetBit(i); bit == 1 {

	// 		price := int32(math.Floor(float64(offer.Price)/float64(priceRangeWidth))) * (priceRangeWidth)

	// 		if rang, ok = priceRange[price]; !ok {
	// 			rang = &types.PriceRange{
	// 				Start: int(price),
	// 				End:   int(price + priceRangeWidth),
	// 				Count: 0,
	// 			}
	// 			priceRange[price] = rang
	// 			priceRanges = append(priceRanges, rang)
	// 		}

	// 		rang.Count++

	// 	}
	// }

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
