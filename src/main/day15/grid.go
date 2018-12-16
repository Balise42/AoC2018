package day15

import (
	"fmt"
	"math"
)

type Grid struct {
	desc            [][]byte
	unit            [][]*Unit
	width           int
	height          int
	Round           int
	strength        int
	initialNumElves int
}

func NewGrid(lines []string, strength int) *Grid {
	gridWidth := len(lines[0])
	gridHeight := len(lines)
	grid := Grid { make([][]byte, gridWidth), make([][]*Unit, gridWidth), gridWidth, gridHeight, 0, strength, 0}
	for i := 0; i<gridHeight; i++ {
		grid.desc[i] = make([]byte, gridHeight)
		grid.unit[i] = make([]*Unit, gridHeight)
	}

	for y, line := range lines {
		for x, char := range line {
			grid.desc[x][y] = byte(char)
			if grid.desc[x][y] == 'G' || grid.desc[x][y] == 'E' {
				grid.unit[x][y] = NewUnit(grid.desc[x][y], x, y)
				if grid.desc[x][y] == 'E' {
					grid.initialNumElves++
				}
			}
		}
	}
	return &grid
}

func (g Grid) Print() {
	fmt.Println("Round ", g.Round)
	for y := 0; y<g.height; y++ {
		for x := 0; x<g.width; x++ {
			fmt.Print(string(g.desc[x][y]))
		}
		fmt.Println("")
	}
	for y := 0; y<g.height; y++ {
		for x := 0; x<g.width; x++ {
			if g.unit[x][y] != nil {
				fmt.Println(x, y, g.unit[x][y].HP)
			}
		}
	}
}

func (g *Grid) PlayRound() bool{
	gameContinues := true
	for y := 0; y<g.height; y++ {
		for x := 0; x<g.width; x++ {
			v := g.desc[x][y]
			if v == 'E' || v == 'G' {
				gameContinues = gameContinues && g.unit[x][y].Activate(g)
			}
		}
	}
	if gameContinues {
		g.cleanUp()
	}

	return gameContinues
}

func (g *Grid) cleanUp() {
	for x, col := range g.desc {
		for y, v := range col {
			if v == 'E' || v == 'G' {
				g.unit[x][y].CleanUp()
			}
		}
	}
	g.Round++
}

func (g *Grid) GetAttackTarget(x int, y int) (int, int) {
	// Assumption: units never enter walls and board is surrounded by walls, so indices are always valid
	curHP := 201
	posX := -1
	posY := -1
	attacker := g.unit[x][y].Side
	if g.isEnemy(x, y-1, attacker) {
		if g.unit[x][y-1].HP < curHP {
			posX = x
			posY = y-1
			curHP = g.unit[x][y-1].HP
		}
	}
	if g.isEnemy(x-1, y, attacker) {
		if g.unit[x-1][y].HP < curHP {
			posX = x-1
			posY = y
			curHP = g.unit[x-1][y].HP
		}
	}
	if g.isEnemy(x+1, y, attacker) {
		if g.unit[x+1][y].HP < curHP {
			posX = x+1
			posY = y
			curHP = g.unit[x+1][y].HP
		}
	}
	if g.isEnemy(x, y+1, attacker) {
		if g.unit[x][y+1].HP < curHP {
			posX = x
			posY = y+1
			curHP = g.unit[x][y+1].HP
		}
	}
	return posX, posY
}

func (g *Grid) Attack(x int, y int) {
	if g.unit[x][y].Side == 'E' {
		g.unit[x][y].HP = g.unit[x][y].HP - 3
	} else {
		g.unit[x][y].HP = g.unit[x][y].HP - g.strength
	}
	if g.unit[x][y].HP <= 0 {
		g.unit[x][y] = nil
		g.desc[x][y] = '.'
	}
}

type Pos struct {
	X int
	Y int
}

func (p Pos) lessThan(other Pos) bool {
	return p.Y < other.Y || (p.Y == other.Y && p.X < other.X)
}

func (p Pos) Manhattan(other Pos) int {
	deltaX := p.X - other.X
	if deltaX < 0 {
		deltaX = -deltaX
	}
	deltaY := p.Y - other.Y
	if deltaY < 0 {
		deltaY = -deltaY
	}
	return deltaX + deltaY
}

func (g *Grid) MoveUnit(x int, y int) bool {
	targets := make([]Pos, 0)
	side := g.unit[x][y].Side
	for x, col := range g.desc {
		for y := range col {
			if g.isEnemy(x, y, side) {
				targets = append(targets, Pos{x, y})
			}
		}
	}

	if len(targets) == 0 {
		return false
	}

	inRange := make([]Pos, 0)
	for _, p := range targets {
		inRange = g.appendIfValid(inRange, p.X-1, p.Y)
		inRange = g.appendIfValid(inRange, p.X, p.Y-1)
		inRange = g.appendIfValid(inRange, p.X, p.Y+1)
		inRange = g.appendIfValid(inRange, p.X+1, p.Y)
	}

	if len(inRange) == 0 {
		//found a target, but can't go there
		return true
	}

	g.MoveToRange(x, y, inRange)
	// we may or may not have been able to move, but we had found a target, so we continue
	return true
}

func (g *Grid) isEnemy(x int, y int, side byte) bool {
	return (g.desc[x][y] == 'G' || g.desc[x][y] == 'E') && g.desc[x][y] != side
}

func (g *Grid) appendIfValid(inRange []Pos, x int, y int) []Pos {
	if g.desc[x][y] == '.' {
		inRange = append(inRange, Pos {x, y})
	}
	return inRange
}

