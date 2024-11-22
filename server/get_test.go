package main

import (
	"bytes"
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

func getRandomCarType() []byte {
	switch rand.Intn(4) {
	case 0:
		return carTypeSmallBytes
	case 1:
		return carTypeSportsBytes
	case 2:
		return carTypeLuxuryBytes
	case 3:
		return carTypeFamilyBytes
	}

	return []byte{}
}

func BenchmarkSortOrder(b *testing.B) {
	// run the Fib function b.N times
	sortOrder := getRandomSortOrder()

	if sortOrder != "price-asc" && sortOrder != "price-desc" {
		b.Errorf("Invalid sort order")
	}

	getRandomSortOrder()
}

func BenchmarkCarType(b *testing.B) {
	// run the Fib function b.N times
	carTypeBytes := getRandomCarType()

	var carType = -1

	if bytes.Equal(carTypeBytes, carTypeSmallBytes) {
		carType = carTypeSmall
	} else if bytes.Equal(carTypeBytes, carTypeSportsBytes) {
		carType = carTypeSports
	} else if bytes.Equal(carTypeBytes, carTypeLuxuryBytes) {
		carType = carTypeLuxury
	} else if bytes.Equal(carTypeBytes, carTypeFamilyBytes) {
		carType = carTypeFamily
	} else if len(carTypeBytes) > 0 {
		return
	}

	_ = carType
}
