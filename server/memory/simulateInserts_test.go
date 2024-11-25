package memory

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// SimulateAttributes performs the simulation for hasVollkasko and isFamilyCar with timing
func TestSimulateAttributes(t *testing.T) {
	const numEntries = 20000

	// Initialize the BitArrays
	hasVollkasko := NewBitArray(numEntries)
	isFamilyCar := NewBitArray(numEntries)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Measure time for inserting random bits
	startInsert := time.Now()
	for i := 0; i < numEntries; i++ {
		if rand.Intn(2) == 1 {
			hasVollkasko.SetBit(i)
		}
		if rand.Intn(2) == 1 {
			isFamilyCar.SetBit(i)
		}
	}
	elapsedInsert := time.Since(startInsert)
	fmt.Printf("Time for inserting random bits: %v\n", elapsedInsert)

	// Measure time for logical AND operation
	startAnd := time.Now()
	result := LogicalAnd(hasVollkasko, isFamilyCar)
	elapsedAnd := time.Since(startAnd)
	fmt.Printf("Time for logical AND operation: %v\n", elapsedAnd)

	// Measure time for counting set bits in the result
	startCount := time.Now()
	count := 0
	for i := 0; i < numEntries; i++ {
		bit, _ := result.GetBit(i)
		if bit == 1 {
			count++
		}
	}
	elapsedCount := time.Since(startCount)
	fmt.Printf("Time for counting bits: %v\n", elapsedCount)

	// Output results
	fmt.Printf("Number of '1' bits in hasVollkasko: %d\n", hasVollkasko.CountSetBits())
	fmt.Printf("Number of '1' bits in isFamilyCar: %d\n", isFamilyCar.CountSetBits())
	fmt.Printf("Number of '1' bits after AND: %d\n", count)

	// Combined operation timing
	totalElapsed := elapsedInsert + elapsedAnd + elapsedCount
	fmt.Printf("Total time for all operations: %v\n", totalElapsed)
}
