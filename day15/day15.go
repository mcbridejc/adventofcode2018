package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
	// "github.com/faiface/pixel"
	// "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const AttackDamage = 3

type Direction int
const (
	None Direction = iota
	North
	East
	South
	West
)

func Dir2Str(dir Direction) string {
	switch dir {
	case North:
		return "North"
	case South:
		return "South"
	case East: 
		return "East"
	case West:
		return "West"
	case None:
		return "None"
	default:
		return "INVALID"
	}
}

type Character struct {
	position [2]int
	isElf bool // Elf if true, goblin if false
	hitpoints int
	awaitingMove bool
}

type GridCell struct {
	wall bool
	occupant *Character
}

type WorldMap struct {
	grid [][]GridCell
	width int
	height int
	characters []*Character
	turnCount int
}

func (world *WorldMap) Copy() *WorldMap {
	var copy WorldMap
	for row := 0; row <world.height; row++ {
		for col := 0; col < world.width; col++ {
			copy.SetCell(col, row, world.grid[row][col].wall)
		}
	}
	for _, c := range world.characters {
		// Make a copy, and store the pointer
		newChar := *c
		copy.characters = append(copy.characters, &newChar)
	}
	copy.turnCount = world.turnCount
	return &copy
}

func (world *WorldMap) SetCell (x, y int, wall bool) {
	// Resize storage as necessary to accomodate set location
	if x >= world.width {
		for i := 0; i < world.height; i++ {
			world.grid[i] = append(world.grid[i], make([]GridCell, x - world.width + 1)...)
		}
		world.width = x + 1
	}
	if y >= world.height {
		linesToAdd := y - world.height + 1
		for i := 0; i<linesToAdd; i++ {
			world.grid = append(world.grid, make([]GridCell, world.width))
		}
		world.height = y + 1
	}

	world.grid[y][x] = GridCell{wall, nil}
}


func (world *WorldMap) AddCharacter (x, y int, isElf bool) {
	const StartingHitpoints = 200
	char := Character{[2]int{x, y}, isElf, StartingHitpoints, false}
	world.characters = append(world.characters, &char)
	world.grid[y][x].occupant = &char
}

func (m *WorldMap) EmptyNeighbors(x, y int) ([][2]int, []Direction) {
	coords := make([][2]int, 0, 4)
	directions := make([]Direction, 0, 4)
	if x > 0 {
		coords = append(coords, [2]int{x-1, y})
		directions = append(directions, West)
	}
	if y > 0 {
		coords = append(coords, [2]int{x, y-1})
		directions = append(directions, North)
	}
	if x < m.width - 1 {
		coords = append(coords, [2]int{x+1, y})
		directions = append(directions, East)
	}
	if y < m.height - 1 {
		coords = append(coords, [2]int{x, y+1})
		directions = append(directions, South)
	}

	emptyCoords := make([][2]int, 0, 4)
	emptyDirections := make([]Direction, 0, 4)
	for i, c := range coords {
		cell := m.grid[c[1]][c[0]]
		if cell.wall { continue }
		if cell.occupant != nil { continue }
		emptyCoords = append(emptyCoords, c)
		emptyDirections = append(emptyDirections, directions[i])
	}
	return emptyCoords, emptyDirections
}

func (m *WorldMap) Neighbors(x, y int) [][2]int {
	coords := make([][2]int, 0, 4)
	if x > 0 {
		coords = append(coords, [2]int{x-1, y})
	}
	if y > 0 {
		coords = append(coords, [2]int{x, y-1})
	}
	if x < m.width - 1 {
		coords = append(coords, [2]int{x+1, y})
	}
	if y < m.height - 1 {
		coords = append(coords, [2]int{x, y+1})
	}
	return coords
}

func (world *WorldMap) InRange(x, y int, attackerIsElf bool) bool {
	neighbors := world.Neighbors(x, y)
	for _, p := range neighbors {
		cell := world.grid[p[1]][p[0]]
		if cell.occupant != nil {
			if attackerIsElf != cell.occupant.isElf {
				return true
			}
		}
	}
	return false // No enemy found in neighbors
}

