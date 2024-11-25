package memory

import (
	"checkmate/optimization"
	"checkmate/types"
	"encoding/json"
	"fmt"
	"sort"
)

func QuerySearchResults(opts *types.GetParams) (*types.QueryResponse, error) {

	// -- -- -- required -- -- --

	// region Checking
	result := whereRegionBoundsMatch(opts.RegionID)
	temp := NewBitArray(result.size)

	// check startDate
	daysStart := MillisecondsToDays(opts.TimeRangeStart)
	StartTree.BitArrayGreaterEqual(daysStart, temp)
	LogicalAndInPlace(result, temp)

	// offer.start >= request.startDate
	// request.startDate <= offer.start

	temp.Clear()

	// check EndDate
	daysEnd := MillisecondsToDays(opts.TimeRangeEnd)
	EndTree.BitArrayLessEqual(daysEnd, temp)
	LogicalAndInPlace(result, temp)

	// check daysAmount
	amountDays := opts.NumberDays

	if int(amountDays) > len(DaysIndexMap) {
		// TODO: dynamically calculate the days map
		return nil, fmt.Errorf("amount of days is too high")
	}

	LogicalAndInPlace(result, &DaysIndexMap[amountDays])

	// -- -- -- optional -- -- --

	vollkaskoInital := result
	priceRangeInital := result
	carTypeInital := result
	numberSeatsInital := result
	freeKilometersInital := result

	// CarType
	if opts.CarType.Valid {
		carTypeInital = LogicalAnd(result, GetCarTypeIndex(opts.CarType.String))
	}

	// MaxPrice (exclusive)
	if opts.MaxPrice.Valid {
		priceRangeInital = NewBitArray(result.size)
		PriceTree.BitArrayLessThan(opts.MaxPrice.Int32, priceRangeInital)
	}

	// MinPrice (inclusive)
	if opts.MinPrice.Valid {
		priceRangeMinInital := NewBitArray(result.size)

		PriceTree.BitArrayGreaterEqual(opts.MinPrice.Int32, priceRangeMinInital)
		LogicalAndInPlace(priceRangeInital, priceRangeMinInital)
	}

	// MinFreeKilometer (inclusive)
	if opts.MinFreeKilometer.Valid {
		freeKilometersInital = NewBitArray(result.size)
		KilometerTree.BitArrayGreaterEqual(opts.MinFreeKilometer.Int32, freeKilometersInital)
	}

	// OnlyVallkasko
	if opts.OnlyVollkasko.Bool && opts.OnlyVollkasko.Valid {
		vollkaskoInital = LogicalAnd(result, whereHasVollkaskoIsTrue())
	}

	// MinNumberSeats (inclusive)
	if opts.MinNumberSeats.Valid && opts.MinNumberSeats.Int32 > 0 {
		seats, err := GetNumberOfSeatsIndex(int(opts.MinNumberSeats.Int32))
		if err != nil {
			return nil, err
		}
		numberSeatsInital = LogicalAnd(result, seats)
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

	// result = pagination(result, opts.Page, opts.PageSize)

	searchResults, err := collectOfferJSONSorted(result, OfferMap, opts.SortOrder == 0)

	searchResults = paginationArray(searchResults, opts.Page, opts.PageSize)

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

func paginationArray[T any](in []T, pageNumber uint, pageSize uint) []T {
	startIndex := pageNumber * pageSize
	endIndex := startIndex + pageSize

	if startIndex > uint(len(in)) {
		return []T{}
	}

	if endIndex > uint(len(in)) {
		endIndex = uint(len(in))
	}

	return in[startIndex:endIndex]
}

func whereRegionBoundsMatch(regionID uint) *BitArray {
	min, max, min2, max2 := optimization.GetRegionBounds(regionID)

	temp1 := NewBitArray(RegionTree.Size)
	temp2 := NewBitArray(RegionTree.Size)
	temp3 := NewBitArray(RegionTree.Size)

	RegionTree.BitArrayGreaterEqual(int32(min), temp1) // temp1 = x > min
	RegionTree.BitArrayLessEqual(int32(max), temp2)    // temp2 = x < max
	LogicalAndInPlace(temp1, temp2)                    // temp1 = (x > min) AND (x < max)

	temp2.Clear()

	RegionTree.BitArrayGreaterEqual(int32(min2), temp2) // temp2 = x > min2
	RegionTree.BitArrayLessEqual(int32(max2), temp3)    // temp3 = x < max2
	LogicalAndInPlace(temp2, temp3)                     // temp2 = (x > min2) AND (x < max2)

	// Combine both ranges with OR
	LogicalOrInPlace(temp1, temp2) // temp1 = ((x > min) AND (x < max)) OR ((x > min2) AND (x < max2))

	// Return the final result
	return temp1
}

func whereHasVollkaskoIsTrue() *BitArray {
	return &VollkaskoIndex
}

func GetNumberOfSeatsIndex(amount int) (*BitArray, error) {
	if amount >= len(MinSeatIndexMap) {
		return nil, fmt.Errorf("amount of seats is too high")
	}
	return &MinSeatIndexMap[amount], nil
}

func GetCarTypeIndex(cartype string) *BitArray {
	switch cartype {
	case "family":
		return &FamilyCarIndex
	case "sports":
		return &SportsCarIndex
	case "luxury":
		return &LuxuryCarIndex
	case "small":
		return &SmallCarIndex
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
			} else {
				fmt.Println("Offer not found for IID", i)
				return nil, fmt.Errorf("Offer not found for IID %d", i)
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
