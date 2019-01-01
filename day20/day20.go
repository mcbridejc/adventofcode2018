package main

import (
	"fmt"
	"io/ioutil"
 	"strings"
)

type Room struct {
	N *Room
	S *Room
	E *Room
	W *Room
	distance int
	visited bool
}

type Coordinate struct {
	x int
	y int
}


/** Split options separated by |, respecting any nested option groups */
func SplitOptions(sequence string) []string {
	groups := make([]string, 0)
	startIdx := 0
	idx := 0
	parenDepth := 0
	for idx < len(sequence) {
		if parenDepth == 0 && sequence[idx] == '|' {
			groups = append(groups, sequence[startIdx:idx])
			startIdx = idx + 1
		} else if sequence[idx] == '(' {
			parenDepth++
		} else if sequence[idx] == ')' {
			parenDepth--
		}
		idx++
	}
	// Add the last group if there is one (last character could have been |,
	// in which case there will be no last group)
	if startIdx < len(sequence) {
		groups = append(groups, sequence[startIdx:])
	}
	return groups
}

/** Consumes the next part of the sequence
* This may return: 
* 1. A zero length array, indicating the end of the sequence has been reached
* 2. A array of length 1, indicating a single series of movements (no alternatives)
* 3. An array of length > 1, indicating there are multiple options for the next segment
	 In this case, you must call ConsumeSequence again on each element, as it may contain
	 nested options (e.g. '(NEE|(WEN|NWE))' will return {'NEE', '(WEN|NWE)'}
*/
func ConsumeSequence(sequence *string) []string {
	idx := 0
	var ret []string
	if len(*sequence) == 0 {
		ret = make([]string, 0)
	} else if (*sequence)[0] == '(' {
		// Consume everything up to the paired closing paren
		openCount := 1
		for openCount > 0 {
			idx++
			if (*sequence)[idx] == '(' {
				openCount++
			}
			if (*sequence)[idx] == ')' {
				openCount--
			}
		}
		// the character at 0 is open paren, and at idx is the closing paren.
		// We want everything inbetween
		optionSeq := (*sequence)[1:idx]
		// Remove everything up to and including the closing paren
		*sequence = (*sequence)[idx+1:]
		// Now split it on '|' and return each option
		ret = SplitOptions(optionSeq)
	} else {
		// Consume everything up to the end or up to an open paren
		for idx < len(*sequence) {
			if (*sequence)[idx] == '(' {
				break
			}
			idx++
		}
		retSeq := (*sequence)[:idx]
		*sequence = (*sequence)[idx:]
		ret = []string{retSeq}
	}
	// if len(ret) > 0 {
	// 	fmt.Printf("Returning subseq ['%s']\n", strings.Join(ret, "', '"))
	// } else {
	// 	fmt.Println("Break")
	// }
	
	return ret
}

func FindOrCreateRoom(atlas map[Coordinate]*Room, loc Coordinate) *Room {
	room, ok := atlas[loc]
	if !ok {
		room = &Room{}
		atlas[loc] = room
	}
	return room
}

func WalkPath(atlas map[Coordinate]*Room, start Coordinate, sequence string) (endPos Coordinate) {
	curPos := start
	curRoom := atlas[start]
	for {
		// Get the next steps or set of options. Returned segments are removed from sequence
		nextSegment := ConsumeSequence(&sequence)
		if len(nextSegment) == 0 {
			break
		} else if len(nextSegment) > 1 {
			var endPos Coordinate
			for _, s := range nextSegment {
				// Each of the paths will end up in the same location
				endPos = WalkPath(atlas, curPos, s)
			}
			curPos = endPos
			curRoom = atlas[curPos]
		} else {
			for _, char := range nextSegment[0] {
				switch(char) {
				case 'N': 
					curPos.y--
					nextRoom := FindOrCreateRoom(atlas, curPos)
					curRoom.N = nextRoom
					nextRoom.S = curRoom
				case 'E': 
					curPos.x++
					nextRoom := FindOrCreateRoom(atlas, curPos)
					curRoom.E = nextRoom
					nextRoom.W = curRoom
				case 'S': 
					curPos.y++
					nextRoom := FindOrCreateRoom(atlas, curPos)
					curRoom.S = nextRoom
					nextRoom.N = curRoom
				case 'W': 
					curPos.x--
					nextRoom := FindOrCreateRoom(atlas, curPos)
					curRoom.W = nextRoom
					nextRoom.E = curRoom
				}
				curRoom = atlas[curPos]
			}
		}
	}
	return curPos // return position after navigating the sequence
}

// Walk down the tree of rooms, annotating the distance when it is first reach
// Return the longest distance
func AnnotateDistances(room *Room, distance int) int {
	// If we reach a room that has already been visited from a shorter path, we
	// dont need to follow its exits
	if room.visited && room.distance <= distance {
		return distance - 1
	}
	neighbors := make([]*Room, 0, 4)
	if room.N != nil { neighbors = append(neighbors, room.N)}
	if room.S != nil { neighbors = append(neighbors, room.S)}
	if room.E != nil { neighbors = append(neighbors, room.E)}
	if room.W != nil { neighbors = append(neighbors, room.W)}

	maxDistance := 0
	room.distance = distance
	room.visited = true
	if len(neighbors) == 0 {
		return distance
	}
	for _, n := range neighbors {
		branchDist := AnnotateDistances(n, distance + 1)
		if branchDist > maxDistance {
			maxDistance = branchDist
		}
	}
	return maxDistance
}

func main() {
	loadFile := "day20_input.txt"

	directionBytes, err := ioutil.ReadFile(loadFile)
	if err != nil {
		panic(err)
	}
	directions := string(directionBytes)
	// Remove any trailing newline
	directions = strings.TrimSuffix(directions, "\n")
	// Remove the first and last character (the ^ and $) because they carry no meaning
	directions = directions[1:len(directions)-1]
	fmt.Printf("Read directions %d long\n", len(directions))

	atlas := make(map[Coordinate]*Room)
	start := Coordinate{0, 0}
	atlas[start] = &Room{}
	WalkPath(atlas, start, directions)

	fmt.Printf("Atlas now has %d rooms\n", len(atlas))

	maxDistance := AnnotateDistances(atlas[start], 0)
	fmt.Println("Maximum distance found: ", maxDistance)

	part2count := 0
	for _, room := range atlas {
		if room.distance >= 1000 {
			part2count++
		}
	}
	fmt.Println("The number of rooms with a distance >= 1000 is ", part2count)
}