func (world *WorldMap) SortCharacters() {
	sort.Slice(world.characters, func(i, j int) bool {
		a := world.characters[i]
		b := world.characters[j]
		if a.position[1] == b.position[1] {
			return a.position[0] < b.position[0]
		} else {
			return a.position[1] < b.position[1]
		}
	})
}

func (world *WorldMap) MoveCharacter(char *Character, dir Direction) {
	if dir == None {
		return
	}
	x := char.position[0]
	y := char.position[1]
	world.grid[y][x].occupant = nil
	switch dir {
	case North:
		y -= 1
	case South:
		y += 1
	case East:
		x += 1
	case West: 
		x -= 1
	}
	world.grid[y][x].occupant = char
	char.position = [2]int{x, y}
}

func (world *WorldMap) ElfCount() int {
	elfCount := 0
	for _, c := range world.characters {
		if c.isElf {
			elfCount++
		}
	}
	return elfCount
}

func (world *WorldMap) CheckForWinner() (score int, finished bool, winnerIsElf bool) {
	elfCount := 0
	goblinCount := 0
	for _, c := range world.characters {
		if c.isElf {
			elfCount++
		} else {
			goblinCount++
		}
	}
	if elfCount == 0 || goblinCount == 0 {
		finished = true
		totalHP := 0
		// Check for corner case: has the current turn been completed? 
		completedTurnCount := world.turnCount
		if !world.IsTurnComplete() {
			completedTurnCount--
		}
		fmt.Printf("Number of completed turns: %d", completedTurnCount)
		for _, c := range world.characters {
			fmt.Println(c.position, c.hitpoints)
			totalHP += c.hitpoints
		}
		score = completedTurnCount * totalHP
		if goblinCount == 0 {
			winnerIsElf = true
		}
	}
	return score, finished, winnerIsElf
}

func (world *WorldMap) IsTurnComplete() bool {
	completed := true
	for _, c := range world.characters {
		if c.awaitingMove {
			completed = false
			break
		}
	}
	return completed
}

func (world *WorldMap) KillCharacter(char *Character) {
	// Remove the pointer in the grid cell
	world.grid[char.position[1]][char.position[0]].occupant = nil
	// Remove the character from the list
	for i, check := range world.characters {
		if check == char {
			world.characters = append(world.characters[:i], world.characters[i+1:]...)
			break
		}
	}
}

func (world *WorldMap) GetNeighboringCell(x int, y int, dir Direction) *GridCell {
	switch dir {
	case North:
		y -= 1
	case South:
		y += 1
	case East:
		x += 1
	case West: 
		x -= 1
	}
	return &world.grid[y][x]
}

func (world *WorldMap) MakeNextMove(elfBonus int) {
	// Character list should be sorted by position already
	// Find the next character that hasn't been moved this turn
	var char *Character
	for _, c := range world.characters {
		if c.awaitingMove {
			char = c
			break
		}
	}
	// If no characters found, clear their flags and resort
	if char == nil {
		fmt.Println("Re-sorting characters")
		world.SortCharacters()
		for _, c := range world.characters {
			c.awaitingMove = true
		}
		world.turnCount++
		char = world.characters[0]
	}
	
	char.awaitingMove = false
	chosenDir := ChooseDirection(world, char)
	world.MoveCharacter(char, chosenDir)
	// Check neighbor cells for target
	targets := make(map[Direction]*Character)
	for _, dir := range []Direction{North, West, East, South} {
		cell := world.GetNeighboringCell(char.position[0], char.position[1], dir)
		if cell.occupant != nil && cell.occupant.isElf != char.isElf {
			targets[dir] = cell.occupant
		}
	}
	// Find lowest hitpoint target (if there's a tie we may have to break it 
	// with direction priority
	minHP := 200
	for _, t := range targets {
		if t.hitpoints < minHP {
			minHP = t.hitpoints
		}
	}
	selectTargets := make(map[Direction]*Character)
	for dir, t := range targets {
		if t.hitpoints == minHP {
			selectTargets[dir] = t
		}
	}
	var finalTarget *Character
	for _, dir := range []Direction{North, West, East, South} {
		if t, present := selectTargets[dir]; present {
			finalTarget = t
			break
		}
	}
	if finalTarget != nil {
		finalTarget.hitpoints -= AttackDamage
		if char.isElf {
			finalTarget.hitpoints -= elfBonus
		}
		if finalTarget.hitpoints <= 0 {
			world.KillCharacter(finalTarget)
		}
	}
}

