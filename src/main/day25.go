package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

type FixPoint struct {
	X int
	Y int
	Z int
	T int
}

func main() {
	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	fixPoints := make([]*FixPoint, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		toks := strings.Split(line, ",")
		x, _ := strconv.Atoi(toks[0])
		y, _ := strconv.Atoi(toks[1])
		z, _ := strconv.Atoi(toks[2])
		t, _ := strconv.Atoi(toks[3])
		fixPoints = append(fixPoints, &FixPoint{x, y, z, t})
	}

	fmt.Println(getNumComponents(fixPoints))
}

func getNumComponents(fixPoints []*FixPoint) int {
	graph := make(map[*FixPoint][]*FixPoint)
	visited := make(map[*FixPoint]int)
	for _, u := range fixPoints {
		graph[u] = make([]*FixPoint, 0)
		visited[u] = -1
	}
	for i := 0; i<len(fixPoints); i++ {
		for j := i+1; j<len(fixPoints); j++ {
			if isNeighbor(fixPoints[i], fixPoints[j]) {
				graph[fixPoints[i]] = append(graph[fixPoints[i]], fixPoints[j])
				graph[fixPoints[j]] = append(graph[fixPoints[j]], fixPoints[i])
			}
		}
	}

	numComp := 0
	for u, c := range visited {
		if c == -1 {
			numComp++
			visit(u, graph, numComp, visited)
		}
	}

	return numComp
}

func visit(point *FixPoint, graph map[*FixPoint][]*FixPoint, numComp int, visited map[*FixPoint]int) {
	toVisit := []*FixPoint{point}
	for len(toVisit) > 0 {
		u := toVisit[0]
		toVisit = toVisit[1:]
		if visited[u] == -1 {
			toVisit = append(toVisit, graph[u]...)
			visited[u] = numComp
		}
	}
}



func isNeighbor(u *FixPoint, v *FixPoint) bool {
	deltaX := getDelta(u.X, v.X)
	deltaY := getDelta(u.Y, v.Y)
	deltaZ := getDelta(u.Z, v.Z)
	deltaT := getDelta(u.T, v.T)
	return deltaX + deltaY + deltaZ + deltaT <= 3

}

func getDelta(i int, j int) int {
	delta := i - j
	if delta < 0 {
		return -delta
	}
	return delta
}
