package main

import (
	"fmt"
	"strconv"
)

type State struct {
	posA int64
	posB int64
	scores []int
}

func (s *State) Step() {
	nextScore := s.scores[s.posA] + s.scores[s.posB]
	if nextScore > 9 {
		s.scores = append(s.scores, nextScore / 10)
	}
	s.scores = append(s.scores, nextScore % 10)
	s.posA = (s.posA + int64(s.scores[s.posA]) + 1) % int64(len(s.scores))
	s.posB = (s.posB + int64(s.scores[s.posB]) + 1) % int64(len(s.scores))
}

func RunPart1(numRecipes int) string {
	state := State{0, 1, []int{3, 7}}
	for len(state.scores) < numRecipes + 10 {
		state.Step()
	}

	result := ""
	for i := numRecipes; i < numRecipes + 10; i++ {
		result += strconv.Itoa(state.scores[i])
	}
	return result
}

func CheckForPattern(offset int, scores []int, pattern []int) bool {
	match := true
	if offset < 0 || len(scores) - offset < len(pattern) {
		return false
	}
	for i := 0; i<len(pattern); i += 1 {
		if len(scores) < len(pattern) {
			match = false
			break
		}
		if scores[i + offset] != pattern[i] {
			match = false
			break
		}
	}
	return match
}

func RunPart2(input string) int {
	sequence := make([]int, 0)
	for i := 0; i<len(input); i++ {
		val, _ :=  strconv.Atoi(input[i:i+1])
		sequence = append(sequence, val)
	}
	fmt.Println("Sequence: ", sequence)
	state := State{0, 1, []int{3, 7}}
	for {
		state.Step()
		// Check for the pattern at two places, because we may have added 2 recipes to the scoreboard in 
		// last iteration
		if CheckForPattern(len(state.scores) - len(sequence), state.scores, sequence) {
			return len(state.scores) - len(sequence)
		}
		if CheckForPattern(len(state.scores) - len(sequence) - 1, state.scores, sequence) {
			return len(state.scores) - len(sequence) - 1
		}
	}
}

func TestPart1(numRecipes int, exp string) {
	result := RunPart1(numRecipes)
	fmt.Printf("%d -> %s: ", numRecipes, result)
	if result != exp {
		fmt.Printf("FAIL!\n")
	} else {
		fmt.Printf("PASS!\n")
	}
}

func TestPart2(input string, exp int) {
	result := RunPart2(input)
	fmt.Printf("%s -> %d: ", input, result)
	if result != exp {
		fmt.Printf("FAIL!\n")
	} else {
		fmt.Printf("PASS!\n")
	}
}


func main() {

	fmt.Println("Part 1:\n------ ")
	TestPart1(9, "5158916779")
	TestPart1(5, "0124515891")
	TestPart1(18, "9251071085")
	TestPart1(2018, "5941429882")

	// This is the random puzzle input
	input := 990941
	inputStr := "990941"
	result := RunPart1(input)
	fmt.Printf("Next 10: %s\n", result)

	fmt.Println("Part 2:\n------")
	TestPart2("51589", 9)
	TestPart2("01245", 5)
	TestPart2("92510", 18)
	TestPart2("59414", 2018)
	result2 := RunPart2(inputStr)
	fmt.Printf("Part 2 answer: %d\n", result2)
}
