package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	p [4]int
	constellation *Constellation
}

type Constellation struct {
	points []*Point
}

func (c *Constellation) Add(p *Point) { 
	c.points = append(c.points, p)
}

func ReadInput(filepath string) []*Point {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	points := make([]*Point, 0)
	
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		p := Point{}
		fmt.Sscanf(scanner.Text(), "%d,%d,%d,%d", &p.p[0], &p.p[1], &p.p[2], &p.p[3])
		points = append(points, &p)
	}
	return points
}

func IntAbs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}

// Return manhattan distance between two points
func Distance(a *Point, b *Point) int {
	sum := 0
	for i := 0; i < 4; i++ {
		sum += IntAbs(a.p[i] - b.p[i])
	}	
	return sum
}

func BuildConstellations(points []*Point) []*Constellation {
	constellations := make([]*Constellation, 0)
	for _, p0 := range points {
		if p0.constellation == nil {
			c := Constellation{[]*Point{p0}}
			constellations = append(constellations, &c)
			p0.constellation = &c
		}
		for _, p1 := range points {
			if Distance(p0, p1) <= 3 {
				if p1.constellation == nil {
					// add to p0's constellation
					p1.constellation = p0.constellation
					p0.constellation.Add(p1)
				} else if p1.constellation != p0.constellation {
					// Move all from p1's constellation to p0's constellation
					c0 := p0.constellation
					c1 := p1.constellation
					for _, move_point := range c1.points {
						c0.Add(move_point)
						move_point.constellation = c0
					}
					c1.points = make([]*Point, 0)
				}
			}
		}
	}
	trimmed := make([]*Constellation, 0)
	for _, c := range constellations {
		if len(c.points) > 0 {
			trimmed = append(trimmed, c)
		}
	}
	return trimmed
}


func main() {
	points := ReadInput("day25_input.txt")
	fmt.Printf("Read %d points\n", len(points))
	constellations := BuildConstellations(points)
	fmt.Printf("Found %d constellations\n", len(constellations))
	// for _, c := range constellations {
	// 	fmt.Println(len(c.points))
	// }
}