package main

import (
	"os"
	"bufio"
	"regexp"
	"strconv"
	"fmt"
	"math"
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
func computeMaxBotsPointDistance(nanobots []*Nanobot) int {
	minX, maxX, minY, maxY, minZ, maxZ, minR, maxR := getBounds(nanobots)
	fmt.Println(minX, maxX, minY, maxY, minZ, maxZ, minR, maxR)

	sampledLocaleMax := make(map[*Nanobot]int)
	samplingInterval := minR / 30
	for x := minX; x <= maxX; x += samplingInterval {
		for y := minY; y <= maxZ; y += samplingInterval {
			for z := minZ; z <= maxZ; z += samplingInterval {
				num := numIntersectionsCenter(&Nanobot{x, y, z, 0}, nanobots, 869)
				if num >= 870 {
					sampledLocaleMax[&Nanobot{x, y, z, 0}] = num
				}
			}
		}
	}

	return findMax(sampledLocaleMax, nanobots, samplingInterval)
}

func findMax(locals map[*Nanobot]int, nanobots []*Nanobot, samplingInterval int) int {
	globalMax := 0
	var globalMaxPoint *Nanobot
	for k, v := range locals {
		fmt.Println(k, v)
		localMax, localMaxPoint := findLocalMax(k, v, nanobots, globalMax)
		if localMax > globalMax {
			globalMax = localMax
			globalMaxPoint = localMaxPoint
		} else if localMax == globalMax && closer(localMaxPoint, globalMaxPoint) {
			globalMaxPoint = localMaxPoint
		}
	}
	fmt.Println(globalMax, globalMaxPoint)
	return distNanobots(globalMaxPoint, &Nanobot{0, 0, 0, 0})
}

//94622655 too low
//96277887 too low
//97602486 not right

func findLocalMax(n *Nanobot, currMax int, nanobots []*Nanobot, globalMax int) (int, *Nanobot) {
	localMax := currMax
	localMaxPoint := n

	i := 100

	for x := n.X - i; x <= n.X+i; x++ {
		for y := n.Y - i; y <= n.Y+i; y++ {
			for z := n.Z - i; z <= n.Z+i; z++ {
				candidate := &Nanobot{x, y, z, 0}
				numIntCandidate := numIntersectionsCenter(candidate, nanobots, localMax)
				if numIntCandidate > localMax {
					localMaxPoint = candidate
					localMax = numIntCandidate
				} else if numIntCandidate == localMax && closer(candidate, localMaxPoint) {
					localMaxPoint = candidate
				}
			}
		}
	}

	return localMax, localMaxPoint
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

func closer(n1 *Nanobot, n2 *Nanobot) bool {
	if n2 == nil {
		return true
	}
	return distNanobots(n1, &Nanobot{0, 0, 0, 0}) < distNanobots(n2, &Nanobot{0, 0, 0, 0})
}

func numIntersectionsCenter(nanobot *Nanobot, bots []*Nanobot, globalMax int) int {
	res := 0
	maxNumLeft := 1000

	for _, n := range bots {
		if distNanobots(nanobot, n) <= n.Range {
			res++
		} else {
			maxNumLeft--
		}
		if maxNumLeft+res < globalMax {
			break
		}
	}
	return res
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
