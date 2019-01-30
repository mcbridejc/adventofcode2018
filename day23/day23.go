package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Nanobot struct {
	x int
	y int
	z int
	r int
}
func ReadInput(filepath string) []Nanobot {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	result := make([]Nanobot, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var x, y, z, r int
		fmt.Sscanf(scanner.Text(), "pos=<%d,%d,%d>, r=%d\n", &x, &y, &z, &r);
		result = append(result, Nanobot{x, y, z, r})
	}
	return result
}

func IntAbs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func CubeDistance(cx int, cy int, cz int, cubeSize int, x int, y int, z int) int {
	x0 := cx - cubeSize / 2
	x1 := cx + cubeSize / 2
	y0 := cy - cubeSize / 2
	y1 := cy + cubeSize / 2
	z0 := cz - cubeSize / 2
	z1 := cz + cubeSize / 2
	
	distance := 0
	if x < x0 {
		distance += x0 - x
	} else if x > x1 {
		distance += x - x1
	}
	if y < y0 {
		distance += y0 - y
	} else if y > y1 {
		distance += y - y1
	}
	if z < z0 {
		distance += z0 - z
	} else if z > z1 {
		distance += z - z1
	}
	return distance
}

func CubeScore(x, y, z, size int, bots []Nanobot) int {
	score := 0
	for _, b := range bots {
		d := CubeDistance(x, y, z, size, b.x, b.y, b.z)
		if d <= b.r {
			score += 1
		}
	}
	return score
}

type Score struct {
	x int
	y int
	z int
	size int
	score int
}

/** Break up the volume into cubes, and score each cube based on how many bots 
are within range *of the cube*, i.e. of any single location within the cube 
volume. The maximum score for any location within the cube must be less than or
equal to the cube score. 

The returned list is sorted so that the highest score is at element 0. 
*/
func ScaledSearch(minX, maxX, minY, maxY, minZ, maxZ, scale int, bots []Nanobot) []Score {
	scores := make([]Score, 0)
	for x := minX; x <= maxX; x += scale {
		for y := minY; y <= maxY; y += scale {
			for z:= minZ; z <= maxZ; z += scale {
				score := CubeScore(x, y, z, scale, bots)
				scores = append(scores, Score{x, y, z, scale, score})
			}
		}
	}
	// Sort by score
	sort.Slice(scores, func (i, j int) bool { return scores[i].score > scores[j].score})
	return scores
}



func Recurse(score Score, maxValue int, bots []Nanobot, collectLocations bool) (int, []Score) {
	const factor = 32
	// For each volume, break it up into factor^3 sub volumes, and score
	// each of these. maxValue provides the current known best score, so any 
	// volume with a score less than this can be ignored
	// force even
	size := (score.size / factor) - ((score.size / factor) % 2)
	if size <= 4 {
		size = 1
	}
	
	//fmt.Printf("Searching Area of size %d with sub-size %d\n", score.size, size)

	contenders := make([]Score, 0)
	
	subScores := ScaledSearch(
		score.x - score.size/2, 
		score.x + score.size/2,
		score.y - score.size/2,
		score.y + score.size/2,
		score.z - score.size/2,
		score.z + score.size/2,
		size,
		bots)
	
	if size == 1 {
		// If we're looking at individual points, make our list
		for _, s := range subScores {
			if s.score > maxValue {
				// We're increasing max value, throw away any points we've already collected
				fmt.Printf("Increasing maxValue from %d to %d\n", maxValue, s.score)
				maxValue = s.score
				if collectLocations {
					contenders = make([]Score, 0)
					contenders = append(contenders, s)
				}
			}  else if collectLocations && s.score == maxValue {
				// We're equal to the max, so include this point and any others we've already found
				contenders = append(contenders, s)
				//fmt.Printf("Up to %d contenders\n", len(contenders))
			} else {
				//skip
			}
		}
		if len(contenders) > 0 {
			c := contenders[0]
			fmt.Printf("Found %d contenders in subregion (%d, %d, %d), score %d\n", len(contenders), score.x, score.y, score.z, c.score)
		}
		return maxValue, contenders
	}
	
	for _, s := range subScores {
		if collectLocations && s.score >= maxValue || s.score > maxValue {
			// This sub-volume score is >= maxValue, so it *may* contain points to include
			newMaxValue, subContenders := Recurse(s, maxValue, bots, collectLocations)
			if len(subContenders) > 0 {
				if newMaxValue > maxValue {
					// If we increase the maxvalue, we need to clear our contender list
					contenders = make([]Score, 0)
				}
				contenders = append(contenders, subContenders...)
			}
			maxValue = newMaxValue
		}
	}
	return maxValue, contenders
}


