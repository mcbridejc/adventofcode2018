package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func Annihilate(a string, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b) && b != a 
}

// Inefficient/naive version that removes pairs and then re-evaluates every pair in the new
// string until no more pairs are found
func ReactSlowRecursive(s string) string {
	for {
		num_cuts := 0
		new_s := ""
		for i:=0; i<len(s); i+=1 {
			if i == len(s) - 1 {
				new_s = new_s + s[len(s)-1:]
				continue
			}

			a := s[i:i+1]	
			b := s[i+1:i+2]
			if Annihilate(a, b) {
				i += 1 // skip both
				num_cuts += 1
				continue
			} else {
				new_s = new_s + s[i:i+1]
			}
		}

		if num_cuts == 0 {
			break
		}

		s = new_s
	}
	return s
}

// Much smarter version that iterates once through the string, and only removes
// from the end of the result string to avoid re-allocating and copying
func React(s string) string {
	new_s := s[0:1]
	i := 1
	for i < len(s) {
		if len(new_s) == 0 {
			new_s = s[i:i+1]
			i += 1
			continue
		}
		a := new_s[len(new_s)-1:]
		b := s[i:i+1]
		if Annihilate(a, b) {
			new_s = new_s[:len(new_s)-1]
		} else {
			new_s = new_s + s[i:i+1]
		}
		i += 1
	}
	return new_s
}

func main() {
	data, err := ioutil.ReadFile("day5_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(data)
	input = strings.Trim(input, "\n")
	s := React(input)

	fmt.Println("Final sequence:\n", s)
	fmt.Printf("%d characters remain\n", len(s))

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	min_length := len(input)
	for i:=0; i<len(alphabet); i += 1 {
		s = strings.Replace(input, alphabet[i:i+1], "", -1)
		s = strings.Replace(s, strings.ToUpper(alphabet[i:i+1]), "", -1)
		s = React(s)
		if len(s) < min_length {
			min_length = len(s)
		}
	}
	fmt.Println("Min length: ", min_length)
}