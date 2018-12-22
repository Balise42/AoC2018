package main

import (
	"fmt"
	"container/heap"
	"math"
)

func main() {
	targetX := 10
	targetY := 785
	depth := 5616

	sizeX := targetX * 5
	sizeY := targetY * 5

	/*targetY = 10
	depth = 510
	sizeX = 16
	sizeY = 16*/

	risk := 0
	erosion := make([][]int, sizeX)
	caveMap := make([][]int, sizeX)

	for i := 0; i < sizeX; i++ {
		erosion[i] = make([]int, sizeY)
		caveMap[i] = make([]int, sizeY)
	}

	erosion[0][0] = depth % 20183
	risk = risk + erosion[0][0]%3
	caveMap[0][0] = erosion[0][0] % 3

	for x := 1; x < sizeX; x++ {
		erosion[x][0] = (x*16807 + depth) % 20183
		if x <= targetX {
			risk = risk + erosion[x][0]%3
		}
		caveMap[x][0] = erosion[x][0] % 3
	}
	for y := 1; y < sizeY; y++ {
		erosion[0][y] = (y*48271 + depth) % 20183
		if y <= targetY {
			risk = risk + erosion[0][y]%3
		}
		caveMap[0][y] = erosion[0][y] % 3
	}
	for x := 1; x < sizeX; x++ {
		for y := 1; y < sizeY; y++ {
			if x == targetX && y == targetY {
				erosion[x][y] = depth % 20183
			} else {
				erosion[x][y] = (erosion[x-1][y]*erosion[x][y-1] + depth) % 20183
			}
			if x <= targetX && y <= targetY {
				risk = risk + erosion[x][y]%3
			}
			caveMap[x][y] = erosion[x][y] % 3
		}
	}

	fmt.Println(risk)
/*	for y := 0; y < sizeY; y++ {
		for x := 0; x < sizeX; x++ {
			fmt.Print(getChar(caveMap[x][y]))
		}
		fmt.Println("")
	}*/
	fmt.Println(getShortestPath(caveMap, targetX, targetY, sizeX, sizeY))
}

func getChar(i int) string {
	if i == 0 {
		return "."
	}
	if i == 1 {
		return "="
	}
	return "|"

}

type caveNode struct {
	X         int
	Y         int
	gear      int // torch = 0; climbing = 1; neither = 2
	neighbors map[*caveNode]int
}


//prio queue from golang doc

type caveNodeQElem struct {
	value    *caveNode
	priority int
	index    int
}

type cavePriorityQueue []*caveNodeQElem

func (pq cavePriorityQueue) Len() int { return len(pq) }


func (pq cavePriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq cavePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *cavePriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*caveNodeQElem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *cavePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *cavePriorityQueue) update(X int, Y int, G int, priority int) {
	for _, v := range *pq {
		if v.value.X == X && v.value.Y == Y && v.value.gear == G {
			v.priority = priority
			heap.Fix(pq, v.index)
		}
	}
}

func getShortestPath(caveMap [][]int, targetX int, targetY int, sizeX int, sizeY int) int {

	nodes := make([][][]*caveNode, sizeX)
	dists := make([][][]int, sizeX)
	for x := 0; x < sizeX; x++ {
		nodes[x] = make([][]*caveNode, sizeY)
		dists[x] = make([][]int, sizeY)
		for y := 0; y < sizeY; y++ {
			dists[x][y] = make([]int, 3)
			nodes[x][y] = make([]*caveNode, 3)
		}
	}

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			nodes[x][y][caveMap[x][y]] = &caveNode{x, y, caveMap[x][y], make(map[*caveNode]int)}
			nodes[x][y][(caveMap[x][y]+1)%3] = &caveNode{x, y, (caveMap[x][y] + 1) % 3, make(map[*caveNode]int)}
			nodes[x][y][caveMap[x][y]].neighbors[nodes[x][y][(caveMap[x][y]+1)%3]] = 7
			nodes[x][y][(caveMap[x][y]+1)%3].neighbors[nodes[x][y][caveMap[x][y]]] = 7
			for g := 0; g<3; g++ {
				dists[x][y][g] = math.MaxInt64
			}
		}
	}

	dists[0][0][0] = 0
	pq := make(cavePriorityQueue, 0)
	heap.Init(&pq)
	start := &caveNodeQElem{value: nodes[0][0][0], priority: 0}
	heap.Push(&pq, start)

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			for g := 0; g < 3; g++ {
				if nodes[x][y][g] != nil {
					if x > 0 && nodes[x-1][y][g] != nil {
						nodes[x][y][g].neighbors[nodes[x-1][y][g]] = 1
					}
					if y > 0 && nodes[x][y-1][g] != nil {
						nodes[x][y][g].neighbors[nodes[x][y-1][g]] = 1
					}
					if x < sizeX-1 && nodes[x+1][y][g] != nil{
						nodes[x][y][g].neighbors[nodes[x+1][y][g]] = 1
					}
					if y < sizeY-1 && nodes[x][y+1][g] != nil{
						nodes[x][y][g].neighbors[nodes[x][y+1][g]] = 1
					}

					if x != 0 || y != 0 || g != 0 {
						heap.Push(&pq, &caveNodeQElem{value: nodes[x][y][g], priority:math.MaxInt64})
					}
				}
			}
		}
	}

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*caveNodeQElem)
		if u.value.X == targetX && u.value.Y == targetY && u.value.gear == 0 {
			break
		}
		for n, d := range u.value.neighbors {
			altD := dists[u.value.X][u.value.Y][u.value.gear] + d
			if altD < dists[n.X][n.Y][n.gear] {
				dists[n.X][n.Y][n.gear] = altD
				pq.update(n.X, n.Y, n.gear, altD)
			}
		}
	}

	return dists[targetX][targetY][0]
}
