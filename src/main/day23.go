package main

import (
	"os"
	"bufio"
	"regexp"
	"strconv"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Nanobot struct {
	X     int
	Y     int
	Z     int
	Range int
}

func main() {
	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	nanobots := make([]*Nanobot, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		nanobots = append(nanobots, parseNanobots(scanner.Text()))
	}

	maxRange := 0
	var largestRangeNanobot *Nanobot = nil
	for _, n := range nanobots {
		if n.Range > maxRange {
			maxRange = n.Range
			largestRangeNanobot = n
		}
	}

	numBotsinRange := 0
	for _, n := range nanobots {
		if distNanobots(n, largestRangeNanobot) <= largestRangeNanobot.Range {
			numBotsinRange++
		}
	}

	fmt.Println(numBotsinRange)

	fmt.Println(computeMaxBotsPointDistance(nanobots))
}

func random(i int, j int) int {
	return int(rand.Int63n(int64(j)-int64(i))) + i
}


func computeMaxBotsPointDistance(nanobots []*Nanobot) int {
	minX, maxX, minY, maxY, minZ, maxZ, _, _ := getBounds(nanobots)
	rand.Seed(time.Now().UnixNano())

	point := &Nanobot{random(minX, maxX), random(minY, maxY), random(minZ, maxZ), 0}
	currValue := numIntersectionsCenter(point, nanobots)
	for i := 0; i < 5000000; i++ {
		tmpPoint := &Nanobot{random(minX, maxX), random(minY, maxY), random(minZ, maxZ), 0}

		tmpValue := numIntersectionsCenter(tmpPoint, nanobots)
		if tmpValue > currValue {
			point = tmpPoint
			currValue = tmpValue
		}
	}
	fmt.Println(point)
	fmt.Println(currValue)

	closestPoint := getClosestPoint(point, currValue, nanobots)
	return distNanobots(closestPoint, &Nanobot{0,0,0,0})
}

func getClosestPoint(nanobot *Nanobot, value int, nanobots []*Nanobot) *Nanobot {
	currNanobot := &Nanobot{nanobot.X, nanobot.Y, nanobot.Z, 0}
	//i know the nanobot has >0 coords

	for {
		currNanobot.X -= 1
		if numIntersectionsCenter(currNanobot, nanobots) < value {
			currNanobot.X += 1
			break
		}
	}

	for {
		currNanobot.Y -= 1
		if numIntersectionsCenter(currNanobot, nanobots) < value {
			currNanobot.Y += 1
			break
		}
	}

	for {
		currNanobot.Z -= 1
		if numIntersectionsCenter(currNanobot, nanobots) < value {
			currNanobot.Z += 1
			break
		}
	}
	return currNanobot
}

func numIntersectionsCenter(nanobot *Nanobot, bots []*Nanobot) int {
	res := 0

	for _, n := range bots {
		if distNanobots(nanobot, n) <= n.Range {
			res++
		}
	}
	return res
}

func getBounds(nanobots []*Nanobot) (int, int, int, int, int, int, int, int) {
	minX := math.MaxInt64
	maxX := math.MinInt64
	minY := minX
	maxY := maxX
	minZ := minX
	maxZ := maxX
	minR := minX
	maxR := 0

	for _, b := range nanobots {
		if b.X < minX {
			minX = b.X
		}
		if b.X > maxX {
			maxX = b.X
		}
		if b.Y < minY {
			minY = b.Y
		}
		if b.Y > maxY {
			maxY = b.Y
		}
		if b.Z < minZ {
			minZ = b.Z
		}
		if b.Z > maxZ {
			maxZ = b.Z
		}
		if b.Range < minR {
			minR = b.Range
		}
		if b.Range > maxR {
			maxR = b.Range
		}
	}
	return minX, maxX, minY, maxY, minZ, maxZ, minR, maxR
}

var nanobotRegex = regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>,\sr=(\d+)`)

func distNanobots(n1 *Nanobot, n2 *Nanobot) int {
	deltaX := n1.X - n2.X
	if deltaX < 0 {
		deltaX = -deltaX
	}
	deltaY := n1.Y - n2.Y
	if deltaY < 0 {
		deltaY = -deltaY
	}
	deltaZ := n1.Z - n2.Z
	if deltaZ < 0 {
		deltaZ = -deltaZ
	}
	return deltaX + deltaY + deltaZ
}

func parseNanobots(line string) *Nanobot {
	toks := nanobotRegex.FindStringSubmatch(line)
	x, _ := strconv.Atoi(toks[1])
	y, _ := strconv.Atoi(toks[2])
	z, _ := strconv.Atoi(toks[3])
	r, _ := strconv.Atoi(toks[4])

	return &Nanobot{x, y, z, r}
}
