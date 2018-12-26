package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ReadInput(filepath string) (string, map[string]string) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	// Look for initial state on first line
	re := regexp.MustCompile("initial state: ([#.]+)")
	match := re.FindStringSubmatch(scanner.Text())
	if match == nil {
		panic(fmt.Sprintf("No initial state found in %s", scanner.Text()))
	}
	initState := match[1]

	stateTable := make(map[string]string)

	re = regexp.MustCompile("([#.]{5}) => ([#.])")
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}
		stateTable[match[1]] = match[2]
	}

	return initState, stateTable
}

func Evolve(state string, stateTable map[string]string) (string, int) {
	firstIndex := strings.Index(state, "#")
	lastIndex := strings.LastIndex(state, "#")
	prependCount := 0
	appendCount := 0

	if firstIndex < 4 {
		prependCount = 4 - firstIndex;
	}
	if lastIndex >= len(state) - 4 {
		appendCount = len(state) + 4 - lastIndex
	}
	firstIndex += prependCount
	lastIndex += prependCount

	state = strings.Repeat(".", prependCount) + state + strings.Repeat(".", appendCount)

	newState := []byte(strings.Repeat(".", len(state) + prependCount + appendCount))
	for i := firstIndex - 2; i<=lastIndex + 2; i += 1 {
		subStr := state[i-2:i+3]
		
		// Assume missing entries result in a '.' so that we can use input
		// that only defines the states that result in '#' (i.e. the example)
		tableVal, ok := stateTable[subStr]
		if !ok {
			tableVal = "."
		}
		
		newState[i] = tableVal[0]
	}

	newStateStr := string(newState)
	newFirstIndex := strings.Index(newStateStr, "#")
	newLastIndex := strings.LastIndex(newStateStr, "#")
	shift := prependCount - newFirstIndex
	return newStateStr[newFirstIndex:newLastIndex+1], shift
}

func score(zeroIndex int64, state string) int64 {
	plantChecksum := int64(0)
	for i := 0; i < len(state); i += 1 { 
		if state[i] == '#' {
			plantChecksum += int64(i) - int64(zeroIndex)
		}
	}
	return plantChecksum
}

func main() {
	flag.Parse()
	inputFile := flag.Args()[0]

	fmt.Println("Reading from ", inputFile)
	initState, stateTable := ReadInput(inputFile)

	// Print the input, just to see that it was read appropriately
	fmt.Printf("Initial state: '%s'\n", initState)
	for k, v := range stateTable {
		fmt.Printf("%s -> %s\n", k, v)
	}

	// We would do so many generations, but its not computationally feasible
	// So we'll just do enough for it to stabililze
	//totalGenerations := int64(50*1000*1000*1000)
	totalGenerations := int64(200)
	
	state := initState
	var zeroIndex int64
	for i := int64(0); i < totalGenerations; i += 1 {
		if i%1000000 == 0 {
			fmt.Println(i/1000000)
		}
		shift := 0
		state, shift = Evolve(state, stateTable)
		zeroIndex += int64(shift)
		fmt.Println("% 4d: ", i+1, (score(zeroIndex, state)))	
	}

	fmt.Printf("Final state:\n     %s\n", state)
	plantChecksum := score(zeroIndex, state)

	fmt.Println("Part 1\n------")
	fmt.Printf("Plant count: %d\n", plantChecksum)
	fmt.Println("Part 2\n------")
	// The solution here is incredibly unsatisfying and not general for any input. But it appears
	// that all the puzzle seeds result in a simple geometric repeating pattern you can detect and
	// project forward to 50B. So look at the output and sort it. In my case, the answer is 
	// score(200) + (score(200) - score(199)) * (50e9-200) == 38223 + 186*(50e9-200)
	fmt.Printf("You do the math.")
}
