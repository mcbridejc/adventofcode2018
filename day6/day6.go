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
const MaxCountLevel = 500
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

// Get a list of points making up a square with a distance of halfwidth
// from center to edge
func SquareBorder(halfwidth int, cx int, cy int) [][2]int {
	if halfwidth == 0 {
		return [][2]int{{cx, cy}}
	}
	// Pre-allocate the needed size just for efficiency
	totalPoints := (2 * halfwidth + 1) * 4 - 4
	points := make([][2]int, totalPoints)
	for x := -1 * halfwidth; x <= halfwidth; x += 1 {
		points = append(points, [2]int{cx + x, cy - halfwidth})
		points = append(points, [2]int{cx + x, cy + halfwidth})
	}
	for y := -1 * halfwidth + 1; y <= halfwidth-1; y += 1 {
		points = append(points, [2]int{cx + halfwidth, cy + y})
		points = append(points, [2]int{cx - halfwidth, cy + y})
	}
	return points
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

func IntAbs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func IntMax(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func FindOwningPointIndex(points [][2]int, x int, y int) int {
	var closestPointIdx int
	closestDistance := math.MaxInt32
	var uniqueDistance bool
	for i, p := range points {
		dx := IntAbs((p[0] - x))
		dy := IntAbs((p[1] - y))
		d := dx + dy // Manhattan distance
		//d := dx*dx + dy*dy // squared distance is fine for comparison
		if d < closestDistance {
			closestDistance = d
			closestPointIdx = i
			uniqueDistance = true
		} else if d == closestDistance {
			uniqueDistance = false
		}
	}

	if uniqueDistance {
		return closestPointIdx
	} else {
		return -1
	}
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
	
	/* Take a square three times the size of a bounding square on the points, and 
	match its border pixels to their nearest point. These points own an infinite
	number of cells */
	cx := (xMax + xMin) / 2
	cy := (yMax + yMin) / 2
	infinite_distance := 2 * IntMax(xMax - xMin, yMax - yMin)
	infinite_points := make([]bool, len(points))

	squarePoints := SquareBorder(infinite_distance, cx, cy)
	for _, sp := range squarePoints {
		owner_idx := FindOwningPointIndex(points, sp[0], sp[1])
		if owner_idx >= 0 {
			infinite_points[owner_idx] = true
		}
	}
	
	// Mark the infinite points black for reference
	for idx, isInfinite := range infinite_points {
		if isInfinite {
			MarkNode(gridImage, points[idx][0], points[idx][1], color.RGBA{0, 0, 0, 255})
		}
	}
	
	/* Now count of all the cells accoring to their owner, excluding owners with
	infinite points */
	cellCount := make([]int, len(points))
	cellLists := make([][][2]int, len(points))
	for x := -1 * infinite_distance; x <= infinite_distance; x += 1 {
		for y := -1 * infinite_distance; y <= infinite_distance; y += 1 {
			owner_idx := FindOwningPointIndex(points, cx + x, cy + y)
			if owner_idx < 0 || infinite_points[owner_idx] {
				continue
			}
			cellCount[owner_idx] += 1
			cellLists[owner_idx] = append(cellLists[owner_idx], [2]int{cx + x, cy+y})
		}
	}

	// Find the point with the most owned cells
	maxIdx := 0
	maxValue := 0
	for i, v := range cellCount {
		if v > maxValue {
			maxValue = v
			maxIdx = i
		}
	}

	// Mark selected point
	MarkNode(gridImage, points[maxIdx][0], points[maxIdx][1],  color.RGBA{255, 0, 0, 255})

	// Fill in cells belonging to the largest region
	for _, p := range cellLists[maxIdx] {
		if ContainsPoint(points, p) {
			continue // Dont overwrite original point marking
		}
		MarkNode(gridImage, p[0], p[1], color.RGBA{240, 178, 122, 180})
	}

	imageFile, err := os.Create("grid1.png")
	if err != nil {
		panic(err)
	}
	png.Encode(imageFile, gridImage)

	fmt.Printf("Part 1\n------\n")
	fmt.Printf("Biggest region is point (%d, %d) with %d cells\n", points[maxIdx][0], points[maxIdx][1], cellCount[maxIdx])

	// PART 2
	fmt.Printf("Part 2\n------\n")

	// count up in outward layers until we dont count anymore
	layer := 0
	part2Count := 0
	part2CellList := make([][2]int, 0)
	for {
		eventCount := 0
		countPoints := SquareBorder(layer, cx, cy)
		for _, p0 := range countPoints {
			dSum := 0
			for _, p1 := range points {
				dSum += IntAbs(p1[0] - p0[0]) + IntAbs(p1[1] - p0[1])
			}
			if dSum < 10000 {
				part2Count += 1
				part2CellList = append(part2CellList, p0)
				eventCount += 1
			}
		}
		if eventCount == 0 && layer > 2*IntMax(cx, cy) {
			break
		}
		layer += 1
	}
	
	grid2Image := CreateGridImage(xMax + 10, yMax + 10) 
	for _, p := range points {
		MarkNode(grid2Image, p[0], p[1],  color.RGBA{0, 255, 0, 255})
	}
	//Fill in cells belonging to the largest region
	for _, p := range part2CellList {
		if ContainsPoint(points, p) {
			continue // Dont overwrite original point marking
		}
		MarkNode(grid2Image, p[0], p[1], color.RGBA{240, 178, 122, 180})
	}

	imageFile, err = os.Create("grid2.png")
	if err != nil {
		panic(err)
	}
	png.Encode(imageFile, grid2Image)

	fmt.Printf("Found %d points\n", part2Count)
}