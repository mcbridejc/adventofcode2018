package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)


const GridSize = 3
const MaxCountLevel = 5000
const RequiredConsecutiveZeroCountLayers = 10

func GetInput(filepath string) (points [][2]int) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	// Read all lines and sort them alphabetically
	// (which has the ultimate effect of sorting them chronologically)
	for scanner.Scan() {
		//var point [2]int
		var x int
		var y int
		line := scanner.Text()
		fmt.Sscanf(line, "%d, %d", &x, &y)
		points = append(points, [2]int{x, y})
	}
	return
} 

func CreateGridImage(width int, height int) *image.RGBA {
	imageRect := image.Rect(0, 0, width*3+1, height*3+1)
	gridImage := image.NewRGBA(imageRect)
	
	// White out the image
	draw.Draw(gridImage, imageRect, &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.ZP, draw.Src)

	gridColor := color.RGBA{255, 255, 255, 255} //color.RGBA{10, 10, 10, 255}
	for row := 0; row <= imageRect.Dy(); row += GridSize {
		for col := 0; col < imageRect.Dx(); col += 1 {
			gridImage.SetRGBA(col, row, gridColor)	
		}
	}
	for col := 0; col <= imageRect.Dx(); col += GridSize {
		for row :=0; row < imageRect.Dy(); row += 1 {
			gridImage.SetRGBA(col, row, gridColor)
		}
	}
	return gridImage
}

func MarkNode(image *image.RGBA, x int, y int, c color.RGBA) {
	for px := x * GridSize; px < (x+1) * GridSize; px += 1 {
		for py := y * GridSize; py < (y+1)*GridSize; py += 1 {
			image.SetRGBA(px, py, c)
		}
	}
}

/* Find the convex hull of a set of 2D points using the Jarvis March / Gift Wrapping algo
*
* Essentially: Pick the point with the lower x value to start, as this is guaranteed to be 
* part of the convex hull, then try all other points as the next preferring points that 
* fall to the left of the current candidate line or that fall on the same line but are closer
*/
func JarvisConvexHull(points [][2]int) [][2]int {
	hull := make([][2]int, 0)
	currentPoint := [2]int{math.MaxInt32, math.MaxInt32}
	// Select point with smallest y
	for _, p := range points {
		if p[0] < currentPoint[0] {
			currentPoint = p
		}
	}

	for {
		hull = append(hull, currentPoint)
		var candidate [2]int
		candidate = points[0]
		for i := 1; i < len(points); i += 1 {
			A := [2]float64{float64(currentPoint[0]), float64(currentPoint[1])}
			B := [2]float64{float64(candidate[0]), float64(candidate[1])}
			P := [2]float64{float64(points[i][0]), float64(points[i][1])}

			if A == B || A == P{
				continue
			}
			// Given point P, and line A->B, the sign of the cross product of vectors
			// A->B and A->P gives which side of the line AB P falls on
			cross := (P[0] - A[0]) * (B[1] - A[1]) - (P[1] - A[1])*(B[0] - A[0])

			eps := 1e-5
			if cross < -eps {
				candidate = points[i]
			} else if math.Abs(cross) < eps {
				fmt.Println("Zero cross for ", math.Abs(cross), A, B, P)
				dCur := math.Pow(A[0] - B[0], 2) + math.Pow(A[1] - B[1], 2)
				dNew := math.Pow(A[0] - P[0], 2) + math.Pow(A[1] - P[1], 2)
				if dNew < dCur {
					candidate = points[i]
				}
			}
		}
		if candidate == hull[0] {
			break
		} else {
			currentPoint = candidate
		}
	}
	return hull
}

func min(list []int) int {
	min := math.MaxInt32
	for _, a := range(list) {
		if a < min {
			min = a
		}
	}
	return min
}
func max(list []int) int {
	max := math.MinInt32
	for _, a := range(list) {
		if a > max {
			max = a
		}
	}
	return max
}

func ContainsPoint(list [][2]int, p [2]int) bool {
	for _, test_p := range list {
		if test_p == p {
			return true
		}
	}
	return false
}

