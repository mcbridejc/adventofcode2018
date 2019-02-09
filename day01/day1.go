package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)
func main() {
	f, err := os.Open("day1_input.txt")
	if err != nil {
		panic(err)
	}
	
	scanner := bufio.NewScanner(f)

	input := make([]int, 0)
	for scanner.Scan() {
		num, _ := strconv.Atoi(scanner.Text())
		input = append(input, num)
	}

	sum := 0
	for _, num := range input {
		sum += num
	}
	fmt.Printf("Sum: %d\n", sum)


	freq := 0
	freqs := make(map[int]bool)
	i := 0
	for {
		_, found := freqs[freq]
		if found {
			fmt.Printf("First repeat: %d\n", freq)
			break
		}
		freqs[freq] = true
	
		num := input[i]
		freq += num	
		i = (i+1) % len(input)
		//fmt.Println(freq)
	}
}