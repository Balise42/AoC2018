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

	seenIt := make(map[string]int)

	period := -1
	firstIter := 0

	numIters := 1000000000

	for i := 1; i<1000; i++ {
		forest.evolve()
		if seenIt[forest.flatten()] != 0 {
			period = i - seenIt[forest.flatten()]
			firstIter = seenIt[forest.flatten()]
			break
		}
		seenIt[forest.flatten()] = i
	}

	equalIter := firstIter + (numIters - firstIter) % period

	forest = NewForest(lines, height, width)
	for i := 0; i<equalIter; i++ {
		forest.evolve()
	}
	fmt.Println(forest.computeValue())

}

func NewForest(lines []string, height int, width int) Forest {
	forest := Forest {make([][]byte, width), width, height }
	for x := 0; x < width; x++ {
		forest.desc[x] = make([]byte, height)
		for y := 0; y<height; y++ {
			forest.desc[x][y] = lines[y][x]
		}
	}
	return forest
}

func (f* Forest) evolve() {
	newState := make([][]byte, f.width)

	for i := 0; i<f.width; i++ {
		newState[i] = make([]byte, f.height)
	}

	for x := 0; x<f.width; x++ {
		for y := 0; y < f.height; y++ {
			tree, lumber := f.getNeighbors(x, y)
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
					newState[x][y] = '|'
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
	f.desc = newState
}

func (f Forest) computeValue() int64 {
	wooded := int64(0)
	lumber := int64(0)
	for x := 0; x<f.width; x++ {
		for y := 0; y<f.height; y++ {
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
func (f * Forest) getNeighbors(x int, y int) (int, int) {

	var dxmin, dxmax, dymin, dymax int
	if x > 0 {
		dxmin = -1
	} else {
		dxmin = 0
	}

	if y > 0 {
		dymin = -1
	} else {
		dymin = 0
	}

	if x < f.width - 1 {
		dxmax = 1
	} else {
		dxmax = 0
	}

	if y < f.height - 1 {
		dymax = 1
	} else {
		dymax = 0
	}

	tree := 0
	lumber := 0
	for dx := dxmin; dx <= dxmax; dx++ {
		for dy := dymin; dy<= dymax; dy++ {
			if dx != 0 || dy != 0 {
				if f.desc[x+dx][y+dy] == '#' {
					lumber++
				}
				if f.desc[x+dx][y+dy] == '|' {
					tree++
				}
			}
		}
	}
	return tree, lumber
}

func (f Forest) printForest() {
	for y := 0; y<f.height; y++ {
		for x := 0; x<f.width; x++ {
			fmt.Print(string(f.desc[x][y]))
		}
		fmt.Println("")
	}
}
func (f Forest) flatten() string {
	res := ""
	for x := 0; x<f.width; x++ {
		for y := 0; y<f.height; y++ {
			res = res + string(f.desc[x][y])
		}
	}
	return res
}
