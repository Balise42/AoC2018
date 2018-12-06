package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"strings"
)

type point struct {
	X int
	Y int
}

var LIMIT = 10000

func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	coords := make([][]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		x, y := parseCoords(scanner.Text())
		coords = append(coords, []int{x, y})
	}

	fmt.Println(computeSafeZoneSize(coords))
}

func parseCoords(line string) (int, int) {
	coords := strings.Split(line, ",")
	x, _ := strconv.ParseInt(coords[0], 10, 64)
	y, _ := strconv.ParseInt(coords[1][1:], 10, 64)
	return int(x), int(y)
}

func computeSafeZoneSize(coords [][]int) int {
	stack := make([]point, 1)

	//this is going to get horrendous (for some definition of horrendous that's probably 100M) memorywise. eh.
	visited := make(map[point]bool)

	stack[0] = findStart(coords)
	sum := 0

	for len(stack) > 0 {
		currPoint := stack[0]
		stack = stack[1:]
		if visited[currPoint] {
			continue
		}
		distCurrPoint := computeDistanceSum(currPoint, coords)
		if distCurrPoint < LIMIT {
			sum++
			stack = addNeighbors(currPoint, stack)
		}
		visited[currPoint] = true
	}
	return sum
}

func findStart(coords [][]int) point {
	for _, p := range coords {
		if computeDistanceSum(point{p[0], p[1]}, coords) < LIMIT {
			return point {p[0], p[1]}
		}
	}
	return  point {coords[0][0], coords[0][1]}

}

func addNeighbors(p point, points []point) []point {
	// we only need straight line neighbors, the diagonals are handled as straight line from the neighbors
	neighbors := []point{
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
	}
	return append(neighbors, points...)
}

func computeDistanceSum(p point, coords [][]int) int {
	sum := 0
	for _, coord := range coords {
		sum = sum + computeDistance(p, coord)
	}
	return sum
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