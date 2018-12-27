package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"time"
	"sort"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
) 

type Track int
const (
	Blank Track = 0
	Horizontal Track = 1
	Vertical Track = 2
	RightCurve Track = 3
	LeftCurve Track = 4
	Intersection Track = 5
)

type CartDir int
const (
	Up CartDir = 0
	Right CartDir = 1
	Down CartDir = 2
	Left CartDir = 3
)

type Cart struct {
	x int
	y int
	dir CartDir
	turnCount int
	scheduledForRemoval bool
}

type Map struct {
	locs [][]Track
	width int
	height int
}

func NewMap() (*Map) {
	var m Map
	m.locs = make([][]Track, 0)
	return &m
}

func (m *Map) Set(x, y int, shape Track) {
	// Resize storage as necessary to accomodate set location
	if x >= m.width {
		for i := 0; i < m.height; i++ {
			m.locs[i] = append(m.locs[i], make([]Track, x - m.width + 1)...)
		}
		m.width = x + 1
	}
	if y >= m.height {
		linesToAdd := y - m.height + 1
		for i := 0; i<linesToAdd; i++ {
			m.locs = append(m.locs, make([]Track, m.width))
		}
		m.height = y + 1
	}

	m.locs[y][x] = shape
}

func (m *Map) Get (x, y int)(Track) {
	shape := m.locs[y][x]
	return shape
}

type CartList struct {
	carts []*Cart
}

func (a CartList) Len() int { 
	return len(a.carts) 
}
func (list CartList) Less(i, j int) bool { 
	a := list.carts[i]
	b := list.carts[j]
	if a.y == b.y {
		return a.x < b.x
	} else {
		return a.y < b.y
	}
}
func (a CartList) Swap(i, j int) { 
	a.carts[i], a.carts[j] = a.carts[j], a.carts[i] 
}
func (a *CartList) Remove(c *Cart) {
	index := -1
	for i, check := range a.carts {
		if c == check {
			index = i
			break
		}
	}
	if index > -1 {
		a.carts = append(a.carts[:index], a.carts[index+1:]...)
		// a.carts[index:] = a.carts[index+1:]  // shift left
		// a.carts = a.carts[:len(a.carts)-1] // truncate last element
	}
}

func ReadInput(filepath string) (*Map, CartList) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	m := NewMap()
	carts := make([]*Cart, 0)
	line := 0
	for scanner.Scan() {
		for i, s := range scanner.Text() {
			if s == '>' {
				carts = append(carts, &Cart{i, line, Right, 0, false})
				s = '-'
			} else if s == '<' {
				carts = append(carts, &Cart{i, line, Left, 0, false})
				s = '-'
			} else if s == '^' {
				carts = append(carts, &Cart{i, line, Up, 0, false})
				s = '|'
			} else if s == 'v' {
				carts = append(carts, &Cart{i, line, Down, 0, false})
				s = '|'
			}

			var trackType Track
			if s == '|' {
				trackType = Vertical
			} else if s == '-' {
				trackType = Horizontal
			} else if s == '/' {
				trackType = RightCurve 
			} else if s == '\\' {
				trackType = LeftCurve
			} else if s == '+' {
				trackType = Intersection
			}
			m.Set(i, line, trackType)
		}
		line += 1
	}
	return m, CartList{carts}
}

