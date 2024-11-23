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
	printBounds()
	printParsedStructure()
}

func GetRegionBounds(region uint) (uint, uint, uint, uint) {
	r := bounds[region]

	if region == 6 {
		return 56, 57, 121, 124
	}

	return r.Min, r.Max, r.Min, r.Max
}

func printRegionStructure(region Region, level int) {
	// Indentation based on level to show hierarchy
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}

	// Print the current region
	fmt.Printf("%sRegion ID: %d, Name: %s\n", indent, region.ID, region.Name)

	// Recursively print subregions
	for _, subregion := range region.Subregions {
		printRegionStructure(subregion, level+1)
	}
}

func printParsedStructure() {
	var root Region
	err := json.Unmarshal(regionsFile, &root)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("Parsed Region Structure:")
	printRegionStructure(root, 0) // Start from the root at level 0
}
