package main

import (
	_ "flag"
	"fmt"
	"math"
	_ "os"
	_ "runtime/pprof"
)

// Problem input
const DEPTH = 8787
const TARGET_X = 10
const TARGET_Y = 725

// Example input
// const DEPTH = 510
// const TARGET_X = 10
// const TARGET_Y = 10

func NewCave(width, height int) [][]int64 {
	c := make([][]int64, 0)
	for x := 0; x < width; x++ {
		c = append(c, make([]int64, height))
	}
	return c
}

type Tool int;
const (
	None Tool = iota
	Climb
	Torch
)

type GraphNode struct {
	// each node represents a state defined by the location, and the tool 
	// equipped
	x int
	y int
	tool Tool

	distance int // The shorted time found so far to reach this node
	edges map[*GraphNode]int // node is the key, cost of traversing is value 
}

type NodeState struct {
	x int
	y int
	tool Tool
}

type NodeCollection map[NodeState]*GraphNode;

func (array *NodeCollection) FindOrCreate(x, y int, tool Tool) *GraphNode {
	state := NodeState{x, y, tool}
	foundNode, found := (*array)[state]
	if found {
		return foundNode
	}
	n := GraphNode{}
	n.x = x
	n.y = y
	n.tool = tool
	n.distance = math.MaxInt32
	n.edges = make(map[*GraphNode]int)
	(*array)[state] = &n
	return &n
}

func ToolAllowed(gridType int, tool Tool) bool {
	switch gridType {
	case 0:
		return tool == Torch || tool == Climb
	case 1:
		return tool == None || tool == Climb
	case 2:
		return tool == None || tool == Torch
	default:
		return false
	}
}

const SWITCH_TIME = 7
const MOVE_TIME = 1

// Traverse outward, assuming we've just reached node after a time 
// distance has elapsed
func AnnotateDistance(distance int, node *GraphNode) {
	distanceQ := make([]int, 0)
	nodeQ := make([]*GraphNode, 0)
	
	distanceQ = append(distanceQ, distance)
	nodeQ = append(nodeQ, node)
	idx := 0
	for idx < len(distanceQ) {
		distance = distanceQ[idx]
		node = nodeQ[idx]
		idx++
		if distance >= node.distance {
			continue
		}
		node.distance = distance
		for next, cost := range node.edges {
			distanceQ = append(distanceQ, distance+cost)
			nodeQ = append(nodeQ, next)
		}
	}
}

func main() {

	// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

    // flag.Parse()
    // if *cpuprofile != "" {
    //     f, err := os.Create(*cpuprofile)
    //     if err != nil {
    //         panic(err)
    //     }
    //     pprof.StartCPUProfile(f)
    //     defer pprof.StopCPUProfile()
    // }

	// For part 2, we need to compute the map beyond the target, as this may be 
	// part of the fastest route. This amount of extra is a total SWAG. 
	width := TARGET_X * 5
	height := TARGET_Y * 2
	cave := NewCave(width, height)

	// Logic assumes TARGET location cannot fall on first row/col
	if TARGET_X == 0 || TARGET_Y == 0 {
		panic("Unhandled target position")
	}

	for x := 1; x < width; x++ {
		cave[x][0] = (int64(x) * 16807)
	}
	for y := 1; y < height; y++ {
		cave[0][y] = (int64(y) * 48271)
	}

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			if x != TARGET_X || y != TARGET_Y {
				geoA := (cave[x-1][y] + DEPTH) % 20183
				geoB := (cave[x][y-1] + DEPTH) % 20183
				cave[x][y] = ((geoA % 20183)  * geoB) % 20183
			}
		}
	}

	// Convert to region type
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cave[x][y] = ((cave[x][y] + DEPTH) % 20183) % 3
		}
	}
	// Sum the "risk" for part 1
	risk := 0
	for y := 0; y <= TARGET_Y; y++ {
		for x := 0; x <= TARGET_X; x++ {
			category := cave[x][y]
			switch category {
			case 0:
				fmt.Printf(".")
			case 1:
				fmt.Printf("=")
			case 2:
				fmt.Printf("|")
			default:
				fmt.Println("Unknown category ", category, cave[x][y], x, y)
				panic("Unknown value")
			}
			risk += int(category)
			cave[x][y] = category
		}
		fmt.Printf("\n")
	}

	fmt.Println("Part 1\n------")
	fmt.Println("Risk: ", risk)

	// Build a graph of all possible states
	graphNodes := make(NodeCollection, 0)
	startNode := graphNodes.FindOrCreate(0, 0, Torch)
	targetNode := graphNodes.FindOrCreate(TARGET_X, TARGET_Y, Torch)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gridType := int(cave[x][y])

			nodes := make([]*GraphNode, 0, 2)

			for _, tool := range []Tool{None, Climb, Torch} {
				if ToolAllowed(gridType, tool) {
					nodes = append(nodes, graphNodes.FindOrCreate(x, y, tool))
				}
			}
			// Note: len(nodes) will always be 2
			// Add edges for switching tools within this room
			nodes[0].edges[nodes[1]] = SWITCH_TIME
			nodes[1].edges[nodes[0]] = SWITCH_TIME
			
			// Find allowable neighbors and add edges for those
			for _, n := range nodes {
				for _, delta := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
					nX := n.x + delta[0]
					nY := n.y + delta[1]
					if nX < 0 || nX >= width || nY < 0 || nY >= height {
						continue
					}
					nType := int(cave[nX][nY])
					if ToolAllowed(nType, n.tool) {
						// Add edge to transition to neighboring room with same tool
						nNode := graphNodes.FindOrCreate(nX, nY, n.tool)
						n.edges[nNode] = MOVE_TIME
					}
				}
			}
		}
	}
	fmt.Printf("Graph constructed\n")
	fmt.Printf("Created %d nodes\n", len(graphNodes))

	// Now we have a graph. Traverse it to populate all nodes with a distance from start. 
	AnnotateDistance(0, startNode)
	fmt.Printf("It took %d seconds to reach target\n", targetNode.distance)
}