func MoveCart(tracks *Map, c *Cart) {
	t := tracks.Get(c.x, c.y)
	switch c.dir {
	case Up:
		if t == Horizontal {
			panic(fmt.Sprintf("Kart going up on horizontal track @ (%d, %d)", c.x, c.y))
		} else if t == Vertical {
			c.y -= 1
		} else if t == RightCurve {
			c.dir = Right
			c.x += 1
		} else if t == LeftCurve {
			c.dir = Left
			c.x -= 1
		} else if t == Intersection {
			switch c.turnCount % 3 {
			case 0: // Turn left
				c.dir = Left
				c.x -= 1
			case 1: // Go straigt
				c.y -= 1
			case 2: // Turn Right
				c.dir = Right
				c.x += 1
			}
			c.turnCount += 1
		}
	case Right:
		if t == Horizontal {
			c.x += 1
		} else if t == Vertical {
			panic(fmt.Sprintf("Kart going right on vertical track @ (%d, %d)", c.x, c.y))
		} else if t == RightCurve {
			c.dir = Up
			c.y -= 1
		} else if t == LeftCurve {
			c.dir = Down
			c.y += 1
		} else if t == Intersection {
			switch c.turnCount % 3 {
			case 0: // Turn left
				c.dir = Up
				c.y -= 1
			case 1: // Go straigt
				c.x += 1
			case 2: // Turn Right
				c.dir = Down
				c.y += 1
			}
			c.turnCount += 1
		}
	case Down:
		if t == Horizontal {
			panic(fmt.Sprintf("Kart going down on horizontal track @ (%d, %d)", c.x, c.y))
		} else if t == Vertical {
			c.y += 1
		} else if t == RightCurve {
			c.dir = Left
			c.x -= 1
		} else if t == LeftCurve {
			c.dir = Right
			c.x += 1
		} else if t == Intersection {
			switch c.turnCount % 3 {
			case 0: // Turn left
				c.dir = Right
				c.x += 1
			case 1: // Go straigt
				c.y += 1
			case 2: // Turn Right
				c.dir = Left
				c.x -= 1
			}
			c.turnCount += 1
		}
	case Left:
		if t == Horizontal {
			c.x -= 1
		} else if t == Vertical {
			panic(fmt.Sprintf("Kart going left on vertical track @ (%d, %d)", c.x, c.y))
		} else if t == RightCurve {
			c.dir = Down
			c.y += 1
		} else if t == LeftCurve {
			c.dir = Up
			c.y -= 1
		} else if t == Intersection {
			switch c.turnCount % 3 {
			case 0: // Turn left
				c.dir = Down
				c.y += 1
			case 1: // Go straigt
				c.x -= 1
			case 2: // Turn Right
				c.dir = Up
				c.y -= 1
			}
			c.turnCount += 1
		}
	}
}

func RunTick(tracks *Map, carts *CartList) (collision bool) {
	collision = false

	sort.Sort(carts)

	collidedCarts := make([]*Cart, 0)
	for i, c := range carts.carts {
		MoveCart(tracks, c)
		for j, other := range carts.carts {
			if i == j {
				continue
			}
			if c.x == other.x && c.y == other.y && !other.scheduledForRemoval {
				fmt.Printf("Collision @ %d,%d\n", c.x, c.y)
				c.scheduledForRemoval = true
				other.scheduledForRemoval = true
				collidedCarts = append(collidedCarts, c)
				collidedCarts = append(collidedCarts, other)
				collision = true
			}
		}
	}
	for _, c := range collidedCarts {
		carts.Remove(c)
	}
	return collision
}

func DrawTrackSegment(imd *imdraw.IMDraw, x, y float64, size float64, shape Track) {
	switch(shape) {
	case Horizontal:
		imd.Push(pixel.V(x, y+size/2))
		imd.Push(pixel.V(x+size, y+size/2))
		imd.Line(size/4)
	case Vertical:
		imd.Push(pixel.V(x+size/2, y))
		imd.Push(pixel.V(x+size/2, y+size))
		imd.Line(size/4)
	case Intersection:
		imd.Push(pixel.V(x+size/2, y))
		imd.Push(pixel.V(x+size/2, y+size))
		imd.Line(size/4)
		imd.Push(pixel.V(x, y+size/2))
		imd.Push(pixel.V(x+size, y+size/2))
		imd.Line(size/4)
	case RightCurve:
		imd.Push(pixel.V(x+size/2, y))
		imd.Push(pixel.V(x+size, y+size/2))
		imd.Line(size/4)
		imd.Push(pixel.V(x, y+size/2))
		imd.Push(pixel.V(x+size/2, y+size))
		imd.Line(size/4)
	case LeftCurve:
		imd.Push(pixel.V(x+size/2, y))
		imd.Push(pixel.V(x, y+size/2))
		imd.Line(size/4)
		imd.Push(pixel.V(x+size/2, y+size))
		imd.Push(pixel.V(x+size, y+size/2))
		imd.Line(size/4)
	}
	return
}

