package memory

import (
	"checkmate/optimization"
	"checkmate/types"
	"encoding/json"
	"fmt"
	"sort"
)

func QuerySearchResults(opts *types.GetParams) (*types.QueryResponse, error) {

	result := NewBitArray(DEFAULT_BITLENGTHSIZE)
	temp := NewBitArray(DEFAULT_BITLENGTHSIZE)

	// -- -- -- required -- -- --

	// region Checking
	LogicalOrInPlace(result, whereRegionBoundsMatch(opts.RegionID))

	// check startDate
	daysStart := millisecondsToDays(opts.TimeRangeStart)
	startTree.BitArrayGreaterEqual(daysStart, temp)
	LogicalAndInPlace(result, temp)

	// offer.start >= request.startDate
	// request.startDate <= offer.start

	temp.Clear()

	// check EndDate
	daysEnd := millisecondsToDays(opts.TimeRangeEnd)
	endTree.BitArrayLessThanEqual(daysEnd, temp)
	LogicalAndInPlace(result, temp)

	// check daysAmount
	amountDays := daysEnd - daysStart

	if int(amountDays) > len(daysIndexMap) {
		// TODO: dynamically calculate the days map
		return nil, fmt.Errorf("amount of days is too high")
	}

	LogicalAndInPlace(result, &daysIndexMap[amountDays])

	// -- -- -- optional -- -- --

	vollkaskoInital := result
	priceRangeInital := result
	carTypeInital := result
	numberSeatsInital := result
	freeKilometersInital := result

	// CarType
	if opts.CarType.Valid {
		carTypeInital = LogicalAnd(result, whereCarTypIs(opts.CarType.String))
	}

	if opts.MinPrice.Valid || opts.MaxPrice.Valid {
		priceRangeInital = result.Copy()
	}

	// MaxPrice
	if opts.MaxPrice.Valid {
		priceTree.BitArrayLessThan(opts.MaxPrice.Int32, priceRangeInital)
	}

	// MinPrice
	if opts.MinPrice.Valid {
		copy := priceRangeInital.Copy()
		priceTree.BitArrayGreaterEqual(opts.MinPrice.Int32, copy)
		LogicalAndInPlace(priceRangeInital, copy)
	}

	// MinFreeKilometer
	if opts.MinFreeKilometer.Valid {
		freeKilometersInital = result.Copy()
		kilometerTree.BitArrayGreaterEqual(opts.MinFreeKilometer.Int32, freeKilometersInital)
	}

	// OnlyVallkasko
	if opts.OnlyVollkasko.Bool && opts.OnlyVollkasko.Valid {
		vollkaskoInital = LogicalAnd(result, whereHasVollkaskoIsTrue())
	}

	// MinNumberSeats
	if opts.MinNumberSeats.Valid && opts.MinNumberSeats.Int32 > 0 {
		numberSeatsInital = LogicalAnd(result, whereNumberOfSeatsIs(int(opts.MinNumberSeats.Int32)))
	}

	response := getPriceRangeAggregation(opts, priceRangeInital, carTypeInital, numberSeatsInital, freeKilometersInital, vollkaskoInital)

	if opts.CarType.Valid {
		LogicalAndInPlace(result, carTypeInital)
	}

	if opts.MinPrice.Valid || opts.MaxPrice.Valid {
		LogicalAndInPlace(result, priceRangeInital)
	}

	if opts.MinFreeKilometer.Valid {
		LogicalAndInPlace(result, freeKilometersInital)
	}

	if opts.OnlyVollkasko.Bool && opts.OnlyVollkasko.Valid {
		LogicalAndInPlace(result, vollkaskoInital)
	}

	if opts.MinNumberSeats.Valid && opts.MinNumberSeats.Int32 > 0 {
		LogicalAndInPlace(result, numberSeatsInital)
	}

	result = pagination(result, opts.Page, opts.PageSize)

	searchResults, err := collectOfferJSONSorted(result, offerMap, opts.SortOrder == 0)

	if err != nil {
		return nil, err
	}

	response.Offers = searchResults

	return response, nil

}

