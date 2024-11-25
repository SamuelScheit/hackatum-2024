package main

import (
	"bufio"
	"bytes"
	"checkmate/memory"
	"checkmate/optimization"
	"checkmate/routes"
	"checkmate/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	memory.Init()
	optimization.Init()
	go routes.Serve()

	entries, err := os.ReadDir(".")

	if err != nil {
		panic(err)
	}

	for _, entry := range entries {

		if strings.Contains(entry.Name(), ".log") {

			openLog(entry.Name())
		}
	}
}

type LogEntry struct {
	RequestType string `json:"requestType"`
	Timestamp   string `json:"timestamp"`
	Log         Log    `json:"log"`
}

type Log struct {
	ID             string           `json:"id"`
	StartTime      string           `json:"start_time"`
	Duration       float64          `json:"duration"`
	ExpectedResult *json.RawMessage `json:"expected_result"`
	ActualResult   *json.RawMessage `json:"actual_result"`
	SearchConfig   *json.RawMessage `json:"search_config"`
	WriteConfig    *json.RawMessage `json:"write_config"`
}

type SearchConfig struct {
	ID                string     `json:"ID"`
	RegionID          int        `json:"RegionID"`
	StartRange        string     `json:"StartRange"` // Alternativ time.Time, falls als Datumsformat verwendet
	EndRange          string     `json:"EndRange"`   // Alternativ time.Time, falls als Datumsformat verwendet
	NumberDays        int        `json:"NumberDays"`
	CarType           *string    `json:"CarType,omitempty"`
	OnlyVollkasko     *bool      `json:"OnlyVollkasko,omitempty"`
	MinFreeKilometer  *int       `json:"MinFreeKilometer,omitempty"`
	MinNumberSeats    *int       `json:"MinNumberSeats,omitempty"`
	MinPrice          *float64   `json:"MinPrice,omitempty"`
	MaxPrice          *float64   `json:"MaxPrice,omitempty"`
	Pagination        Pagination `json:"Pagination"`
	Order             string     `json:"Order"`
	PriceBucketWidth  int        `json:"PriceBucketWidth"`
	FreeKmBucketWidth int        `json:"FreeKmBucketWidth"`
}

type GetResponse struct {
	Offers              []GetOffer              `json:"Offers"`
	CarTypeCounts       map[string]int          `json:"CarTypeCounts"`
	FreeKilometerRanges []GetRangeCount         `json:"FreeKilometerRanges"`
	PriceRanges         []GetRangeCount         `json:"PriceRanges"`
	SeatsCounts         map[string]int          `json:"SeatsCounts"`
	VollkaskoCount      GetVollkaskoCountStruct `json:"VollkaskoCount"`
}

type GetOffer struct {
	OfferID       string `json:"OfferID"`
	IsDataCorrect bool   `json:"IsDataCorrect"`
}

type GetRangeCount struct {
	Start int `json:"Start"`
	End   int `json:"End"`
	Count int `json:"Count"`
}

type GetVollkaskoCountStruct struct {
	TrueCount  int `json:"TrueCount"`
	FalseCount int `json:"FalseCount"`
}

type Pagination struct {
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
}

