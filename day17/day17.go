package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type DirtType int
const (
	Sand DirtType = iota
	Clay
	StaticWater
	FlowingWater
)

type DirtMap struct {
	clayMap map[[2]int]bool
	waterMap map[[2]int]bool
	minX int
	maxX int
	minY int
	maxY int
}

func NewDirtMap() *DirtMap {
	var dm DirtMap
	dm.clayMap = make(map[[2]int]bool)
	dm.waterMap = make(map[[2]int]bool)
	dm.minX = 1000000000
	dm.maxX = 0
	dm.minY = 1000000000
	dm.maxY = 0
	return &dm
}

func (dm DirtMap) Get(x, y int) DirtType {
	
	if static, present := dm.waterMap[[2]int{x, y}]; present {
		if static {
			return StaticWater
		} else {
			return FlowingWater
		}
	}
	if _, present := dm.clayMap[[2]int{x, y}]; present {
		return Clay
	}
	return Sand
}

func (dm *DirtMap) Set(x int, y int, dtype DirtType) {
	switch dtype {
	case Clay:
		dm.clayMap[[2]int{x, y}] = true
		if y < dm.minY {
			dm.minY = y
		}
		if y > dm.maxY {
			dm.maxY = y
		}
		if x < dm.minX {
			dm.minX = x
		}
		if x > dm.maxX {
			dm.maxX = x
		}
	case FlowingWater:
		dm.waterMap[[2]int{x, y}] = false
	case StaticWater:
		dm.waterMap[[2]int{x, y}] = true
	}
}

func ReadInput(filepath string) *DirtMap {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	dmap := NewDirtMap()
	re := regexp.MustCompile("([xy])=([0-9]+), [xy]=([0-9]+)..([0-9]+)")
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if match == nil {
			fmt.Println("Couldn't parse line: ", scanner.Text())
			continue
		}
		singletonAxis := match[1]
		u, _ := strconv.Atoi(match[2])
		v0, _ := strconv.Atoi(match[3])
		v1, _ := strconv.Atoi(match[4])
		for v := v0; v <= v1; v++ {
			var x, y int
			if singletonAxis == "x" {
				x = u
				y = v
			} else {
				x = v
				y = u
			}
			dmap.Set(x, y, Clay)
		}
	}
	return dmap
}

type Position [2]int

func FillWater(dm *DirtMap, x int, y int, delta int) (int, bool) {
	// Seek until obstructed, or until we find a spot where we are free to drop
	// Depending on the outcome, all of the locations we 
	leak := false
	var dx int
	for dx = x+delta; ; dx += delta {
		dtype := dm.Get(dx, y)
		if dtype == Clay {
			// We hit an edge, without spilling. 
			break
		}
		belowDtype := dm.Get(dx, y+1)
		if belowDtype == Sand {
			// We found a spill point
			leak = true
			break
		}
	}

	// In cases where we hit clay instead of a leak, the last dx value actually isn't water
	if !leak { 
		dx -= delta
	}
	// Return the last x position the water reached, and whether or not it will 
	// fall from this location
	return dx, leak
}

func DropWater(dm *DirtMap, x int, y int) []Position{
	// Drop until we hit something besides sand, marking each tile we fall 
	// through as flowing water
	for {
		if y > dm.maxY {
			// We fell of the map
			return make([]Position, 0)
		}
		dtypeBelow := dm.Get(x, y+1)
		if y < dm.minY {
			y++
			continue
		}
		if dtypeBelow == Sand {
			dm.Set(x, y, FlowingWater)
		} else if dtypeBelow == FlowingWater {
			// We landed on an already overflowed container
			// We need not continue
			dm.Set(x, y, FlowingWater)
			return make([]Position, 0)	
		}else {
			break
		}
		y++
	}

	// We landed on something solid. Now keep filling until we overflow. 
	for {
		// Seek left until obstructed, or until we find another drop
		leftBound, leftLeak := FillWater(dm, x, y, -1)
		rightBound, rightLeak := FillWater(dm, x, y, 1)
		leak := leftLeak || rightLeak
		
		// Set as water, either static or flowing based on whether there's an overflow
		for i := leftBound; i <= rightBound; i++ {
			if leak {
				dm.Set(i, y, FlowingWater)
			} else {
				dm.Set(i, y, StaticWater)
			}
		}

		if leak {
			newDropSources := make([]Position, 0)
			if leftLeak {
				newDropSources = append(newDropSources, Position{leftBound, y})
			}
			if rightLeak {
				newDropSources = append(newDropSources, Position{rightBound, y})
			}
			return newDropSources
		} else {
			// Move up and fill again
			y--
		}
	}
}

func main() {
	dirtMap := ReadInput("day17_input.txt")

	fmt.Printf("Read map with x %d..%d and y %d..%d\n", dirtMap.minX, dirtMap.maxX, dirtMap.minY, dirtMap.maxY)

	dropSources := []Position{{500, 0}}
	iterCount := 0
	for len(dropSources) > 0 {
		nextSource := dropSources[0]
		dropSources = dropSources[1:]
		// Each iteration may return 0 to 2 new drop locations to iterate on
		newDropSources := DropWater(dirtMap, nextSource[0], nextSource[1])
		if len(newDropSources) > 0 {
			dropSources = append(dropSources, newDropSources...)
		}
		iterCount++
		if iterCount % 100 == 0 {
			fmt.Println("Iteration: ", iterCount, "(", len(dropSources), ")")
		}
	}

	fmt.Printf("Found %d water tiles\n", len(dirtMap.waterMap))
	staticCount := 0
	for _, w := range dirtMap.waterMap {
		if w {staticCount++}
	}
	fmt.Printf("Found %d static water tiles\n", staticCount)
	waterCnt := 0
	for row := dirtMap.minY-2; row <= dirtMap.maxY+2; row++ {
		for col := dirtMap.minX-2; col < dirtMap.maxX+2; col++ {
			switch dirtMap.Get(col, row) {
			case Sand:
				fmt.Printf(".")
			case Clay:
				fmt.Printf("#")
			case StaticWater:
				fmt.Printf("~")
				waterCnt++
			case FlowingWater:
				fmt.Printf("|")
				waterCnt++
			}
		}
		fmt.Printf("\n")
	}
}