func DrawTrack(tracks *Map, size float64) (*imdraw.IMDraw) {
	gridSize := size / 150.0
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 0.5, 0.5)
	imd.EndShape = imdraw.RoundEndShape
	for y := 0.0; y < float64(tracks.height); y += 1.0 {
		for x := 0.0; x < float64(tracks.width); x += 1.0 {
			shape := tracks.Get(int(x), int(y))
			DrawTrackSegment(imd, x*gridSize, size - y*gridSize, gridSize, shape)
		}
	}
	return imd
}

func DrawCarts(carts []*Cart, size int) *imdraw.IMDraw {
	gridSize := float64(size) / 150.0
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(0.5, 0.5, 1.0)
	for _, cart := range carts {
		x := float64(cart.x) * gridSize
		y := float64(size) - float64(cart.y)*gridSize
		switch cart.dir {
		case Up:
			imd.Push(pixel.V(x, y), pixel.V(x+gridSize, y), pixel.V(x+gridSize/2, y+gridSize))
		case Right:
			imd.Push(pixel.V(x, y), pixel.V(x, y+gridSize), pixel.V(x+gridSize, y+gridSize/2))
		case Down:
			imd.Push(pixel.V(x, y+gridSize), pixel.V(x+gridSize, y+gridSize), pixel.V(x+gridSize/2, y))
		case Left:
			imd.Push(pixel.V(x+gridSize, y), pixel.V(x+gridSize, y+gridSize), pixel.V(x, y+gridSize/2))
		}
		imd.Polygon(0)
	}
	return imd
}

func entry() {
	interactive := flag.Bool("interactive", false, "Run GUI in step-by-step mode")
	inputFile := flag.String("file", "day13_input.txt", "The input file")
	flag.Parse()
	
	fmt.Println("Reading input from ", *inputFile)
	tracks, carts := ReadInput(*inputFile)
	fmt.Printf("Size of map: %dx%d\n", tracks.width, tracks.height)
	fmt.Printf("Number of carts: %d\n", len(carts.carts))

	if *interactive {
		// Create a UI for visualizing
		cfg := pixelgl.WindowConfig{
			Title:  "Cart Crash",
			Bounds: pixel.R(0, 0, 1050, 1050),
		}
		window, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		tick := 0
		imd := DrawTrack(tracks, 1050)
		window.Clear(pixel.RGB(1.0, 1.0, 1.0))
		//imd.Draw(window)
		for (!window.Closed()) {
			if window.JustPressed(pixelgl.KeyEnter) {
				tick += 1
				fmt.Println("Iteration ", tick)
				RunTick(tracks, &carts)
				window.Clear(pixel.RGB(1.0, 1.0, 1.0))
				imd.Draw(window)
				cartImd := DrawCarts(carts.carts, 1050)
				cartImd.Draw(window)
				window.Update()
			} else {
				//time.Sleep(time.Duration(0.5*1e9))
			}
			window.Update()
		}
	} else {
		for tick := 0; tick < 30000; tick += 1 {
			collision := RunTick(tracks, &carts)
			if collision { fmt.Printf("Iteration: %d, carts remaining: %d\n", tick+1, len(carts.carts)) }
			if len(carts.carts) == 1 {
				//RunTick(tracks, &carts)
				fmt.Printf("Last cart remains at %d, %d\n", carts.carts[0].x, carts.carts[0].y)
				return
			}
		}
	}
	
}

func main() {
	// Run via pixelGL so it can hold "the original thread" for OS/UI interactions
	// It will run our main code
	pixelgl.Run(entry)
}