func openLog(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)

	go func() {
		time.Sleep(10 * time.Second)

		pprof.StopCPUProfile()
		fmt.Println("CPU profile stopped")
		os.Exit(0)
	}()

	buf := make([]byte, 0, 1024*1024*10) // Erhöhen Sie den Puffer auf 1 MB
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024*10)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		i += 1
		fmt.Println(i)

		// Parse each JSON line into a LogEntry
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Println("Failed to parse log entry:", i)
			panic(err)
		}

		if entry.RequestType == "PUSH" {
			handlePost(*entry.Log.WriteConfig)
		} else if entry.RequestType == "READ" {
			handleGet(*entry.Log.SearchConfig, &entry.Log)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func ConvertSearchToParams(config SearchConfig) (types.GetParams, error) {
	// Konvertiere Start- und Endzeitpunkt in Unix-Timestamps
	startTime, err := time.Parse(time.RFC3339, config.StartRange)
	if err != nil {
		return types.GetParams{}, err
	}

	endTime, err := time.Parse(time.RFC3339, config.EndRange)
	if err != nil {
		return types.GetParams{}, err
	}

	// Mapping der Sortierung
	sortOrder := 0 // Standardmäßig unsortiert
	if config.Order == "price-asc" {
		sortOrder = 1
	} else if config.Order == "price-desc" {
		sortOrder = -1
	}

	// Baue GetParams
	getParams := types.GetParams{
		RegionID:              uint(config.RegionID),
		TimeRangeStart:        startTime.Unix(),
		TimeRangeEnd:          endTime.Unix(),
		NumberDays:            uint(config.NumberDays),
		SortOrder:             sortOrder,
		Page:                  uint(config.Pagination.Page),
		PageSize:              uint(config.Pagination.PageSize),
		PriceRangeWidth:       uint(config.PriceBucketWidth),
		MinFreeKilometerWidth: uint(config.FreeKmBucketWidth),
		MinNumberSeats: sql.NullInt32{
			Int32: int32Value(config.MinNumberSeats),
			Valid: config.MinNumberSeats != nil,
		},
		MinPrice: sql.NullInt32{
			Int32: int32(*config.MinPrice),
			Valid: config.MinPrice != nil,
		},
		MaxPrice: sql.NullInt32{
			Int32: int32(*config.MaxPrice),
			Valid: config.MaxPrice != nil,
		},
		CarType: sql.NullString{
			String: stringValue(config.CarType),
			Valid:  config.CarType != nil,
		},
		OnlyVollkasko: sql.NullBool{
			Bool:  boolValue(config.OnlyVollkasko),
			Valid: config.OnlyVollkasko != nil,
		},
		MinFreeKilometer: sql.NullInt32{
			Int32: int32Value(config.MinFreeKilometer),
			Valid: config.MinFreeKilometer != nil,
		},
	}

	return getParams, nil
}

// Hilfsfunktionen zum Extrahieren von Werten aus Pointern
func int32Value(ptr *int) int32 {
	if ptr == nil {
		return 0
	}
	return int32(*ptr)
}

func boolValue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

func stringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func handlePost(logData json.RawMessage) {
	var logEntry struct {
		Offers []struct {
			OfferID        string  `json:"OfferID"`
			RegionID       int     `json:"RegionID"`
			CarType        string  `json:"CarType"`
			NumberDays     int     `json:"NumberDays"`
			NumberSeats    int     `json:"NumberSeats"`
			StartTimestamp string  `json:"StartTimestamp"`
			EndTimestamp   string  `json:"EndTimestamp"`
			Price          float64 `json:"Price"`
			HasVollkasko   bool    `json:"HasVollkasko"`
			FreeKilometers int     `json:"FreeKilometers"`
		} `json:"Offers"`
	}

	if err := json.Unmarshal(logData, &logEntry); err != nil {
		fmt.Printf("Failed to parse log data: %v\n", err)
		return
	}

	// Transform offers to the required format
	var transformedOffers []map[string]interface{}
	for _, offer := range logEntry.Offers {
		transformedOffer := map[string]interface{}{
			"ID":                   offer.OfferID,
			"Data":                 "string",
			"mostSpecificRegionID": offer.RegionID,
			"startDate":            parseTimestampToMillis(offer.StartTimestamp),
			"endDate":              parseTimestampToMillis(offer.EndTimestamp),
			"numberSeats":          offer.NumberSeats,
			"price":                offer.Price,
			"carType":              offer.CarType,
			"hasVollkasko":         offer.HasVollkasko,
			"freeKilometers":       offer.FreeKilometers,
		}
		transformedOffers = append(transformedOffers, transformedOffer)
	}

	// Create the final request payload
	payload := map[string]interface{}{
		"offers": transformedOffers,
	}

	// Serialize payload to JSON
	reqBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to serialize request body: %v\n", err)
		return
	}

	// Make the POST request
	url := "http://127.0.0.1:80/api/offers"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to make POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body), ",")
}

// Helper function to convert a timestamp string to milliseconds
func parseTimestampToMillis(timestamp string) int64 {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		fmt.Printf("Failed to parse timestamp: %v\n", err)
		return 0
	}
	return parsedTime.UnixMilli()
}