func pagination(in *BitArray, pageNumber uint, pageSize uint) *BitArray {
	startIndex := pageNumber * pageSize
	endIndex := startIndex + pageSize

	result := NewBitArray(in.size)
	count := uint(0)

	for i := 0; i < in.size; i++ {
		bit, _ := in.GetBit(i)
		if bit == 1 {
			if count >= startIndex && count < endIndex {
				result.SetBit(i)
			}
			count++
			if count >= endIndex {
				break
			}
		}
	}

	return result
}

func whereRegionBoundsMatch(regionID uint) *BitArray {
	min, max, min2, max2 := optimization.GetRegionBounds(regionID)

	temp1 := NewBitArray(regionTree.Size)
	temp2 := NewBitArray(regionTree.Size)
	temp3 := NewBitArray(regionTree.Size)

	regionTree.BitArrayGreaterEqual(int32(min), temp1)  // temp1 = x > min
	regionTree.BitArrayLessThanEqual(int32(max), temp2) // temp2 = x < max
	LogicalAndInPlace(temp1, temp2)                     // temp1 = (x > min) AND (x < max)

	temp2.Clear()

	regionTree.BitArrayGreaterEqual(int32(min2), temp2)  // temp2 = x > min2
	regionTree.BitArrayLessThanEqual(int32(max2), temp3) // temp3 = x < max2
	LogicalAndInPlace(temp2, temp3)                      // temp2 = (x > min2) AND (x < max2)

	// Combine both ranges with OR
	LogicalOrInPlace(temp1, temp2) // temp1 = ((x > min) AND (x < max)) OR ((x > min2) AND (x < max2))

	// Return the final result
	return temp1
}

func whereHasVollkaskoIsTrue() *BitArray {
	return &vollkaskoIndex
}

func whereNumberOfSeatsIs(amount int) *BitArray {
	switch amount {
	case 1:
		return &exactlyOneSeatCarIndex
	case 2:
		return &exactlyTwoSeatCarIndex
	case 3:
		return &exactlyThreeSeatCarIndex
	case 4:
		return &exactlyFourSeatCarIndex
	case 5:
		return &exactlyFiveSeatCarIndex
	case 6:
		return &exactlySixSeatCarIndex
	case 7:
		return &exactlySevenSeatCarIndex
	case 8:
		return &exactlyEightSeatCarIndex
	}
	panic("Invalid number of seats: " + string(amount))
}

func whereCarTypIs(cartype string) *BitArray {
	switch cartype {
	case "family":
		return &familyCarIndex
	case "sports":
		return &sportsCarIndex
	case "luxury":
		return &luxuryCarIndex
	case "small":
		return &smallCarIndex
	}
	panic("invalid Cartype")
}

func collectOfferJSON(ba *BitArray, offerMap map[int32]types.SearchResultOffer) (string, error) {
	var results []types.SearchResultOffer

	for i := 0; i < ba.size; i++ {
		bit, _ := ba.GetBit(i)
		if bit == 1 {
			if offer, exists := offerMap[int32(i)]; exists {
				results = append(results, offer)
			}
		}
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	return string(jsonData), nil
}

func collectOfferJSONSorted(ba *BitArray, offerMap map[int32]*types.Offer, sortAscending bool) ([]types.SearchResultOffer, error) {
	var results []types.Offer

	// Collect offers matching the BitArray
	for i := 0; i < ba.size; i++ {
		bit, _ := ba.GetBit(i)
		if bit == 1 {
			if offer, exists := offerMap[int32(i)]; exists {
				results = append(results, *offer)
			}
		}
	}

	// Sort the results based on the Price and the sorting order
	sort.Slice(results, func(i, j int) bool {
		if results[i].Price == results[j].Price {
			return results[i].ID < results[j].ID
		}
		if sortAscending {
			return results[i].Price < results[j].Price
		}
		return results[i].Price > results[j].Price
	})
	// OptimizedSearchResultOffer

	convertedResults := []types.SearchResultOffer{}
	for _, sortedRow := range results {
		convertedResults = append(convertedResults, types.SearchResultOffer{
			ID:   sortedRow.ID,
			Data: sortedRow.Data,
		})
	}

	return convertedResults, nil
}
