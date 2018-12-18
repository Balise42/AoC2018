package main

import (
	"bufio"
	"fmt"
	"os"
)

type Forest struct {
	desc [][]byte
	width int
	height int
}

func main() {
	height := 50
	width := 50

	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	forest := NewForest(lines, height, width)

	for i := 0; i<10; i++ {
		forest.evolve()
	}

	fmt.Println(forest.computeValue())
}

func NewForest(lines []string, height int, width int) Forest {
	forest := Forest {make([][]byte, width), width, height }
	//TODO
	return forest
}

func (f* Forest) evolve() {
	newState := make([][]byte, f.width)

	for i := 0; i<f.width; i++ {
		newState[i] = make([]byte, f.height)
	}

	for x := 0; x<f.width; x++ {
		for y := 0; y < f.height; y++ {
			tree, lumber := f.getNeighbors()
			if f.desc[x][y] == '.' {
				if tree >= 3 {
					newState[x][y] = '|'
				} else {
					newState[x][y] = '.'
				}
			}
			if f.desc[x][y] == '|' {
				if lumber >= 3 {
					newState[x][y] = '#'
				} else {
					newState[x][y] = '.'
				}
			}
			if f.desc[x][y] == '#' {
				if lumber == 0 || tree == 0 {
					newState[x][y] = '.'
				} else {
					newState[x][y] = '#'
				}
			}
		}
	}

}

func (f Forest) computeValue() int64 {
	wooded := int64(0)
	lumber := int64(0)
	for x := 0; x<f.width; x++ {
		for y := 0; y<f.height; x++ {
			if f.desc[x][y] == '#' {
				lumber++
			}
			if f.desc[x][y] == '|' {
				wooded++
			}
		}
	}
	return wooded * lumber
}
