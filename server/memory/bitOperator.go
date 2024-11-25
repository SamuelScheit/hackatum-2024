package memory

import (
	"fmt"
	"math/bits"
)

// BitArray represents a dynamic array of bits, backed by a slice of uint64
type BitArray struct {
	data []uint64 // Each uint64 stores 64 bits
	size int      // Total number of bits
}

// NewBitArray creates a new BitArray with the specified number of bits
func NewBitArray(size int) *BitArray {
	numWords := (size + 63) / 64 // Calculate the number of uint64s needed
	return &BitArray{
		data: make([]uint64, numWords),
		size: size,
	}
}

// ensureCapacity expands the BitArray if the given position is out of range
func (ba *BitArray) ensureCapacity(pos int) {
	if pos >= ba.size {
		newSize := pos + 1
		numWords := (newSize + 63) / 64 // Calculate the number of uint64s needed
		if numWords > len(ba.data) {
			newData := make([]uint64, numWords)
			copy(newData, ba.data)
			ba.data = newData
		}
		ba.size = newSize
	}
}

// Clear resets all bits in the BitArray to 0
func (ba *BitArray) Clear() {
	for i := range ba.data {
		ba.data[i] = 0
	}
	ba.size = len(ba.data) * 64
}

func (ba *BitArray) CopyFrom(other *BitArray) {
	ba.size = other.size
	ba.data = make([]uint64, len(other.data))
	copy(ba.data, other.data)
}

// SetBit sets the bit at position pos to 1, expanding the BitArray if needed
func (ba *BitArray) SetBit(pos int) {
	if pos < 0 {
		panic("position must be non-negative")
	}
	ba.ensureCapacity(pos)
	wordIndex := pos / 64 // Determine which uint64
	bitIndex := pos % 64  // Determine which bit within the uint64
	ba.data[wordIndex] |= 1 << bitIndex
}

// GetBit gets the value of the bit at position pos (returns 0 or 1)
func (ba *BitArray) GetBit(pos int) (uint64, error) {
	if pos < 0 {
		return 0, fmt.Errorf("position must be non-negative")
	}
	if pos >= ba.size {
		return 0, nil // Bits beyond the current size are implicitly 0
	}
	wordIndex := pos / 64
	bitIndex := pos % 64
	return (ba.data[wordIndex] >> bitIndex) & 1, nil
}

// LogicalAnd performs a logical AND operation between two BitArrays
func LogicalAnd(ba1, ba2 *BitArray) *BitArray {
	minSize := ba1.size
	if ba2.size < minSize {
		minSize = ba2.size
	}

	result := NewBitArray(minSize)

	for i := 0; i < len(ba1.data) && i < len(ba2.data); i++ {
		result.data[i] = ba1.data[i] & ba2.data[i]
	}

	result.size = minSize

	return result
}

func LogicalAndInPlace(ba1, ba2 *BitArray) {
	// Perform word-by-word AND operation
	ba2Len := len(ba2.data)
	for i := range ba1.data {
		if i >= ba2Len {
			ba1.data = ba1.data[:i]
			ba1.size = i * 64
			break
		}
		ba1.data[i] = ba1.data[i] & ba2.data[i]
	}
}

func LogicalOr(ba1, ba2 *BitArray) *BitArray {
	maxSize := ba1.size
	if ba2.size > maxSize {
		maxSize = ba2.size
	}

	result := NewBitArray(maxSize)

	for i := range result.data {
		word1 := uint64(0)
		word2 := uint64(0)

		if i < len(ba1.data) {
			word1 = ba1.data[i]
		}
		if i < len(ba2.data) {
			word2 = ba2.data[i]
		}

		result.data[i] = word1 | word2
	}

	return result
}

func LogicalOrInPlace(ba1, ba2 *BitArray) {
	// Ensure ba1 is large enough to accommodate the result
	if ba2.size > ba1.size {
		ba1.ensureCapacity(ba2.size - 1)
	}

	// Perform word-by-word OR operation
	for i := range ba2.data {
		if i < len(ba1.data) {
			ba1.data[i] |= ba2.data[i]
		} else {
			// If ba1 has fewer words, append the remaining words from ba2
			ba1.data = append(ba1.data, ba2.data[i])
		}
	}
}

// PrintBits displays the entire bit array as a binary string (for debugging purposes)
func (ba *BitArray) PrintBits() string {
	bitStr := ""
	for i := 0; i < ba.size; i++ {
		bit, _ := ba.GetBit(i)
		bitStr += fmt.Sprintf("%d", bit)
	}
	return bitStr
}

// countSetBits counts the number of bits set in a BitArray
func (ba *BitArray) CountSetBits() int {
	count := 0
	for i := range ba.data {
		count += bits.OnesCount64(ba.data[i])
	}
	return count
}

func (ba *BitArray) CountUnsetBits() int {
	count := 0
	for i := range ba.data {
		count += 64 - bits.OnesCount64(ba.data[i])
	}
	return count
}

func (ba *BitArray) Copy() *BitArray {
	newData := make([]uint64, len(ba.data))
	copy(newData, ba.data)
	return &BitArray{
		data: newData,
		size: ba.size,
	}
}
