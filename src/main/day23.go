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

type OctreeNode struct {
	MinX     int
	MaxX     int
	MinY     int
	MaxY     int
	MinZ     int
	MaxZ     int
	NumBots  int
	Parent   *OctreeNode
	Children []*OctreeNode
}

func computeMaxBotsPointDistance(nanobots []*Nanobot) int {
	minX, maxX, minY, maxY, minZ, maxZ, _, _ := getBounds(nanobots)
	octree := &OctreeNode{minX, maxX, minY, maxY, minZ, maxZ, 0, nil, nil}

	for _, nanobot := range nanobots {
		fmt.Println(nanobot)
		insertNanobotInOctree(octree, nanobot)
	}

	return getMaxNumberBots(octree)
}
func getMaxNumberBots(node *OctreeNode) int {
	numBots := node.NumBots

	maxChildren := 0
	for _, child := range node.Children {
		tmpMax := getMaxNumberBots(child)
		if tmpMax > maxChildren {
			maxChildren = tmpMax
		}
	}
	return numBots + maxChildren
}

func insertNanobotInOctree(node *OctreeNode, nanobot *Nanobot) {
	if isEntirelyIn(node, nanobot) {
		node.NumBots++
	} else if isEntirelyOut(node, nanobot) {
		// do nothing actually
	} else {
		splitNodeIfNecessary(node)
		for _, child := range node.Children {
			insertNanobotInOctree(child, nanobot)
		}
	}
}

func splitNodeIfNecessary(n *OctreeNode) {
	if len(n.Children) != 0 {
		return
	}

	var midXl, midXr, midYl, midYr, midZl, midZr int
	if n.MinX == n.MaxX || n.MaxX == n.MinX + 1 {
		midXl = n.MinX
		midXr = n.MaxX
	} else {
		midXl = (n.MinX + n.MaxX)/2
		midXr = (n.MinX + n.MaxX)/2 + 1
	}
	if n.MinY == n.MaxY || n.MaxY == n.MinY + 1{
		midYl = n.MinY
		midYr = n.MaxY
	} else {
		midYl = (n.MinY + n.MaxY)/2
		midYr = (n.MinY + n.MaxY)/2 + 1
	}
	if n.MinZ == n.MaxZ || n.MaxZ == n.MinZ + 1{
		midZl = n.MinZ
		midZr = n.MaxZ
	} else {
		midZl = (n.MinZ + n.MaxZ)/2
		midZr = (n.MinZ + n.MaxZ)/2 + 1
	}

	n.Children = []*OctreeNode{
		{n.MinX, midXl, n.MinY, midYl, n.MinZ, midZl, 0, n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, n.MinY, midYl, n.MinZ, midZl, 0, n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, midYr, n.MaxY, n.MinZ, midZl, 0, n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, n.MinY, midYl, midZr, n.MaxZ, 0, n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, midYr, n.MaxY, n.MinZ, midZl, 0, n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, n.MinY, midYl, midZr, n.MaxZ, 0, n, make([]*OctreeNode, 0)},
		{n.MinX, midXl, midYr, n.MaxY, midZr, n.MaxZ, 0, n, make([]*OctreeNode, 0)},
		{midXr, n.MaxX, midYr, n.MaxY, midZr, n.MaxZ, 0, n, make([]*OctreeNode, 0)},
	}
}

func isEntirelyIn(node *OctreeNode, nanobot *Nanobot) bool {
	corners := getCorners(node)

	for _, c := range corners {
		if distNanobots(nanobot, c) > nanobot.Range {
			return false
		}
	}
	return true
}

func isEntirelyOut(node *OctreeNode, nanobot *Nanobot) bool {
	corners := getCorners(node)
	for _, c := range corners {
		if distNanobots(nanobot, c) <= nanobot.Range {
			return false
		}
	}

	if isCenterNodeInNanobot(node, nanobot) {
		return false
	}

	if isBoundingBoxOut(node, nanobot) {
		return true
	}

	/*for x := node.MinX; x<= node.MaxX; x++ {
		for y:= node.MinY; y<= node.MaxY; y++ {
			if distNanobots(nanobot, &Nanobot{x, y, node.MinZ, 0}) <= nanobot.Range {
				return false
			}
			if distNanobots(nanobot, &Nanobot{x, y, node.MaxZ, 0}) <= nanobot.Range {
				return false
			}
		}

		for z:= node.MinZ; z<= node.MaxZ; z++ {
			if distNanobots(nanobot, &Nanobot{x, node.MinY, z, 0}) <= nanobot.Range {
				return false
			}
			if distNanobots(nanobot, &Nanobot{x, node.MaxY, z, 0}) <= nanobot.Range {
				return false
			}
		}
	}

	for y:= node.MinY; y<= node.MaxY; y++ {
		for z:= node.MinZ; z<= node.MaxZ; z++ {
			if distNanobots(nanobot, &Nanobot{node.MinX, y, z, 0}) <= nanobot.Range {
				return false
			}
			if distNanobots(nanobot, &Nanobot{node.MaxX, y, z, 0}) <= nanobot.Range {
				return false
			}
		}
	}*/

	return true
}
func isBoundingBoxOut(node *OctreeNode, nanobot *Nanobot) bool {
	return node.MaxX < nanobot.X - nanobot.Range || node.MinX > nanobot.X + nanobot.Range || node.MaxY < nanobot.Y - nanobot.Range  || node.MinY > nanobot.Y + nanobot.Range || node.MaxZ < nanobot.Z - nanobot.Range || node.MinZ > nanobot.Z + nanobot.Range
}
func isCenterNodeInNanobot(node *OctreeNode, nanobot *Nanobot) bool {
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
			minX = b.X-b.Range
		}
		if b.X+b.Range > maxX {
			maxX = b.X + b.Range
		}
		if b.Y - b.Range < minY {
			minY = b.Y - b.Range
		}
		if b.Y + b.Range > maxY {
			maxY = b.Y + b.Range
		}
		if b.Z - b.Range < minZ {
			minZ = b.Z - b.Range
		}
		if b.Z + b.Range > maxZ {
			maxZ = b.Z + b.Range
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
