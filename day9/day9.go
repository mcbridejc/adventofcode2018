package main

import (
	"fmt"
)

type Marble struct {
	ccw *Marble
	cw *Marble
	value int
}

func NewMarbleCircle() *Marble {
	var marble Marble
	marble.ccw = &marble
	marble.cw = &marble
	return &marble
}

// Get the marble N positions clockwise from target marble
func (m *Marble) SeekCw(n int) (*Marble) {
	var ret *Marble
	ret = m
	for ; n > 0 ; n -= 1 {
		ret = m.cw
	}
	return ret
}

// Get the marble N positions clockwise from target marble
func (m *Marble) SeekCcw(n int) (*Marble) {
	var ret *Marble
	ret = m
	for ; n > 0 ; n -= 1 {
		ret = ret.ccw
	}
	return ret
}

// Insert a new marble on the clockwise side of m, and return the new marble
func (m *Marble) Insert(value int) (*Marble) {
	cw := m.cw
	// inserting between m and cw
	var newMarble Marble
	newMarble.value = value
	
	m.cw = &newMarble
	cw.ccw = &newMarble
	newMarble.cw = cw
	newMarble.ccw = m
	return &newMarble
}

// Remove the marble from the ring, and return it's clockwise neighbor
func (m *Marble) Remove() (*Marble) {
	ccw := m.ccw
	cw := m.cw
	m.ccw = nil
	m.cw = nil
	ccw.cw = cw
	cw.ccw = ccw
	return cw
}

func PrintCircle(currentMarble *Marble) {
	fmt.Printf("%d ", currentMarble.value)
	m := currentMarble.SeekCw(1)
	for ; m != currentMarble; m = m.SeekCw(1) {
		fmt.Printf("%d ", m.value)
	}
	fmt.Printf("\n")
}

func PlayMarbles(numPlayers int, lastMarble int) (highScore int) {
	// Initialize a circular linked list with one zero marble in it
	currentMarble := NewMarbleCircle()
	playerScores := make([]int, numPlayers)
	nextPlayer := 0
	
	for nextMarbleValue := 1; nextMarbleValue <= lastMarble; nextMarbleValue += 1 {
		if nextMarbleValue%23 == 0 {
			// Seek 7 positions clockwise, and remove that marble
			// Add the current marble and the removed marble to the players score
			playerScores[nextPlayer] += nextMarbleValue
			currentMarble = currentMarble.SeekCcw(7)
			playerScores[nextPlayer] += currentMarble.value
			currentMarble = currentMarble.Remove()
		} else {
			// Seek 1 position clockwise, insert a new marble clockwise of that, 
			// and make the new marble active
			currentMarble = currentMarble.SeekCw(1).Insert(nextMarbleValue)
		}
		nextPlayer = (nextPlayer+1)%numPlayers
		if(nextMarbleValue < 10) {
		}
	}

	for _, score := range playerScores {
		if score > highScore {
			highScore = score
		}
	}
	return highScore
}


func TestMarbles(numPlayers int, lastMarble int, expHighScore int) {
	highScore := PlayMarbles(numPlayers, lastMarble)
	fmt.Printf("%d players; last marble is worth %d: ", numPlayers, lastMarble)
	if highScore == expHighScore {
		fmt.Printf("PASS\n")
	} else {
		fmt.Printf("FAILED (got %d, expected %d)\n", highScore, expHighScore)
	}
}

func main() {
	/* Run examples as test cases */
	fmt.Println("Examples: ")
	TestMarbles(10, 1618, 8317)
	TestMarbles(13, 7999, 146373)
	TestMarbles(17, 1104, 2764)
	TestMarbles(21, 6111, 54718)
	TestMarbles(30, 5807, 37305)

	
	highScore := PlayMarbles(412, 71646)
	fmt.Println("\nPart 1 answer: ", highScore)

	highScore = PlayMarbles(412, 7164600)
	fmt.Println("\nPart 2 answer: ", highScore)
}