type PointSet map[[2]int]bool

func Expand(world *WorldMap, visited PointSet, borderCoords [][2]int) (newBorderCoords [][2]int) {
	newBorderCoords = make([][2]int, 0)
	for _, p := range borderCoords {
		emptyNeighbors, _ := world.EmptyNeighbors(p[0], p[1])
		for _, e := range emptyNeighbors {
			// Skip if its already visited
			if _, present := visited[e]; present {
				continue
			}
			newBorderCoords = append(newBorderCoords, e)
		}
	}
	return newBorderCoords
}

// Return true if a comes before b reading order
func ReadingCompare(a, b [2]int) bool {
	if a[1] == b[1] {
		return a[0] < b[0]
	} else {
		return a[1] < b[1]
	}
}

func CheckForAnyInRange(world *WorldMap, coordMap map[Direction][][2]int, isElf bool) map[Direction][2]int {
	result := make(map[Direction][2]int)
	for dir, points := range coordMap {
		for _, p := range points {
			if world.InRange(p[0], p[1], isElf) {
				curPoint, found := result[dir]
				if !found || (found && ReadingCompare(p, curPoint)) {
					result[dir] = [2]int{p[0], p[1]}		
				}
			}
		}
	}
	return result
}

func ChooseDirection(world *WorldMap, char *Character) Direction {
	// The algorithm is essentially to move outward along all possible paths until
	// an in-range cell is found, or all cells are exhausted. 
	// The only thing we care about is the next move, so we only need to find which 
	// of the neighbor squares leads to a shorter path
	// A common set of already visited points is kept, because if we reach a cell 
	// via some route, then any path that passes through later from one of 
	// the other starting directions must result in a longer path to any location

	// First off, check if we are already in range
	x, y := char.position[0], char.position[1]
	if world.InRange(x, y, char.isElf) {
		return None
	}

	// Initialize a list of border coordinates
	// Each neighbor that is empty created a border coordinate
	coords, directions := world.EmptyNeighbors(x, y)
	if len(coords) == 0 {
		// No neighboring cells are free, so we can't move
		return None
	}
	borderCoordsMap := make(map[Direction][][2]int)
	for i := 0; i<len(coords); i++ {
		borderCoordsMap[directions[i]] = [][2]int{coords[i]}
	}

	routeFoundMap := CheckForAnyInRange(world, borderCoordsMap, char.isElf)
	visitedSet := make(PointSet)

	for ; len(routeFoundMap) == 0; {
		for dir, borderCoords := range borderCoordsMap {
			borderCoordsMap[dir] = Expand(world, visitedSet, borderCoords)
		}
		
		numBorderCells := 0
		for _, v := range borderCoordsMap {
			numBorderCells += len(v)
		}
		if numBorderCells == 0 {
			//fmt.Println("Ran out of cells and didn't find any attack positions")
			return None
		}
		// Mark as visited so we don't re-visit in any future expand round
		for _, borderCoords := range borderCoordsMap {
			for _, p := range borderCoords {
				visitedSet[p] = true
			}
		}
		routeFoundMap = CheckForAnyInRange(world, borderCoordsMap, char.isElf)
	}

	// When we get here, we broke because at least one initial direction 
	// was expanded into a "in-range" cell. If there are more than one, we 
	// want to choose them based on priority: North > West > East > South
	// Problem description calls this "Reading order"
	bestPoint := [2]int{0, 0}
	selectedDir := None
	for _, dir := range []Direction{North, West, East, South} {
		if point, present := routeFoundMap[dir]; present {
			if selectedDir == None || ReadingCompare(point, bestPoint) {
				selectedDir = dir
				bestPoint = point
			}
		}
	}
	return selectedDir
}

