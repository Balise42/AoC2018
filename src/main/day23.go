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

	current := &Nanobot{0,0,0,0}
	for k := 0; k<5000; k++{
		currCover, currPen := penalty(current, nanobots)
		penalties := make(map[*Nanobot]int)
		for i := 1; i<=maxR/2 && current.X -i >= minX; i++ {
			nano := &Nanobot {current.X - i, current.Y, current.Z, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen {
				penalties[nano] = newPen
				break
			}
		}

		for i := 1; i<=maxR/2 && current.X + i <= maxX; i++ {
			nano := &Nanobot{current.X + i, current.Y, current.Z, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen  {
				penalties[nano] = newPen
				break
			}
		}

		for i := 1; i<=maxR/2 && current.Y -i >= minY; i++ {
			nano := &Nanobot{current.X, current.Y - i, current.Z, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen {
				penalties[nano] = newPen
				break
			}
		}

		for i := 1; i<=maxR/2 && current.Y + i <= maxY; i++ {
			nano := &Nanobot{current.X, current.Y + i, current.Z, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen {
				penalties[nano] = newPen
				break
			}
		}

		for i := 1; i<=maxR/2 && current.Z -i >= minZ ; i++ {
			nano := &Nanobot{current.X, current.Y, current.Z - i, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen {
				penalties[nano] = newPen
				break
			}
		}

		for i := 1; i<=maxR/2 && current.Z +i <= maxZ; i++ {
			nano := &Nanobot{current.X, current.Y, current.Z + i, 0}
			newCover, newPen := penalty(nano, nanobots)
			if newCover > currCover || newPen < currPen {
				penalties[nano] = newPen
				break
			}
		}

		var nextCandidate *Nanobot = nil
		for k, v := range penalties {
			if v < currPen {
				nextCandidate = k
				currPen = v
			}
		}

		if nextCandidate == nil {
			break
		}
		fmt.Println(currPen)
	}
	fmt.Println(minX, maxX, minY, maxY, minZ, maxZ, minR, maxR)
	fmt.Println((maxX-minX)/minR, (maxY - minY)/minR, (maxZ - minZ)/minR)
	return 0
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

func penalty(nanobot *Nanobot, bots []*Nanobot) (int, int) {
	sum := 0
	numOk := 0
	for _, n := range bots {
		if distNanobots(nanobot, n) > n.Range {
			sum += distNanobots(nanobot, n) - n.Range
		} else {
			numOk++
		}
	}
	return numOk, sum
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