/*
* The algorithm I settled on is essentially to subdivide the total volume into sub-volumes,
* and score each one. We can then go through the sub-volumes pyramid style, drilling down 
* to a single point. It's important that we choose the sub-volumes with the highest score
* first, because these are the most likely to contain high score points. Once we find a 
* point with a certain score, we can safely ignore any sub-volume whose score is less than
* that, as it cannot contain any points with a score >= the subvolume score. 
*/
func main() {
	bots := ReadInput("day23_input.txt")

	var largestBot Nanobot
	for _, b := range bots {
		if b.r > largestBot.r {
			largestBot = b
		}
	}

	inRangeCount := 0
	minX := math.MaxInt32
	minY := math.MaxInt32
	minZ := math.MaxInt32
	maxX := math.MinInt32
	maxY := math.MinInt32
	maxZ := math.MinInt32
	meanX := 0
	meanY := 0
	meanZ := 0
	for _, b := range bots {
		d := IntAbs(b.x - largestBot.x) + IntAbs(b.y - largestBot.y) + IntAbs(b.z - largestBot.z)
		meanX += b.x
		meanY += b.y
		meanZ += b.z

		if d <= largestBot.r {
			inRangeCount += 1
		}
		if b.x > maxX {
			maxX = b.x
		}
		if b.y > maxY {
			maxY = b.y
		}
		if b.z > maxZ {
			maxZ = b.z
		}
		if b.x < minX {
			minX = b.x
		}
		if b.y < minY {
			minY = b.y
		}
		if b.z < minZ {
			minZ = b.z
		}
	}
	meanX = meanX / len(bots)
	meanY = meanY / len(bots)
	meanZ = meanZ / len(bots)

	fmt.Printf("Number in range for part 1: %d\n", inRangeCount)
	
	
	fmt.Printf("Number range x: (%d, %d), y: (%d, %d), z: (%d, %d)\n", minX, maxX, minY, maxY, minZ, maxZ)

	// The score at the center is a good heuristic to allow us to skip any 
	// sub volumes with a lower score, so it will be used as the starting "maxValue"
	centerScore := CubeScore(meanX, meanY, meanZ, 0, bots)
	fmt.Printf("Center: (%d, %d, %d), score: %d\n", meanX, meanY, meanZ, centerScore)
	

	// Do pyramid style search
	scores := ScaledSearch(minX, maxX, minY, maxY, minZ, maxZ, 4*1024*1024, bots)
	fmt.Println("Top 20 scores:")
	for i, s := range scores {
		fmt.Println(s.score)
		if i > 20 {
			break
		} 
	}

	topSize := 0
	if maxX - minX > topSize {
		topSize = maxX - minX
	}
	if maxY - minY > topSize {
		topSize = maxY - minY
	}
	if maxZ - minZ > topSize {
		topSize = maxZ - minZ
	}
	topVolume := Score{(minX + maxX)/2, (minY + maxY)/2, (minZ + maxZ)/2, topSize, 0}
	fmt.Printf("Initial volume: (%d, %d, %d) size: %d\n", topVolume.x, topVolume.y, topVolume.z, topVolume.size)
	// Do two passes: 
	// - On first pass, only collect the top score
	// - On second pass, collect all the locations that achieve the top score
	// This is because the algorithm will waste much too much time collecting the many, many
	// locations that meet a lower score, before it figures out there are higher scores
	topScore, _ := Recurse(topVolume, centerScore, bots, false)
	fmt.Println("Top score: ", topScore)
	topScore, locations := Recurse(topVolume, topScore, bots, true)
	fmt.Printf("Found %d locations with score %d\n", len(locations), topScore)
	for _, l := range locations {
		fmt.Println("Location: ", l)
		// Compute the distance to origin
		distance := CubeDistance(0, 0, 0, 0, l.x, l.y, l.z)
		// re-compute the score...just as a sanity check
		recomputeScore := CubeScore(l.x, l.y, l.z, 0, bots)
		fmt.Println("Distance from origin: ", distance)
		fmt.Println("recomputed score: ", recomputeScore)
	}
}