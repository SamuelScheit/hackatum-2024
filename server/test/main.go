package main

import (
	"bufio"
	"bytes"
	"checkmate/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
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
	ID             string                 `json:"id"`
	StartTime      string                 `json:"start_time"`
	Duration       float64                `json:"duration"`
	ExpectedResult map[string]interface{} `json:"expected_result"`
	ActualResult   map[string]interface{} `json:"actual_result"`
	SearchConfig   *json.RawMessage       `json:"search_config"`
	WriteConfig    *json.RawMessage       `json:"write_config"`
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
			panic(err)
		}

		if entry.RequestType == "PUSH" {
			handlePost(*entry.Log.WriteConfig)
		} else if entry.RequestType == "READ" {
			handleGet(*entry.Log.SearchConfig)
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

func handleGet(searchConfig json.RawMessage) {
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

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body), ",")
}
