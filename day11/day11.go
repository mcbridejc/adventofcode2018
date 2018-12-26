package main

import (
	"fmt"
	"math"
)

// Get the power for a given fuel cell coordinate
// x/y are on range (1, 300)
func FuelCellPower(x, y, gridSerial int) (int) {
	rackId := x + 10
	power := rackId * y
	power += gridSerial
	power *= rackId
	// Keep only the hundreds digit (12345 -> 3)
	power = (power / 100) % 10
	return power - 5
}

const GRIDSERIAL = 8561

func FindMaxRegion(size int, grid [][]int) (maxRegionValue, maxRegionX, maxRegionY int) {
	maxRegionValue = math.MinInt32
	// Run 3x3 filter kernel over the grid
	for kx := 1; kx <= 300 - size + 1; kx += 1 {
		for ky := 1; ky <= 300 - size + 1; ky += 1 {
			accum := 0
			for x := kx; x < kx + size; x += 1 {
				for y := ky; y < ky + size; y += 1 {
					accum += grid[x-1][y-1]//FuelCellPower(x, y, serial)
				}
			}
			if accum > maxRegionValue {
				maxRegionValue = accum
				maxRegionX = kx
				maxRegionY = ky
			}
		}
	}
	return maxRegionValue, maxRegionX, maxRegionY
}

func main() {
	testCasePower := FuelCellPower(101, 153, 71)
	if testCasePower != 4 {
		fmt.Printf("FAILED TEST CASE! Test case answer: %d\n", testCasePower)
	}

	fuelCellGrid := make([][]int, 300)
	for x := 0; x < 300; x += 1 {
		fuelCellGrid[x] = make([]int, 300)
		for y := 0; y < 300; y += 1 {
			fuelCellGrid[x][y] = FuelCellPower(x+1, y+1, GRIDSERIAL)
		}
	}
	
	maxRegionValue := math.MinInt32
	maxRegionX := 0
	maxRegionY := 0
	
	// Run 3x3 filter kernel over the grid
	maxRegionValue, maxRegionX, maxRegionY = FindMaxRegion(3, fuelCellGrid)
	fmt.Println("Part 1\n------")
	fmt.Printf("Max region is %d,%d with power=%d\n", maxRegionX, maxRegionY, maxRegionValue)

	maxRegionValue = math.MinInt32
	maxRegionSize := 0
	// Run all possible size filter kernels
	for size := 1; size <= 300; size+= 1 {
		value, x, y := FindMaxRegion(size, fuelCellGrid)
		if value > maxRegionValue {
			maxRegionValue = value
			maxRegionX = x
			maxRegionY = y
			maxRegionSize = size
		}
	}
	fmt.Println("Part 2\n------")
	fmt.Printf("Max region is %d,%d,%d with power=%d\n", maxRegionX, maxRegionY, maxRegionSize, maxRegionValue)
}