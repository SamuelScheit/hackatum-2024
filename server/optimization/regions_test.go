package optimization

import (
	"fmt"
	"testing"
)

func TestRegions(t *testing.T) {
	Init()

	fmt.Println(GetRegionBounds(6))

}
