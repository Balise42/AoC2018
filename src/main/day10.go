package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type PointWithVelocity struct{
	X int64
	Y int64
	Dx int64
	Dy int64
}

type Canvas struct {
	Points []PointWithVelocity
}

func main() {
	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	parserRegex := regexp.MustCompile(`position=<\s?(-?\d+),\s+(-?\d+)> velocity=<\s?(-?\d+),\s+(-?\d+)>`)

	canvas := Canvas{}

	for scanner.Scan() {
		line := scanner.Text()
		coords := parserRegex.FindAllStringSubmatch(line, -1)
		canvas.Points = append(canvas.Points, newPointWithVelocity(coords[0][1:]))
	}

	iteration := int64(1)
	height := canvas.GetHeight()
	for {
		newHeight := canvas.ComputeCanvasAt(iteration).GetHeight()
		if newHeight > height {
			break
		}
		height = newHeight
		iteration++
	}

	fmt.Println(iteration-1)
	canvas.ComputeCanvasAt(iteration-1).Print()
}

func newPointWithVelocity(coords []string) PointWithVelocity {
	pv := PointWithVelocity{}
	pv.X, _ = strconv.ParseInt(coords[0], 10, 64)
	pv.Y, _ = strconv.ParseInt(coords[1], 10, 64)
	pv.Dx, _ = strconv.ParseInt(coords[2], 10, 64)
	pv.Dy, _= strconv.ParseInt(coords[3], 10, 64)
	return pv
}

func (p PointWithVelocity) ComputePointAt(time int64) PointWithVelocity {
	return PointWithVelocity{p.X + p.Dx* time, p.Y + p.Dy * time, p.Dx, p.Dy }
}

func (c Canvas) ComputeCanvasAt(time int64) Canvas {
	res := Canvas{}
	for _, p := range c.Points {
		res.Points = append(res.Points, p.ComputePointAt(time))
	}
	return res
}

func (c Canvas) GetHeight() int64 {
	_, _, minY, maxY := c.boundingBox()

	return maxY - minY
}

func (c Canvas) boundingBox() (int64, int64, int64, int64) {
	minX := int64(math.MaxInt64)
	maxX := int64(-math.MaxInt64)
	minY := int64(math.MaxInt64)
	maxY := int64(-math.MaxInt64)

	for _, p := range c.Points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return minX, maxX, minY, maxY
}

func (c Canvas) Print() {
	minX, maxX, minY, maxY := c.boundingBox()

	lines := make([][]rune, maxY - minY + 1)
	for i := int64(0); i < maxY - minY + 1; i++ {
		lines[i] = make([]rune, maxX - minX + 1)
		for j := int64(0); j < maxX - minX + 1; j++ {
			lines[i][j] = '.'
		}
	}

	for _, p := range c.Points {
		lines[p.Y - minY][p.X - minX] = '#'
	}

	for _, line := range lines {
		fmt.Println(string(line))
	}
}