func (g * Grid) distPath(x int, y int, pos Pos) int {
	if g.desc[x][y] != '.' {
		return math.MaxInt64
	}

	if pos.X == x && pos.Y == y {
		return 0
	}

	reachable := make(map[Pos]bool)
	reachable[Pos{x,y}] = true
	hit := make(map[Pos]bool)
	hit[Pos{x,y}] = true

	dist := 0
	for {
		reachable = g.growReach(reachable, hit)
		dist++
		if len(reachable) == 0 {
			//couldn't reach anything
			break
		}
		for k := range reachable {
			if k == pos {
				return dist
			}
		}
	}
	return math.MaxInt64
}

func (g *Grid) MoveToRange(x int, y int, pos []Pos) {
	reachable := make(map[Pos]bool)
	reachable[Pos{x,y}] = true
	hit := make(map[Pos]bool)
	hit[Pos{x,y}] = true

	dist := 0
	closest := Pos{math.MaxInt64, math.MaxInt64 }
	foundIt := false
	for {
		reachable = g.growReach(reachable, hit)
		dist++
		if len(reachable) == 0 {
			//couldn't reach anything
			break
		}
		for k := range reachable {
			if contains(pos, k) && k.lessThan(closest) {
				closest = k
				foundIt = true
			}
		}
		if foundIt {
			g.MoveToPos(x, y, closest)
			break
		}
	}
}

func (g *Grid) MoveToPos(x int, y int, dest Pos) {

	dist := math.MaxInt64
	var curr Pos

	tmpDist := g.distPath(x, y-1, dest)
	if tmpDist < dist {
		curr.X = x
		curr.Y = y-1
		dist = tmpDist
	}

	tmpDist = g.distPath(x-1, y, dest)
	if tmpDist < dist {
		curr.X = x-1
		curr.Y = y
		dist = tmpDist
	}

	tmpDist = g.distPath(x+1, y, dest)
	if tmpDist < dist {
		curr.X = x+1
		curr.Y = y
		dist = tmpDist
	}

	tmpDist = g.distPath(x, y+1, dest)
	if tmpDist < dist {
		curr.X = x
		curr.Y = y+1
		dist = tmpDist
	}

	if dist <math.MaxInt64 {
		u := g.unit[x][y]
		g.desc[curr.X][curr.Y] = u.Side
		g.desc[x][y] = '.'
		u.PosX = curr.X
		u.PosY = curr.Y
		g.unit[curr.X][curr.Y] = u
		g.unit[x][y] = nil
	} else {
		fmt.Println("coin")
	}
}

func (g Grid) growPaths(oldPaths [][]Pos, dist int, dest Pos) [][]Pos {
	paths := make([][]Pos, len(oldPaths))
	for i := range oldPaths {
		paths[i] = make([]Pos, len(oldPaths[i]))
		copy(paths[i], oldPaths[i])
	}
	for i := 0; i<dist; i++ {
		newPaths := make([][]Pos, 0)
		for _, path := range paths {
			reach := g.getReach(map[Pos]bool{path[len(path)-1]: true})
			for _, pos := range reach {
				//looping will for sure not yield the shortest path
				if !contains(path, pos) && pos.Manhattan(dest) <= dist - i {
					newPath := make([]Pos, len(path), len(path)+1)
					copy(newPath, path)
					newPaths = append(newPaths, append(newPath, pos))
				}
			}
		}
		paths = newPaths
	}
	return paths
}

func (g Grid) getReach(reachable map[Pos]bool) []Pos {
	newReach := make([]Pos, 0)
	for k := range reachable {
		if g.desc[k.X][k.Y-1] == '.' {
			newReach = append(newReach, Pos{k.X, k.Y - 1})
		}
		if g.desc[k.X-1][k.Y] == '.' {
			newReach = append(newReach, Pos{k.X-1, k.Y})
		}
		if g.desc[k.X+1][k.Y] == '.' {
			newReach = append(newReach, Pos{k.X+1, k.Y})
		}
		if g.desc[k.X][k.Y+1] == '.' {
			newReach = append(newReach, Pos{k.X, k.Y + 1})
		}
	}
	return newReach
}



func (g Grid) growReach(reachable map[Pos]bool, hit map[Pos]bool) map[Pos]bool {
	newReach := make(map[Pos]bool)
	for k := range reachable {
		if g.desc[k.X][k.Y-1] == '.' && !hit[Pos{k.X, k.Y-1}] {
			newReach[Pos{k.X, k.Y-1}] = true
			hit[Pos{k.X, k.Y-1}] = true
		}
		if g.desc[k.X-1][k.Y] == '.' && !hit[Pos{k.X-1, k.Y}]{
			newReach[Pos{k.X-1, k.Y}] = true
			hit[Pos{k.X-1, k.Y}] = true
		}
		if g.desc[k.X+1][k.Y] == '.' && !hit[Pos{k.X+1, k.Y}] {
			newReach[Pos{k.X+1, k.Y}] = true
			hit[Pos{k.X+1, k.Y}] = true
		}
		if g.desc[k.X][k.Y+1] == '.' && !hit[Pos{k.X, k.Y+1}] {
			newReach[Pos{k.X, k.Y+1}] = true
			hit[Pos{k.X, k.Y+1}] = true
		}
	}
	return newReach
}

func (g Grid) Outcome() (int, bool) {
	hp := 0
	numElves := 0
	for y := 0; y<g.height; y++ {
		for x := 0; x < g.width; x++ {
			v := g.desc[x][y]
			if v == 'E' || v == 'G' {
				hp = hp + g.unit[x][y].HP
				if v == 'E' {
					numElves++
				}
			}
		}
	}
	return hp * g.Round, numElves == g.initialNumElves
}

func contains(list []Pos, pos Pos) bool {
	for _, p := range list {
		if p == pos {
			return true
		}
	}
	return false
}



