package main

import (
	"os"
	"bufio"
	"regexp"
	"strconv"
	"fmt"
	"math"
	"sort"
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

type OctreeNode struct {
	MinX             int
	MaxX             int
	MinY             int
	MaxY             int
	MinZ             int
	MaxZ             int
	IntersectingBots []*Nanobot
	Parent           *OctreeNode
	Children         []*OctreeNode
}

func (n1 *OctreeNode) lessThan(n2 *OctreeNode) bool {
	if len(n1.IntersectingBots) > len(n2.IntersectingBots) {
		return true
	}
	if len(n2.IntersectingBots) > len(n1.IntersectingBots) {
		return false
	}
	if distNanobots(&Nanobot{n1.MinX, n1.MinY, n1.MinZ, 0}, &Nanobot{0,0,0,0}) < distNanobots(&Nanobot{n2.MinX, n2.MinY, n2.MinZ, 0}, &Nanobot{0,0,0,0}) {
		return true
	}
	if distNanobots(&Nanobot{n1.MinX, n1.MinY, n1.MinZ, 0}, &Nanobot{0,0,0,0}) > distNanobots(&Nanobot{n2.MinX, n2.MinY, n2.MinZ, 0}, &Nanobot{0,0,0,0}) {
		return false
	}
	if n1.MaxX - n1.MinX < n2.MaxX - n2.MinX {
		return true
	}
	return false
}

type byProcessingOrder []*OctreeNode

func (s byProcessingOrder) Len() int {
	return len(s)
}
func (s byProcessingOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byProcessingOrder) Less(i, j int) bool {
	return s[i].lessThan(s[j])
}


func computeMaxBotsPointDistance(nanobots []*Nanobot) int {
	minX, maxX, minY, maxY, minZ, maxZ, _, _ := getBounds(nanobots)
	octree := &OctreeNode{minX, maxX, minY, maxY, minZ, maxZ, nanobots, nil, nil}

	toProcess := byProcessingOrder{octree}

	for len(toProcess) > 0 {
		sort.Sort(byProcessingOrder(toProcess))
		node := toProcess[0]
		toProcess = toProcess[1:]
		if node.MinX == node.MaxX {
			return distNanobots(&Nanobot{node.MinX, node.MinY, node.MinZ, 0}, &Nanobot{0,0,0,0})
		} else {
			splitNode(node)
		}
		toProcess = append(toProcess, node.Children...)
	}
	return 0
}

func splitNode(n *OctreeNode) {

	var midXl, midXr, midYl, midYr, midZl, midZr int
	if n.MinX == n.MaxX || n.MaxX == n.MinX+1 {
		midXl = n.MinX
		midXr = n.MaxX
	} else {
		midXl = (n.MinX + n.MaxX) / 2
		midXr = (n.MinX+n.MaxX)/2 + 1
	}
	if n.MinY == n.MaxY || n.MaxY == n.MinY+1 {
		midYl = n.MinY
		midYr = n.MaxY
	} else {
		midYl = (n.MinY + n.MaxY) / 2
		midYr = (n.MinY+n.MaxY)/2 + 1
	}
	if n.MinZ == n.MaxZ || n.MaxZ == n.MinZ+1 {
		midZl = n.MinZ
		midZr = n.MaxZ
	} else {
		midZl = (n.MinZ + n.MaxZ) / 2
		midZr = (n.MinZ+n.MaxZ)/2 + 1
	}

	n.Children = []*OctreeNode{
		{n.MinX, midXl, n.MinY, midYl, n.MinZ, midZl, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, n.MinY, midYl, n.MinZ, midZl, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, midYr, n.MaxY, n.MinZ, midZl, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, n.MinY, midYl, midZr, n.MaxZ, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, midYr, n.MaxY, n.MinZ, midZl, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, n.MinY, midYl, midZr, n.MaxZ, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, midYr, n.MaxY, midZr, n.MaxZ, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, midYr, n.MaxY, midZr, n.MaxZ, make([]*Nanobot, 0), n, make([]*OctreeNode, 0)},
	}

	for _, c := range n.Children {
		for _, bot := range n.IntersectingBots {
			if intersectsNanoNode(bot, c) {
				c.IntersectingBots = append(c.IntersectingBots, bot)
			}
		}
	}
}

func intersectsNanoNode(nanobot *Nanobot, node *OctreeNode) bool {
	corners := getCorners(node)
	for _, c := range corners {
		if distNanobots(nanobot, c) <= nanobot.Range {
			return true
		}
	}

	if nanobotInNode(&Nanobot{nanobot.X, nanobot.Y, nanobot.Z - nanobot.Range, 0}, node) {
		return true
	}
	if nanobotInNode(&Nanobot{nanobot.X, nanobot.Y, nanobot.Z + nanobot.Range, 0}, node) {
		return true
	}
	if nanobotInNode(&Nanobot{nanobot.X, nanobot.Y - nanobot.Range, nanobot.Z, 0}, node) {
		return true
	}
	if nanobotInNode(&Nanobot{nanobot.X, nanobot.Y + nanobot.Range, nanobot.Z, 0}, node) {
		return true
	}
	if nanobotInNode(&Nanobot{nanobot.X - nanobot.Range, nanobot.Y, nanobot.Z, 0}, node) {
		return true
	}

	if nanobotInNode(&Nanobot{nanobot.X + nanobot.Range, nanobot.Y, nanobot.Z, 0}, node) {
		return true
	}

	return false
}

func nanobotInNode(nanobot *Nanobot, node * OctreeNode) bool {
	return nanobot.X >= node.MinX && nanobot.X <= node.MaxX && nanobot.Y >= node.MinY && nanobot.Y <= node.MaxY && nanobot.Z >= node.MinZ && nanobot.Z <= node.MaxZ
}
func getCorners(n *OctreeNode) []*Nanobot {
	return []*Nanobot{
		{n.MinX, n.MinY, n.MinZ, 0},
		{n.MinX, n.MinY, n.MaxZ, 0},
		{n.MinX, n.MaxY, n.MinZ, 0},
		{n.MaxX, n.MinY, n.MinZ, 0},
		{n.MaxX, n.MinY, n.MaxZ, 0},
		{n.MinX, n.MaxY, n.MaxZ, 0},
		{n.MaxX, n.MaxY, n.MinZ, 0},
		{n.MaxX, n.MaxY, n.MaxZ, 0},
	}
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
		if b.X-b.Range < minX {
			minX = b.X - b.Range
		}
		if b.X+b.Range > maxX {
			maxX = b.X + b.Range
		}
		if b.Y-b.Range < minY {
			minY = b.Y - b.Range
		}
		if b.Y+b.Range > maxY {
			maxY = b.Y + b.Range
		}
		if b.Z-b.Range < minZ {
			minZ = b.Z - b.Range
		}
		if b.Z+b.Range > maxZ {
			maxZ = b.Z + b.Range
		}
		if b.Range < minR {
			minR = b.Range
		}
		if b.Range > maxR {
			maxR = b.Range
		}
	}

	maxDim := maxX - minX
	if maxY - minY > maxDim {
		maxDim = maxY - minY
	}
	if maxZ - maxZ > maxDim {
		maxDim = maxZ - minZ
	}

	maxDimPow2 := 1024
	for maxDimPow2 < maxDim {
		maxDimPow2 *= 2
	}

	maxX = minX + maxDimPow2
	maxY = minY + maxDimPow2
	maxZ = minZ + maxDimPow2

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
