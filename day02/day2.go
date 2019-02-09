package main

import (
	"bufio"
	"fmt"
	"os"
)

func CountDiff(a string, b string) int {
	diff := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

func main() {
	f, err := os.Open("day2_input.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	ids := make([]string, 0)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	count2 := 0
	count3 := 0
	for _, id := range ids {
		countMap := make(map[rune]int)
		for _, ch := range id {
			countMap[ch] += 1
		}
		foundA2 := false
		foundA3 := false
		for _, value := range countMap {
			if value == 2 && !foundA2 {
				count2 += 1
				foundA2 = true
			}
			if value == 3 && !foundA3 {
				count3 += 1
				foundA3 = true
			}
		}
	}

	fmt.Printf("Part 1 checksum: %d\n", count2 * count3)
	fmt.Printf("count-2: %d, count-3: %d\n", count2, count3)

	for _, a := range ids {
		for _, b := range ids {
			if CountDiff(a, b) == 1 {
				fmt.Printf("Found pair:\n%s\n%s\n", a, b)
				s := ""
				for i := 0; i < len(a); i++ {
					if a[i] == b[i] {
						s += string([]byte{a[i]})
					}
				}
				fmt.Printf("Answer: %s\n", s)
			}
		}
	}
}