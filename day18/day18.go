package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type TileClass int
const (
	Empty TileClass = iota
	Trees
	Woodshop
)
type Map [][]TileClass

// Set a tile in the map, growing the map as necessary to accomodate
func (m *Map) Set(x int, y int, tile TileClass) {
	width := m.Width()
	height := m.Height()

	// Resize storage as necessary
	if x >= width {
		arraysToAdd := x - width + 1
		for i := 0; i < arraysToAdd; i++ {
			*m = append(*m, make([]TileClass, height))
		}
		width = x + 1
	}
	if y >= height {
		for i := 0; i < width; i++ {
			(*m)[i] = append((*m)[i], make([]TileClass, y - height + 1)...)
		}
	}
	(*m)[x][y] = tile
}

func (m Map) Width() int {
	return len(m)
}

func (m Map) Height() int {
	height := 0
	if len(m) > 0 {
		height = len((m)[0])
	}
	return height 
}

func (m Map) CountAdjacent(x int, y int) map[TileClass]int {
	ret := make(map[TileClass]int)
	ret[Empty] = 0
	ret[Trees] = 0
	ret[Woodshop] = 0
	deltas := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, 
		{0, -1},           {0, 1}, 
		{1, -1},  {1, 0},  {1, 1}}
	for _, delta := range deltas {
		xp := x + delta[0]
		yp := y + delta[1]
		if xp < 0 || xp >= m.Width() || yp < 0 || yp >= m.Height() {
			continue
		}
		ret[m[xp][yp]] += 1
	}
	return ret
}

func Evolve(m Map) Map {
	var new Map

	// Set the size up front so we allocate just once
	new.Set(m.Width()-1, m.Height()-1, Empty)

	for x := 0; x < m.Width(); x++ {
		for y := 0; y < m.Height(); y++ {
			counts := m.CountAdjacent(x, y)
			switch(m[x][y]) {
			case Empty:
				if counts[Trees] >= 3 {
					new.Set(x, y, Trees)
				} else {
					new.Set(x, y, Empty)
				}
			case Trees:
				if counts[Woodshop] >= 3 {
					new.Set(x, y, Woodshop)
				} else {
					new.Set(x, y, Trees)
				}
			case Woodshop:
				if counts[Woodshop] >= 1 && counts[Trees] >= 1 {
					new.Set(x, y, Woodshop)
				} else {
					new.Set(x, y, Empty)
				}
			}
		}
	}
	return new
}

func MapEq(a Map, b Map) bool {
	if a.Height() != b.Height() || a.Width() != b.Width() {
		return false
	}
	for y := 0; y < a.Height(); y++ {
		for x := 0; x < a.Width(); x++ {
			if a[x][y] != b[x][y] {
				return false
			}
		}
	}
	return true
}

func ReadInput(filepath string) Map {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	
	m := make(Map, 0)
	y := 0
	for scanner.Scan() {
		for x, rn := range scanner.Text() {
			if rn == '.' {
				m.Set(x, y, Empty)
			} else if rn == '|' {
				m.Set(x, y, Trees)
			} else if rn == '#' {
				m.Set(x, y, Woodshop)
			}
		}
		y++
	}
	return m
}

func PrintMap(m Map) {
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			switch(m[x][y]) {
			case Empty:
				fmt.Printf(".")
			case Trees:
				fmt.Printf("|")
			case Woodshop:
				fmt.Printf("#")
			}
		}
		fmt.Printf("\n")
	}
}

func ResourceValue(m Map) int {
	trees := 0
	woodshop := 0
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			switch m[x][y] {
			case Trees:
				trees++
			case Woodshop:
				woodshop++
			}
		}
	}
	return trees * woodshop
}

func main() {
	verbose := flag.Bool("verbose", false, "Print more stuff")
	flag.Parse()
	fmt.Println("Reading")

	m := ReadInput("day18_input.txt")

	fmt.Printf("Read map of size %dx%d\n", m.Width(), m.Height())

	if *verbose {
		PrintMap(m)
	}


	for generation := 0; generation < 10; generation++ {
		m = Evolve(m)
		if generation % 1000 == 0 {
			fmt.Println("Gen ", generation+1)
		}
		if (*verbose) {
			fmt.Println("Generation ", generation)
			PrintMap(m)
		}
	}

	fmt.Printf("Final resource value after 10 iterations: %d\n", ResourceValue(m))

	
	// Try to find the value after a large number of generations, by assuming it will 
	// generate a repeated pattern before then
	m = ReadInput("day18_input.txt")
	pastMaps := make([]Map, 0)
	MaxHistory := 100

	repeat := false
	repeatStart := 0
	repeatPeriod := 0
	for generation := 0; generation < 600; generation++ {
		m = Evolve(m)
		fmt.Printf("gen %d: %d\n", generation + 1, ResourceValue(m))
		for i, pm := range pastMaps {
			if MapEq(pm, m) {
				//fmt.Printf("Map %d == Map %d\n", generation+1, i)
				repeat = true
				repeatPeriod = len(pastMaps) - i
				repeatStart = generation + 1 - repeatPeriod
			}
		}
		if repeat {
		 	break
		}
		pastMaps = append(pastMaps, m)
		if len(pastMaps) > MaxHistory {
			pastMaps = pastMaps[1:]
		}
	}

	repeatingScores := make([]int, 0)
	for i := len(pastMaps) - repeatPeriod; i < len(pastMaps); i++ {
		repeatingScores = append(repeatingScores, ResourceValue(pastMaps[i]))
	}
	for i, v := range repeatingScores {
		fmt.Println(i, ": ", v)
	}

	largeGenerations := 1000000000
	repeatIdx := (largeGenerations - repeatStart) % repeatPeriod

	fmt.Println("Repeat start: ", repeatStart)
	fmt.Println("Repeat period: ", repeatPeriod)
	fmt.Println("RepeatingScore[0]: ", repeatingScores[0], len(repeatingScores))
	fmt.Printf("Predicted resource value after %d iterations: %d\n", largeGenerations, repeatingScores[repeatIdx])
}