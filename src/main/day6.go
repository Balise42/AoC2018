package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

var gridSize = 400

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

	grid := createVoronoi(coords)
	fmt.Println(getLargestArea(grid))
}

func getLargestArea(grid[][] int) (int, int) {
	excluded := make([]int, 0)
	excluded = append(excluded, -1)
	for i := 0; i<gridSize; i++ {
		excluded = appendIfNotExists(excluded, grid[0][i])
		excluded = appendIfNotExists(excluded, grid[gridSize-1][i])
		excluded = appendIfNotExists(excluded, grid[i][0])
		excluded = appendIfNotExists(excluded, grid[i][gridSize-1])
	}

	areas := make(map[int]int)
	for i := 0; i<gridSize; i++ {
		for j := 0; j<gridSize; j++ {
			if !contains(excluded, grid[i][j]) {
				areas[grid[i][j]]++
			}
		}
	}

	maxArea := 0
	maxAreaIndex := -1

	for k, v := range areas {
		if v > maxArea {
			maxArea = v
			maxAreaIndex = k
		}
	}

	return maxAreaIndex, maxArea
}



func appendIfNotExists(list []int, element int) []int {
	if !contains(list, element) {
		return append(list, element)
	}
	return list
}

func contains(list []int, element int) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

func parseCoordinates(line string) (int, int) {
	coords := strings.Split(line, ",")
	x, _ := strconv.ParseInt(coords[0], 10, 64)
	y, _ := strconv.ParseInt(coords[1][1:], 10, 64)
	return int(x), int(y)
}

func createVoronoi(coords [][]int) [][]int {
	grid := make([][]int, gridSize)
	for i := 0; i<gridSize; i++ {
		grid[i] = make([]int, gridSize)
		for j := 0; j<gridSize; j++ {
			grid[i][j] = -1
		}
	}

	for i := 0; i<gridSize; i++ {
		for j := 0; j<gridSize; j++ {
			grid[i][j] = findClosest(i, j, coords)
		}
	}

	return grid
}

func findClosest(x int, y int, coords[][]int) int {
	minDist := gridSize*2+1
	indexMinDist := -1

	for i, point := range coords {
		dist := distance(x, y, point)
		if dist < minDist {
			minDist = dist
			indexMinDist = i
		} else if dist == minDist && indexMinDist != i {
			indexMinDist = -1
		}
	}
	return indexMinDist
}

func distance(x int, y int, point []int) int {
	deltax := x - point[0]
	if deltax < 0 {
		deltax = -deltax
	}
	deltay := y - point[1]
	if deltay < 0 {
		deltay = -deltay
	}
	return deltax + deltay
}
