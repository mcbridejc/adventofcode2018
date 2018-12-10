package main

import (
	"bufio"
	"image/color"
	"image/gif"
	"image"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type LightPoint struct {
	x int
	y int
	vx int
	vy int
}

func ReadInput(filepath string) []*LightPoint {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("position=<\\s*(-?\\d*),\\s*(-?\\d*)> velocity=<\\s*(-?\\d*),\\s*(-?\\d*)>")
	scanner := bufio.NewScanner(f)
	points := make([]*LightPoint, 0)
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if match == nil {
			panic("No match")
		}
		p := LightPoint{}
		p.x, _ = strconv.Atoi(match[1])
		p.y, _ = strconv.Atoi(match[2])
		p.vx, _ = strconv.Atoi(match[3])
		p.vy, _ = strconv.Atoi(match[4])
		points = append(points, &p)
	}
	return points
}


// This is basically a manual search process
// In retrospect, finding when the points came closest together would have 
// found the right time, but I didn't know this going in. 
const ImageWidth = 800
const ImageHeight = 800
const SecPerFrame = 1
const StartTime = 10619
const EndTime = 10619

func DrawFrame(points []*LightPoint, time int) *image.Paletted { 
	imagePalette := []color.Color{color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}}
	img := image.NewPaletted(image.Rect(0, 0, ImageWidth, ImageHeight), imagePalette)

	
	xMax := math.MinInt32
	xMin := math.MaxInt32
	yMax := math.MinInt32
	yMin := math.MaxInt32
	for _, p := range points {
		x := + p.x + time*p.vx
		y := + p.y + time*p.vy
		if x > xMax {
			xMax = x
		}
		if x < xMin { 
			xMin = x
		}
		if y > yMax {
			yMax = y
		}
		if y < yMin {
			yMin = y
		}
	}
	expanse := 0
	if (xMax - xMin) > (yMax - yMin) {
		expanse = xMax - xMin
	} else {
		expanse = yMax - yMin
	}

	scale :=  float64(ImageWidth) / float64(expanse) * 0.9
	xOffset := -float64(xMin) + 0.05 * float64(expanse)
	yOffset := -float64(yMin) + 0.05 * float64(expanse)
	fmt.Println("Expanse: ", expanse)
	fmt.Println("offset: ", xOffset, yOffset)
	for _, p := range points {
		x := int(math.Round((float64(p.x + p.vx * time) + xOffset) * scale))
		y := int(math.Round((float64(p.y + p.vy * time) + yOffset) * scale))
		img.SetColorIndex(int(x), int(y), 1)
	}
	return img
}

func main() {
	
	points := ReadInput("day10_input.txt")
	fmt.Printf("Read %d points\n", len(points))
	
	animation := gif.GIF{}
	for time := StartTime; time <= EndTime; time += SecPerFrame {
	
		fmt.Println("Time: ", time)

		frame := DrawFrame(points, time)
		animation.Image = append(animation.Image, frame)
		animation.Delay = append(animation.Delay, 10)
	}

	// save to out.gif
    f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
    defer f.Close()
    gif.EncodeAll(f, &animation)
}