package routes

import (
	"math/rand"
	"testing"
)

func getRandomSortOrder() string {
	if rand.Intn(2) == 0 {
		return "price-asc"
	} else {
		return "price-desc"
	}
}

func BenchmarkSortOrder(b *testing.B) {
	// run the Fib function b.N times
	sortOrder := getRandomSortOrder()

	if sortOrder != "price-asc" && sortOrder != "price-desc" {
		b.Errorf("Invalid sort order")
	}

}
