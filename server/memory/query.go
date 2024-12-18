package memory

import (
	"checkmate/optimization"
	"checkmate/types"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

func QuerySearchResults(opts *types.GetParams) (*types.QueryResponse, error) {

	// -- -- -- required -- -- --

	start := time.Now()

	// region Checking
	result := whereRegionBoundsMatch(opts.RegionID)

	fmt.Println("Region Bounds Match took: ", time.Since(start))

	// check startDate
	daysStart := MillisecondsToDays(opts.TimeRangeStart)
	startResult := StartTree.BitArrayGreaterEqual(daysStart)
	LogicalAndInPlace(result, startResult)

	// offer.start >= request.startDate
	// request.startDate <= offer.start

	// check EndDate
	daysEnd := MillisecondsToDays(opts.TimeRangeEnd)
	endResult := EndTree.BitArrayLessEqual(daysEnd)
	LogicalAndInPlace(result, endResult)

	// check daysAmount
	amountDays := opts.NumberDays

	if int(amountDays) > len(DaysIndexMap) {
		// TODO: dynamically calculate the days map
		return nil, fmt.Errorf("amount of days is too high")
	}

	LogicalAndInPlace(result, &DaysIndexMap[amountDays])

	fmt.Println("required filters took: ", time.Since(start))

	// -- -- -- optional -- -- --

	vollkaskoInital := result
	priceRangeInital := result
	carTypeInital := result
	numberSeatsInital := result
	freeKilometersInital := result

	start = time.Now()

	// CarType
	if opts.CarType.Valid {
		carTypeInital = LogicalAnd(result, GetCarTypeIndex(opts.CarType.String))
	}

	// MaxPrice (exclusive)
	if opts.MaxPrice.Valid {
		priceRangeInital = PriceTree.BitArrayLessThan(opts.MaxPrice.Int32)

	}

	// MinPrice (inclusive)
	if opts.MinPrice.Valid {
		priceRangeMinInital := PriceTree.BitArrayGreaterEqual(opts.MinPrice.Int32)
		priceRangeInital = LogicalAnd(priceRangeInital, priceRangeMinInital)
	}

	// MinFreeKilometer (inclusive)
	if opts.MinFreeKilometer.Valid {
		freeKilometersInital = KilometerTree.BitArrayGreaterEqual(opts.MinFreeKilometer.Int32)
	}

	// OnlyVallkasko
	if opts.OnlyVollkasko.Bool && opts.OnlyVollkasko.Valid {
		vollkaskoInital = LogicalAnd(result, whereHasVollkaskoIsTrue())
	}

	// MinNumberSeats (inclusive)
	if opts.MinNumberSeats.Valid && opts.MinNumberSeats.Int32 > 0 {
		seats, err := GetMinNumberOfSeatsIndex(int(opts.MinNumberSeats.Int32))
		if err != nil {
			return nil, err
		}
		numberSeatsInital = LogicalAnd(result, seats)
	}

	fmt.Println("optional filters took: ", time.Since(start))

	start = time.Now()

	response := getAggregation(opts, priceRangeInital, carTypeInital, numberSeatsInital, freeKilometersInital, vollkaskoInital)

	fmt.Println("aggregation took: ", time.Since(start))

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

	start = time.Now()

	searchResults, err := collectOfferJSONSorted(result, OfferMap, opts.SortOrder == 0)

	searchResults = paginationArray(searchResults, opts.Page, opts.PageSize)

	if err != nil {
		return nil, err
	}

	response.Offers = searchResults

	fmt.Println("collecting offers took: ", time.Since(start))

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

	temp1 := RegionTree.BitArrayGreaterEqual(int32(min)) // temp1 = x > min
	temp2 := RegionTree.BitArrayLessEqual(int32(max))    // temp2 = x < max
	temp1 = LogicalAnd(temp1, temp2)                     // temp1 = (x > min) AND (x < max)

	temp2 = RegionTree.BitArrayGreaterEqual(int32(min2)) // temp2 = x > min2
	temp3 := RegionTree.BitArrayLessEqual(int32(max2))   // temp3 = x < max2
	temp2 = LogicalAnd(temp2, temp3)                     // temp2 = (x > min2) AND (x < max2)

	// Combine both ranges with OR
	return LogicalOr(temp1, temp2) // temp1 = ((x > min) AND (x < max)) OR ((x > min2) AND (x < max2))
}

func whereHasVollkaskoIsTrue() *BitArray {
	return &VollkaskoIndex
}

func GetMinNumberOfSeatsIndex(amount int) (*BitArray, error) {
	if amount >= len(MinSeatIndexMap) {
		return nil, fmt.Errorf("amount of seats is too high")
	}
	return &MinSeatIndexMap[amount], nil
}

func GetExactNumberOfSeatsIndex(amount int) (*BitArray, error) {
	if amount >= len(MinSeatIndexMap) {
		return nil, fmt.Errorf("amount of seats is too high")
	}
	return &ExactSeatIndexMap[amount], nil
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

func collectOfferJSONSorted(ba *BitArray, offerMap []*types.Offer, sortAscending bool) ([]types.SearchResultOffer, error) {
	var results []types.Offer

	// Collect offers matching the BitArray
	for i := 0; i < ba.size; i++ {
		bit, _ := ba.GetBit(i)
		if bit == 1 {
			offer := offerMap[i]
			results = append(results, *offer)
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
