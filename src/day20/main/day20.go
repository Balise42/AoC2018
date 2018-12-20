package main

import (
	"os"
	"io/ioutil"
	"fmt"
)

type XY struct {
	X int
	Y int
}

type Node struct {
	Coords *XY
	N *Node
	S *Node
	E *Node
	W *Node
}

func main() {
	filename := os.Args[1]
	exp, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	maze := createMaze(exp[1:len(exp)-1])

	visited := make(map[*XY]int)
	visit(visited, maze[XY{0,0}], 0)

	maxDist := 0
	for _, dist := range visited {
		if dist > maxDist {
			maxDist = dist
		}
	}
	fmt.Println(maxDist)

	numRooms := 0
	for _, dist := range visited {
		if dist >= 1000 {
			numRooms++
		}
	}

	fmt.Println(numRooms)
}

func createMaze(s []byte) map[XY]*Node {
	maze := make(map[XY]*Node)
	zero := XY{0,0}
	firstNode := Node{&zero, nil, nil, nil, nil}
	maze[zero] = &firstNode
	explore(maze, s, []*Node{&firstNode}, false)

	return maze
}

func visit(visited map[*XY]int, start * Node, dist int) {
	if start == nil {
		return
	}

	currDist, ok := visited[start.Coords]
	if !ok || currDist > dist {
		visited[start.Coords] = dist
		visit(visited, start.N, dist+1)
		visit(visited, start.S, dist+1)
		visit(visited, start.W, dist+1)
		visit(visited, start.E, dist+1)
	}
}


func explore(maze map[XY]*Node, s []byte, nodes []*Node, sub bool) []*Node {
	if len(s) == 0 {
		return nodes
	}

	if sub {
		branches := findBranches(s)
		newNodes := make([]*Node, 0)
		for _, b := range branches {
			newNodes = append(newNodes, explore(maze, b, nodes, false)...)
		}
		return newNodes
	}


	if s[0] == '(' {
		closing := findClosing(s)
		newNodes := explore(maze, s[1:closing], nodes, true)
		return explore(maze, s[closing+1:], newNodes, false)
	}

	var newNodes []*Node

	if s[0] == 'N' {
		newNodes = exploreOneStep(maze, nodes, 0, -1)
	}
	if s[0] == 'E' {
		newNodes = exploreOneStep(maze, nodes, 1, 0)
	}
	if s[0] == 'S' {
		newNodes = exploreOneStep(maze, nodes, 0, 1)
	}
	if s[0] == 'W' {
		newNodes = exploreOneStep(maze, nodes, -1, 0)
	}

	return explore(maze, s[1:], newNodes, false)
}
func findBranches(bytes []byte) [][]byte {
	branches := make([][]byte, 0)
	currStart := 0
	currGroup := 0
	for i, c := range bytes {
		if c == '(' {
			currGroup++
		} else if c == ')' {
			currGroup--
		} else if c == '|' && currGroup == 0 {
			branches = append(branches, bytes[currStart:i])
			currStart = i + 1
		}
	}
	if currStart < len(bytes) {
		branches = append(branches, bytes[currStart:])
	}
	return branches
}

func exploreOneStep(maze map[XY]*Node, nodes []*Node, deltaX int, deltaY int) []*Node {
	newNodes := make([]*Node, 0)
	for _, curr := range nodes {
		xy := XY{curr.Coords.X + deltaX, curr.Coords.Y + deltaY}
		n, ok := maze[xy]
		if !ok {
			n = &Node{&xy, nil, nil, nil, nil}
			maze[xy] = n
			if deltaX == 0 {
				if deltaY == -1 {
					curr.N = n
					n.S = curr
				} else {
					curr.S = n
					n.N = curr
				}
			} else {
				if deltaX == -1 {
					curr.W = n
					n.E = curr
				} else {
					curr.E = n
					n.W = curr
				}
			}
		}
		newNodes = append(newNodes, n)
	}
	return newNodes
}

func findClosing(s []byte) int {
	numPar := 0
	for i, b := range s {
		if b == '(' {
			numPar++
		} else if b == ')' {
			numPar--
			if numPar == 0 {
				return i
			}
		}
	}
	return -1
}