func ReadWorld(filepath string) *WorldMap {
	world := WorldMap{}
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	row := 0
	for scanner.Scan() {
		for i, s := range scanner.Text() {
			wall := s == '#'
			world.SetCell(i, row, wall)
			if s == 'G' {
				world.AddCharacter(i, row, false) // Add goblin
			} else if s == 'E' {
				world.AddCharacter(i, row, true) // Add Elf
			}
		}
		row++
	}
	return &world
}

func entry() {
	inputFile := flag.String("file", "day15_input.txt", "The input file")
	rate := flag.Float64("rate", 1.0, "The number of ticks per second playrate")
	replay := flag.Bool("replay", false, "Run in replay mode: simulate everything, and allow stepping through history in GUI")
	part2 := flag.Bool("part2", false, "Compute part 2")
	flag.Parse()
	
	world := ReadWorld(*inputFile)
	renderer := InitRenderer(world, 1050, 1050)

	if *replay {
		worldSeries := make([]*WorldMap, 0)
		for {
			// Make the first move (because the current turn is complete)
			world.MakeNextMove(0)
			for !world.IsTurnComplete() {
				world.MakeNextMove(0)
				_, finished, _ := world.CheckForWinner()
				if finished {
					break
				}
			}
			fmt.Println("Completed turn")
			worldSeries = append(worldSeries, world.Copy())
			score, finished, _ := world.CheckForWinner()
			if(finished) {
				frame := 0
				fmt.Printf("The war is over! The final score is %d\n", score)
				fmt.Println("Use arrow keys to step through the fight")
				for !renderer.Closed() {
					
					if renderer.window.JustPressed(pixelgl.KeyRight) {
						if frame < len(worldSeries) - 1 {
							frame++
							fmt.Println("Advancing to frame ", frame)
						} else {
							fmt.Println("You've reached the last frame!")
						}
					} else if renderer.window.JustPressed(pixelgl.KeyLeft) {
						if frame > 0 {
							frame--
							fmt.Println("Rewinging to frame ", frame)
						} else {
							fmt.Println("You're on the first frame!")
						}
					}
					renderer.UpdateWorld(worldSeries[frame])
					// block forever so the GUI stays active. Let user close after reviewing
					time.Sleep(time.Duration(0.1*float64(time.Second)))
				}
				return
			}
		}
	} else if *part2 {
		fmt.Println("Solving part 2: elf bonus")
		world = ReadWorld(*inputFile)
		initialElfCount := world.ElfCount()
		elfBonus := 2
		var score int
		for {
			finished := false
			elfBonus++
			fmt.Println("Trying elf bonus ", elfBonus)
			world = ReadWorld(*inputFile)
				
			for !finished{
				world.MakeNextMove(elfBonus)
				if world.ElfCount() < initialElfCount {
					// an elf died. Abort early
					break
				}
				score, finished, _ = world.CheckForWinner()
			}
			if finished {
				break
			}
		}
		fmt.Println("Required bonus is %d, score is %d\n", elfBonus, score)
	} else {
		return
		tickPeriod := time.Duration(float64(time.Second) / *rate) 
		nextUpdate := time.Now().Add(tickPeriod)
		for !renderer.Closed() {
			if time.Now().Before(nextUpdate) {
				time.Sleep(nextUpdate.Sub(time.Now()))
			}
			nextUpdate = time.Now().Add(tickPeriod)
			world.MakeNextMove(0)
			score, finished, _ := world.CheckForWinner()
			if(finished) {
				fmt.Printf("The war is over! The final score is %d\n", score)
				for {
					// block forever so the GUI stays active. Let user close after reviewing
					time.Sleep(1.0)
				}
			}
			renderer.UpdateWorld(world)
		}
	}
}

func main() {
	// Run via pixelGL so it can hold "the original thread" for OS/UI interactions
	// It will run our main code
	pixelgl.Run(entry)
}