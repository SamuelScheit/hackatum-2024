package optimization

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

type Region struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Subregions []Region `json:"subregions"`
}

type LeafBounds struct {
	Min uint
	Max uint
}

var bounds map[uint]LeafBounds

//go:embed regions.json
var regionsFile []byte

func buildDict() {
	var root Region
	err := json.Unmarshal(regionsFile, &root)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	bounds = make(map[uint]LeafBounds)

	computeLeafBounds(root)
}

func computeLeafBounds(region Region) (uint, uint) {
	if len(region.Subregions) == 0 {
		bounds[region.ID] = LeafBounds{Min: region.ID, Max: region.ID}
		return region.ID, region.ID
	}

	minID, maxID := ^uint(0), uint(0)
	for _, subregion := range region.Subregions {
		subMin, subMax := computeLeafBounds(subregion)
		if subMin < minID {
			minID = subMin
		}
		if subMax > maxID {
			maxID = subMax
		}
	}

	bounds[region.ID] = LeafBounds{Min: minID, Max: maxID}
	return minID, maxID
}

func printBounds() {
	fmt.Println("Region ID to Min/Max Leaf Bounds Mapping:")
	for id, bound := range bounds {
		fmt.Printf("Region ID %d: Min = %d, Max = %d\n", id, bound.Min, bound.Max)
	}
}

func Init() {
	buildDict()
}

func GetRegionBounds(region uint) (uint, uint) {
	r := bounds[region]

	return r.Min, r.Max
}
