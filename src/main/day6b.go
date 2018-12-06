package main

import (
	"os"
	"bufio"
)

type point struct {
	X int
	Y int
}

func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	coords := make([][]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		x, y := parseCoordinates(scanner.Text())
		coords = append(coords, []int{x, y})
	}

	computeSafeZoneSize(coords)
}

func computeSafeZoneSize(coords [][]int) int {
	stack := make([]point, 1)
	visited := make(map[point]bool)

	stack[0] = point{coords[0][0], coords[0][1]}

	sum := computeDistanceSum(stack[0], coords)
	for len(stack) > 0 {
		currPoint := stack[0]
	}
}

func computeDistanceSum(p point, coords [][]int) int {
	sum := 0
	for _, coord := range coords {
		sum = sum + computeDistance(p, coord)
	}
}
func computeDistance(p point, coords []int) int {
	deltaX := p.X - coords[0]
	if deltaX < 0 {
		deltaX = -deltaX
	}
	deltaY := p.Y - coords[1]
	if deltaY < 0 {
		deltaY = -deltaY
	}
	return deltaX + deltaY
}