func handleGet(searchConfig json.RawMessage, log *Log) {
	// Parse the `searchConfig` from the log
	var logEntry struct {
		ID               string   `json:"ID"`
		RegionID         int      `json:"RegionID"`
		StartRange       string   `json:"StartRange"`
		EndRange         string   `json:"EndRange"`
		NumberDays       int      `json:"NumberDays"`
		CarType          *string  `json:"CarType"`
		OnlyVollkasko    *bool    `json:"OnlyVollkasko"`
		MinFreeKilometer *int     `json:"MinFreeKilometer"`
		MinNumberSeats   *int     `json:"MinNumberSeats"`
		MinPrice         *float64 `json:"MinPrice"`
		MaxPrice         *float64 `json:"MaxPrice"`
		Pagination       struct {
			Page     int `json:"Page"`
			PageSize int `json:"PageSize"`
		} `json:"Pagination"`
		Order             string `json:"Order"`
		PriceBucketWidth  int    `json:"PriceBucketWidth"`
		FreeKmBucketWidth int    `json:"FreeKmBucketWidth"`
	}

	if err := json.Unmarshal(searchConfig, &logEntry); err != nil {
		fmt.Printf("Failed to parse search_config: %v\n", err)
		return
	}

	// Transform the logEntry fields into GET query parameters
	params := map[string]string{
		"regionID":              fmt.Sprintf("%d", logEntry.RegionID),
		"timeRangeStart":        fmt.Sprintf("%d", parseTimestampToMillis(logEntry.StartRange)),
		"timeRangeEnd":          fmt.Sprintf("%d", parseTimestampToMillis(logEntry.EndRange)),
		"numberDays":            fmt.Sprintf("%d", logEntry.NumberDays),
		"sortOrder":             logEntry.Order,
		"page":                  fmt.Sprintf("%d", logEntry.Pagination.Page),
		"pageSize":              fmt.Sprintf("%d", logEntry.Pagination.PageSize),
		"priceRangeWidth":       fmt.Sprintf("%d", logEntry.PriceBucketWidth),
		"minFreeKilometerWidth": fmt.Sprintf("%d", logEntry.FreeKmBucketWidth),
	}

	// Handle optional fields
	if logEntry.MinNumberSeats != nil {
		params["minNumberSeats"] = fmt.Sprintf("%d", *logEntry.MinNumberSeats)
	}
	if logEntry.MinPrice != nil {
		params["minPrice"] = fmt.Sprintf("%d", int64(*logEntry.MinPrice))
	}
	if logEntry.MaxPrice != nil {
		params["maxPrice"] = fmt.Sprintf("%d", int64(*logEntry.MaxPrice))
	}
	if logEntry.CarType != nil {
		params["carType"] = *logEntry.CarType
	}
	if logEntry.OnlyVollkasko != nil {
		params["onlyVollkasko"] = fmt.Sprintf("%t", *logEntry.OnlyVollkasko)
	}
	if logEntry.MinFreeKilometer != nil {
		params["minFreeKilometer"] = fmt.Sprintf("%d", *logEntry.MinFreeKilometer)
	}

	// Construct the query string
	queryParts := []string{}
	for key, value := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
	}
	query := strings.Join(queryParts, "&")

	// Construct the URL
	url := fmt.Sprintf("http://127.0.0.1:80/api/offers?%s", query)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to make GET request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Request", logEntry.ID)

	body, _ := io.ReadAll(resp.Body)

	var result types.QueryResponse

	json.Unmarshal(body, &result)

	var expectedResult GetResponse
	var actualLogResult GetResponse

	json.Unmarshal(*log.ExpectedResult, &expectedResult)
	json.Unmarshal(*log.ActualResult, &actualLogResult)

	if len(expectedResult.Offers) != len(result.Offers) {
		fmt.Println("Offers incorrect: Expected:", len(expectedResult.Offers), "Actual:", len(result.Offers), expectedResult.Offers, result.Offers)
		fmt.Println(result.Offers, logEntry.ID)
	} else {
		fmt.Println("Offers correct length", result.Offers)
	}

	for i, value := range result.Offers {

		iid, found := memory.IIDMap[value.ID]

		if !found {
			fmt.Println("IID incorrect", value.ID)
			os.Exit(1)
		}

		offer := memory.OfferMap[iid]

		if i < len(expectedResult.Offers) {
			other := expectedResult.Offers[i]

			if (value.ID != other.OfferID) || (other.IsDataCorrect == false) {
				fmt.Println("Offer incorrect ", value.ID, other.OfferID)

				otherOffer := memory.OfferMap[memory.IIDMap[other.OfferID]]

				spew.Dump(offer)
				spew.Dump(otherOffer)

			}
		}

		fmt.Println("Offer found", value.ID)
		spew.Dump(params)
		spew.Dump(offer)

		startDays := memory.MillisecondsToDays(parseTimestampToMillis(logEntry.StartRange))
		endDays := memory.MillisecondsToDays(parseTimestampToMillis(logEntry.EndRange))
		days := logEntry.NumberDays

		daysBit, err := memory.DaysIndexMap[days].GetBit(int(iid))

		if err != nil {
			fmt.Println("Error getting bit", iid, days)
			os.Exit(1)
		}

		if daysBit == 0 {
			fmt.Println("incorrect in DaysIndexMap", iid, days)
			os.Exit(1)
		} else {
			fmt.Println("found in DaysIndexMap", iid, days)
		}

		found = false

		memory.StartTree.GreaterThanEqual(int32(startDays), func(key int32, iids []int32) {
			found = found || slices.Contains(iids, iid)
		})

		if found {
			fmt.Println("StartDate correct", iid, offer.StartDate, logEntry.StartRange)
		} else {
			fmt.Println("StartDate incorrect", iid, offer.StartDate, logEntry.StartRange)
		}

		found = false

		memory.EndTree.LessThanEqual(int32(endDays), func(key int32, iids []int32) {
			found = found || slices.Contains(iids, iid)
		})

		if found {
			fmt.Println("EndDate correct", iid, offer.EndDate, logEntry.EndRange)
		} else {
			fmt.Println("EndDate incorrect", iid, offer.EndDate, logEntry.EndRange)
		}

		if logEntry.OnlyVollkasko != nil {

			vollkasko, err := memory.VollkaskoIndex.GetBit(int(iid))

			if err != nil {
				fmt.Println("Error getting Vollkasko", iid)
				os.Exit(1)
			}

			if (vollkasko == 1 && offer.HasVollkasko) || (vollkasko == 0 && !offer.HasVollkasko) || (*logEntry.OnlyVollkasko && !offer.HasVollkasko) {
				fmt.Println("Vollkasko correct", iid, vollkasko, offer.HasVollkasko, *logEntry.OnlyVollkasko)
			} else {
				fmt.Println("Vollkasko incorrect", iid, vollkasko, offer.HasVollkasko, *logEntry.OnlyVollkasko)
				os.Exit(1)
			}
		}

		if logEntry.CarType != nil {

			carTypeIndex := memory.GetCarTypeIndex(*logEntry.CarType)

			carType, err := carTypeIndex.GetBit(int(iid))

			if err != nil {
				fmt.Println("Error getting CarType", iid, offer.CarType, *logEntry.CarType)
				os.Exit(1)
			}

			if carType == 1 {
				fmt.Println("CarType correct", iid, carType, offer.CarType, *logEntry.CarType)
			} else {
				fmt.Println("CarType incorrect", iid, carType, offer.CarType, *logEntry.CarType)
				os.Exit(1)
			}
		}

		if logEntry.MinNumberSeats != nil {

			seatsIndex, err := memory.GetNumberOfSeatsIndex(*logEntry.MinNumberSeats)

			if err != nil {
				fmt.Println("Error getting Seats", iid, offer.NumberSeats, *logEntry.MinNumberSeats)
				os.Exit(1)
			}

			seats, err := seatsIndex.GetBit(int(iid))

			if err != nil {
				fmt.Println("Error getting Seats", iid, offer.NumberSeats, *logEntry.MinNumberSeats)
				os.Exit(1)
			}

			if seats == 1 {
				fmt.Println("Seats correct", iid, seats, offer.NumberSeats, *logEntry.MinNumberSeats)
			} else {
				fmt.Println("Seats incorrect", iid, seats, offer.NumberSeats, *logEntry.MinNumberSeats)
				os.Exit(1)
			}
		}

		if (logEntry.MinPrice != nil) && (logEntry.MaxPrice != nil) {

			found = false

			memory.PriceTree.LessThanEqual(int32(*logEntry.MaxPrice), func(key int32, iids []int32) {
				found = found || slices.Contains(iids, iid)
			})

			if found {
				fmt.Println("MaxPrice correct", iid, offer.Price, *logEntry.MaxPrice)
			} else {
				fmt.Println("MaxPrice incorrect", iid, offer.Price, *logEntry.MaxPrice)
				os.Exit(1)
			}

			found = false

			memory.PriceTree.GreaterThanEqual(int32(*logEntry.MinPrice), func(key int32, iids []int32) {
				found = found || slices.Contains(iids, iid)
			})

			if found {
				fmt.Println("MinPrice correct", iid, offer.Price, *logEntry.MinPrice)
			} else {
				fmt.Println("MinPrice incorrect", iid, offer.Price, *logEntry.MinPrice)
				os.Exit(1)
			}

		}

		// region, err := memory.RegionTree.Get(int32(offer.MostSpecificRegionID))

		// if err != nil {
		// 	fmt.Println("Error getting Region", iid, offer.MostSpecificRegionID, logEntry.RegionID)
		// 	os.Exit(1)
		// }

		// if slices.Contains(region, iid) {
		// 	fmt.Println("Region correct", iid, offer.MostSpecificRegionID, logEntry.RegionID)
		// } else {
		// 	fmt.Println("Region incorrect", iid, offer.MostSpecificRegionID, logEntry.RegionID)
		// 	os.Exit(1)
		// }

		found = false
	}

	for key, value := range expectedResult.CarTypeCounts {
		var other int
		if key == "small" {
			other = result.CarTypeCounts.Small
		}

		if key == "sports" {
			other = result.CarTypeCounts.Sports
		}

		if key == "luxury" {
			other = result.CarTypeCounts.Luxury
		}

		if key == "family" {
			other = result.CarTypeCounts.Family
		}

		if value != (other) {
			fmt.Println("CarTypeCount incorrect", key, value, other)
			os.Exit(1)
		}
	}

	if len(expectedResult.SeatsCounts) > len(result.SeatsCount) {
		fmt.Println("SeatsCount incorrect: Expected:", len(expectedResult.SeatsCounts), "Actual:", len(result.SeatsCount), expectedResult.SeatsCounts)
		os.Exit(1)
	}

	for seatTypeString, value := range expectedResult.SeatsCounts {
		numberSeats, _ := strconv.ParseInt(seatTypeString, 10, 32)

		found := false

		for _, seat := range result.SeatsCount {
			if seat.NumberSeats == int(numberSeats) {
				if value != int(seat.Count) {
					fmt.Println("SeatsCount incorrect ", seat.NumberSeats, seat.Count, value, expectedResult.SeatsCounts)
					os.Exit(1)

				} else {
					fmt.Println("SeatsCount correct", seat.NumberSeats, seat.Count, value)
				}

				found = true

				break
			}
		}

		if !found {

			spew.Dump(result.SeatsCount)
			fmt.Println("SeatsCount incorrect", numberSeats, expectedResult.SeatsCounts)
			os.Exit(1)
		}

	}

	for i, value := range expectedResult.FreeKilometerRanges {
		if i >= len(result.FreeKilometerRange) {
			fmt.Println("FreeKilometerRanges incorrect: Expected:", len(expectedResult.FreeKilometerRanges), "Actual:", len(result.FreeKilometerRange), expectedResult.FreeKilometerRanges)
			os.Exit(1)
		}
		other := result.FreeKilometerRange[i]

		if (value.Count != other.Count) || (value.End != other.End) || (value.Start != other.Start) {
			fmt.Print("FreeKilometerRange incorrect")

			fmt.Print(" Count: ", value.Count, other.Count)
			fmt.Print(" Start: ", value.Start, other.Start)
			fmt.Print(" End: ", value.End, other.End, "\n")

			os.Exit(1)
		}
	}

	for i, value := range expectedResult.PriceRanges {
		if i >= len(result.PriceRanges) {
			fmt.Println("PriceRanges incorrect: Expected:", len(expectedResult.PriceRanges), "Actual:", len(result.PriceRanges), expectedResult.PriceRanges)
			os.Exit(1)
		}
		other := result.PriceRanges[i]

		if (value.Count != other.Count) || (value.End != other.End) || (value.Start != other.Start) {
			fmt.Print("PriceRange incorrect ")

			fmt.Print("Count: ", value.Count, other.Count)
			fmt.Print("Start: ", value.Start, other.Start)
			fmt.Print("End: ", value.End, other.End)
			os.Exit(1)

		}
	}

}
