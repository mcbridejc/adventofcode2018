package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Claim struct {
	id int
	top int
	left int
	width int
	height int
}

// Create an iterator based on a scanner, so we don't have to keep a full file's worth
// of claims in memory at once #unnecessaryoptimization
func NewClaimInputIterator(filepath string) (nextFunc func() (claim Claim, ok bool)) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile("#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)")

	nextFunc = func() (Claim, bool) {
		var claim Claim
		ok := scanner.Scan()
		if !ok {
			return claim, false
		}
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			panic(fmt.Sprintf("Regexp failed on line %s", line))
		}
		
		// Matching the regexp should guarantee it parses as an int; not checking error
		claim.id, _ = strconv.Atoi(matches[1])
		claim.left, _ = strconv.Atoi(matches[2])
		claim.top, _ = strconv.Atoi(matches[3])
		claim.width, _ = strconv.Atoi(matches[4])
		claim.height, _ = strconv.Atoi(matches[5])
		return claim, true
	}
	return
}

type Location struct {
	x int
	y int
}
func main() {
	nextClaim := NewClaimInputIterator("day3_input.txt")

	locationCounts := make(map[Location]int)

	for {
		claim, ok := nextClaim()
		if !ok {
			break
		}
		for x := claim.left; x < claim.left + claim.width; x += 1 {
			for y := claim.top; y < claim.top + claim.height; y += 1 {
				locationCounts[Location{x, y}] += 1
			}	
		}
	}

	conflictCount := 0
	for _, count := range locationCounts {
		if count >= 2 {
			conflictCount += 1
		}
	}

	fmt.Printf("Number of squares with multiple claims: %d\n", conflictCount)

	nextClaim = NewClaimInputIterator("day3_input.txt")
	for {
		claim, ok := nextClaim()
		if !ok { 
			break
		}
		conflictFound := false
		for x := claim.left; (x < claim.left + claim.width) && !conflictFound; x += 1 {
			for y := claim.top; (y < claim.top + claim.height) &&  !conflictFound; y += 1 {
				if locationCounts[Location{x, y}] != 1 {
					conflictFound = true
				}
			}	
		}
		if !conflictFound {
			fmt.Printf("No conflict found for patch #%d\n", claim.id)
			// The instructions say there will be only one, so we could break...
			// but may as well finish checking all to validate there's only one
		}
	}
}