func main() {
	points := GetInput("day6_input.txt")

	xMax := math.MinInt32
	xMin := math.MaxInt32
	yMax := math.MinInt32
	yMin := math.MaxInt32

	for _, p := range points {
		if p[0] > xMax {
			xMax = p[0]
		}
		if p[0] < xMin {
			xMin = p[0]
		}
		if p[1] > yMax {
			yMax = p[1]
		}
		if p[1] < yMin {
			yMin = p[1]
		}
	}
	
	fmt.Printf("X range: [%d:%d], Y range: [%d:%d]\n", xMin, xMax, yMin, yMax)

	for i := 0; i < len(points); i += 1 {
		points[i][0] -= xMin - 5
		points[i][1] -= yMin - 5
	}
	gridImage := CreateGridImage(xMax + 10, yMax + 10)

	for _, p := range points {
		MarkNode(gridImage, p[0], p[1],  color.RGBA{0, 255, 0, 255})
	}
	
	/* Find convex hull. The points that make up the convex hull will be the "outside"
	points. All of these points will have an infinite number of points that are closest
	to them, and so are excluded from the solution.
	*/
	hull := JarvisConvexHull(points)
	for _, p := range hull {
		MarkNode(gridImage, p[0], p[1],  color.RGBA{0, 0, 0, 255})
	}

	/* Now count up cells by checking which points they are closest too. 
	There's a question of what range of cells we need to check. Although
	the number of cells is guaranteed to be bounded, cells that fall *outside* 
	of the convex hull can still be closest to points that are *inside* of it,
	which means we have to count somewhat outside of the range of the points. 
	
	The question is how far to count? I believe, in theory, in the limit as the
	distance from a point inside the hull to the edge of the hull goes to 
	zero, the distance to the furthest point that is closest to it goes to
	infinity. So this isn't really bounded (except I suppose by the fact that
	the points are integers and this quantizes the possible distances from hull)

	I'm taking the approach of breaking of counting steadily outward like onion 
	layers until no more points are counted. 
	*/

	cellCount := make([]int, len(points))
	cellLists := make([][][2]int, len(points))
	layer := 0
	cx := (xMax + xMin) / 2
	cy := (yMax + yMin) / 2
	zeroEventLayerCount := 0
	for {
		eventCount := 0
		numCells := (2 * layer + 1) * 4 - 4
		cells := make([][2]int, 0, numCells)
		for x := -1 * layer; x <= layer; x += 1 {
			cells = append(cells, [2]int{cx + x, cy - layer}) // top row
			if layer > 0 {
				cells = append(cells, [2]int{cx + x, cy + layer}) // top row
			}
		}
		if layer > 0 {
			for y := -1*layer + 1; y <= layer-1; y += 1 {
				cells = append(cells, [2]int{cx - layer, cy + y})
				cells = append(cells, [2]int{cx + layer, cy + y})
			}
		}
		for _, cell := range cells {
			x := cell[0]
			y := cell[1]
			var closestPointIdx int
			closestDistance := 50000 // effectively infinity
			for i, p := range points {
				dx := (p[0] - x)
				dy := (p[1] - y)
				//d := dx + dy // Manhattan distance
				d := dx*dx + dy*dy // squared distance is fine for comparison
				if d < closestDistance {
					closestDistance = d
					closestPointIdx = i
				}
			}
			fmt.Println("Closesst point ", points[closestPointIdx])
			if ContainsPoint(hull, points[closestPointIdx]) {
				continue
			}
			cellCount[closestPointIdx] += 1
			cellLists[closestPointIdx] = append(cellLists[closestPointIdx], [2]int{x, y})
			eventCount += 1
			zeroEventLayerCount = 0
		}
		if eventCount == 0 {
			zeroEventLayerCount += 1
			if zeroEventLayerCount > RequiredConsecutiveZeroCountLayers || layer > MaxCountLevel {
				break
			}
		}
		layer += 1
		fmt.Printf("Added %d points. Continuing to layer %d\n", eventCount, layer)
	}

	maxIdx := 0
	maxValue := 0
	for i, v := range cellCount {
		if v > maxValue {
			maxValue = v
			maxIdx = i
		}
	}

	MarkNode(gridImage, points[maxIdx][0], points[maxIdx][1],  color.RGBA{255, 0, 0, 255})

	for _, p := range cellLists[maxIdx] {
		if ContainsPoint(points, p) {
			continue // Dont overwrite original point marking
		}
		MarkNode(gridImage, p[0], p[1], color.RGBA{240, 178, 122, 180})
	}


	imageFile, err := os.Create("grid.png")
	if err != nil {
		panic(err)
	}
	png.Encode(imageFile, gridImage)

	
	fmt.Printf("Biggest region is point (%d, %d) with %d cells\n", points[maxIdx][0], points[maxIdx][1], cellCount[maxIdx])
}