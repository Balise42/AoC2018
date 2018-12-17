package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var lineRegex = regexp.MustCompile(`(x|y)=(\d+),\s(x|y)=(\d+)\.\.(\d+)`)

var size = 2000
var yMin = size
var yMax = 0


func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	ground := make([][]byte, size)

	for i := 0; i<size; i++ {
		ground[i] = make([]byte, size)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		x1, y1, x2, y2 := parseLine(line)
		if y1 < yMin {
			yMin = y1
		}
		if y2 > yMax {
			yMax = y2
		}
		for x := x1; x<=x2; x++ {
			for y := y1; y<= y2; y++ {
				ground[x][y] = '#'
			}
		}
	}

	for x := 0; x<size; x++ {
		for y := 0; y<size; y++ {
			if (ground[x][y]) != '#' {
				ground[x][y] = '.'
			}
		}
	}

	fillWater(ground, 500, yMin)

	numAccessible := 0
	numWater := 0
	for x := 0; x<size; x++ {
		for y := yMin; y<=yMax; y++ {
			if ground[x][y] == '~' || ground[x][y] == '|' {
				numAccessible++
			}
			if ground[x][y] == '~' {
				numWater++
			}
		}
	}
	fmt.Println(numAccessible, numWater)
}

func fillWater(ground [][]byte, xStart int, yStart int) {
	if ground[xStart][yStart] != '.' {
		return
	}
	blockedBottom := false
	var y int

	for y = yStart; y<=yMax; y++ {
		if ground[xStart][y+1] == '#' || ground[xStart][y+1]  == '~' {
			blockedBottom = true
			break
		}
		ground[xStart][y] = '|'
	}

	if blockedBottom {
		xLeft, blockedLeft := isBlockedLeft(ground, xStart, y)
		xRight, blockedRight := isBlockedRight(ground, xStart, y)

		for blockedLeft && blockedRight && y >= yMin {
			for x := xLeft + 1; x<xRight; x++ {
				ground[x][y] = '~'
			}
			y--
			xLeft, blockedLeft = isBlockedLeft(ground, xStart, y)
			xRight, blockedRight = isBlockedRight(ground, xStart, y)
		}

		if y >= yMin {

			for x := xLeft + 1; x < xRight; x++ {
				ground[x][y] = '|'
			}

			if !blockedLeft {
				fillWater(ground, xLeft, y)
			}

			if !blockedRight {
				fillWater(ground, xRight, y)
			}
		}
	}
}

func isBlockedLeft(ground [][]byte, xStart int, y int) (int, bool) {
	blockedLeft := false
	var xLeft int
	for x := xStart-1; x>=0; x-- {
		if ground[x][y] == '#' {
			blockedLeft = true
			xLeft = x
			break
		} else if ground[x][y+1] == '.' || ground[x][y+1] == '|' {
			blockedLeft = false
			xLeft = x
			break
		}
	}
	return xLeft, blockedLeft
}

func isBlockedRight(ground [][]byte, xStart int, y int) (int, bool) {
	blockedRight := false
	var xRight int
	for x := xStart+1; x<=size; x++ {
		if ground[x][y] == '#' {
			blockedRight = true
			xRight = x
			break
		} else if ground[x][y+1] == '.' || ground[x][y+1] == '|' {
			blockedRight = false
			xRight = x
			break
		}
	}
	return xRight, blockedRight
}

func parseLine(line string) (int, int, int, int) {
	toks := lineRegex.FindStringSubmatch(line)
	var x1, x2, y1, y2 int
	if toks[1] == "x" {
		x, _ := strconv.Atoi(toks[2])
		x1 = x
		x2 = x
		y1, _ = strconv.Atoi(toks[4])
		y2, _ = strconv.Atoi(toks[5])
	} else {
		y, _ := strconv.Atoi(toks[2])
		y1 = y
		y2 = y
		x1, _ = strconv.Atoi(toks[4])
		x2, _ = strconv.Atoi(toks[5])
	}
	return x1, y1, x2, y2
}

func printGround(ground [][]byte) {
	for y := yMin; y<yMax; y++ {
		for x := 420; x<580; x++ {
			fmt.Print(string(ground[x][y]))
		}
		fmt.Println("")
	}
}