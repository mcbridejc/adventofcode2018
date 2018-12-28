package main 

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Renderer struct {
	gridSize float64
	displayWidth int
	displayHeight int
	window *pixelgl.Window
	rockSprite *pixel.Sprite
	rockXforms []pixel.Matrix
	elfSprite *pixel.Sprite
	goblinSprite *pixel.Sprite
	gridDrawer *imdraw.IMDraw
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

/** Return a transform matrix to display the given sprite in the cell denoted by (x,y) */
func (ctx *Renderer) spriteXform(sprite *pixel.Sprite, x int, y int) pixel.Matrix {
	scale := ctx.gridSize / sprite.Frame().H()
	offset := pixel.V((float64(x) + 0.5) * ctx.gridSize, float64(ctx.displayHeight) - (float64(y) + 0.5) * ctx.gridSize)
	matrix := pixel.IM.Scaled(pixel.V(0, 0), scale)
	matrix = matrix.Moved(offset)
	return matrix
}

func (ctx *Renderer) Closed() bool {
	return ctx.window.Closed()
}

func (ctx *Renderer) DrawHealthBar(x, y, hitpoints int) {
	imd := imdraw.New(nil)
	x0 := float64(x)*ctx.gridSize + 1.5
	y0 := float64(ctx.displayHeight) - float64(y+1)*ctx.gridSize
	y1 := y0 + math.Ceil(float64(hitpoints)/200 * ctx.gridSize)
	if hitpoints > 120 {
		imd.Color = pixel.RGB(0.0, 1.0, 0.0)
	} else if hitpoints > 60 {
		imd.Color = pixel.RGB(1.0, 0.5, 0.0)
	} else {
		imd.Color = pixel.RGB(1.0, 0.0, 0.0)
	}
	imd.Push(pixel.V(x0, y0), pixel.V(x0, y1))
	imd.Line(3)
	imd.Draw(ctx.window)
}

func (ctx *Renderer) UpdateWorld(world *WorldMap) {
	ctx.window.Clear(pixel.RGB(0.0, 0.0, 0.0))
	for _, xform := range ctx.rockXforms {
		ctx.rockSprite.Draw(ctx.window, xform)
	}
	ctx.gridDrawer.Draw(ctx.window)
	for _, char := range world.characters {
		var sprite *pixel.Sprite
		if char.isElf {
			sprite = ctx.elfSprite
		} else {
			sprite = ctx.goblinSprite
		}
		xform := ctx.spriteXform(sprite, char.position[0], char.position[1])
		sprite.Draw(ctx.window, xform)
		ctx.DrawHealthBar(char.position[0], char.position[1], char.hitpoints)

	}
	ctx.window.Update()
}

func InitRenderer(world *WorldMap, maxWidth int, maxHeight int) *Renderer  {
	var renderer Renderer

	xGridSize := float64(maxWidth) / float64(world.width)
	yGridSize := float64(maxHeight) / float64(world.height)

	renderer.gridSize = math.Min(xGridSize, yGridSize)
	renderer.displayWidth = int(math.Ceil(renderer.gridSize * float64(world.width)))
	renderer.displayHeight = int(math.Ceil(renderer.gridSize * float64(world.height)))

	gridDrawer := imdraw.New(nil)
	gridDrawer.Color = pixel.RGB(0.3, 0.3, 0.3)
	// Draw vertical grid lines
	for i := 1; i < world.width; i++ {
		gridDrawer.Push(pixel.V(float64(i)*renderer.gridSize, 0), pixel.V(float64(i)*renderer.gridSize, float64(renderer.displayHeight)))
		gridDrawer.Line(1.1)
	}
	// Draw horizontal grid lines
	for i := 1; i < world.height; i++ {
		gridDrawer.Push(pixel.V(0, float64(i)*renderer.gridSize), pixel.V(float64(renderer.displayWidth), float64(i)*renderer.gridSize))
		gridDrawer.Line(1.1)
	}
	renderer.gridDrawer = gridDrawer

	cfg := pixelgl.WindowConfig{
		Title: "Elves vs Goblins",
		Bounds: pixel.R(0, 0, float64(renderer.displayWidth), float64(renderer.displayHeight)),
		VSync: true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	renderer.window = window

	rockImg, err := loadPicture("rock.png")
	if err != nil {
		panic(err)
	}
	renderer.rockSprite = pixel.NewSprite(rockImg, pixel.R(0, 0, 128, 128))

	elfImg, err := loadPicture("elf.png")
	if err != nil {
		panic(err)
	}
	renderer.elfSprite = pixel.NewSprite(elfImg, elfImg.Bounds())

	goblinImg, err := loadPicture("goblin.png")
	if err != nil {
		panic(err)
	}
	renderer.goblinSprite = pixel.NewSprite(goblinImg, goblinImg.Bounds())

	for x := 0; x < world.width; x++ {
		for y := 0; y < world.height; y++ {
			if world.grid[y][x].wall {
				renderer.rockXforms = append(renderer.rockXforms, renderer.spriteXform(renderer.rockSprite, x, y))
			}
		}
	}

	fmt.Printf("Added %d rocks\n", len(renderer.rockXforms))

	return &renderer
}