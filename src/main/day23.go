package main

import (
	"os"
	"bufio"
	"regexp"
	"strconv"
	"fmt"
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
	intervalCandidates := make([]*Nanobot, 0)
	maxIntersect := 0

	for _, n := range nanobots {
		intersect := numIntersectionsWithBounds(n, nanobots)
		if intersect == maxIntersect {
			intervalCandidates = append(intervalCandidates, n)
		} else if intersect > maxIntersect {
			maxIntersect = intersect
			intervalCandidates = []*Nanobot{n}
		}
	}

	var candidate *Nanobot = nil

	// i actually know i have a single candidate here; it would require more care in the general case
	for _, i := range intervalCandidates {
		for _, b := range getNanobotBounds(i) {
			if numIntersectionsCenter(b, nanobots) == maxIntersect && closer(b, candidate) {
				candidate = b
			}
		}
	}
	fmt.Println(candidate)

	delta := 0

	if candidate.X < 0 {
		delta = 1
	} else if candidate.X > 0 {
		delta = -1
	}
	if delta != 0 {
		for {
			newCandidate := &Nanobot{candidate.X + delta, candidate.Y, candidate.Z, candidate.Range}
			if closer(newCandidate, candidate) && numIntersectionsCenter(newCandidate, nanobots) == maxIntersect {
				candidate = newCandidate
			} else {
				break
			}
		}
	}

	delta = 0
	if candidate.Y < 0 {
		delta = 1
	} else if candidate.Y > 0 {
		delta = -1
	}
	if delta != 0 {
		for {
			newCandidate := &Nanobot{candidate.X, candidate.Y + delta, candidate.Z, candidate.Range}
			if closer(newCandidate, candidate) && numIntersectionsCenter(newCandidate, nanobots) == maxIntersect {
				candidate = newCandidate
			} else {
				break
			}
		}
	}

	delta = 0
	if candidate.Z < 0 {
		delta = 1
	} else if candidate.Z > 0 {
		delta = -1
	}
	if delta != 0 {
		for {
			newCandidate := &Nanobot{candidate.X, candidate.Y, candidate.Z + delta, candidate.Range}
			if closer(newCandidate, candidate) && numIntersectionsCenter(newCandidate, nanobots) == maxIntersect {
				candidate = newCandidate
			} else {
				break
			}
		}
	}


	// 97771477 too low
	return distNanobots(candidate, &Nanobot{0,0,0,0})
}

func closer(n1 *Nanobot, n2 *Nanobot) bool {
	if n2 == nil {
		return true
	}
	return distNanobots(n1, &Nanobot{0,0,0,0}) < distNanobots(n2, &Nanobot{0,0,0,0})
}

func numIntersectionsWithBounds(nanobot *Nanobot, bots []*Nanobot) int {
	bounds := getNanobotBounds(nanobot)
	res := 0

	for _, b := range bounds {
		num := numIntersectionsCenter(b, bots)
		if num > res {
			res = num
		}
	}

	return res
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

func getNanobotBounds(n *Nanobot) []*Nanobot {
	return []*Nanobot{
		{n.X - n.Range, n.Y - n.Range, n.Z - n.Range, 0},
		{n.X + n.Range, n.Y - n.Range, n.Z - n.Range, 0},
		{n.X - n.Range, n.Y + n.Range, n.Z - n.Range, 0},
		{n.X - n.Range, n.Y - n.Range, n.Z + n.Range, 0},
		{n.X + n.Range, n.Y + n.Range, n.Z - n.Range, 0},
		{n.X - n.Range, n.Y + n.Range, n.Z + n.Range, 0},
		{n.X + n.Range, n.Y - n.Range, n.Z + n.Range, 0},
		{n.X + n.Range, n.Y + n.Range, n.Z + n.Range, 0},
	}
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
