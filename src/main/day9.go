package main

import "fmt"

type Game struct {
	marbles []int
	currPos int
}

func newGame() Game {
	return Game{[]int{0}, 0}
}

func (g *Game) PlayMarble(number int) int {
	if number % 23 != 0 {
		g.addMarble(number)
		return 0
	} else {
		score := number
		score = score + g.removeMarble()
		return score
	}
}

func (g *Game) removeMarble() int {
	pos := g.getRemovalPos()
	value := g.marbles[pos]
	g.marbles = append(g.marbles[0:pos], g.marbles[pos+1:]...)
	g.currPos = pos
	return value
}

func (g *Game) addMarble(number int) {
	pos := g.getNextPos()
	if pos == len(g.marbles) {
		g.marbles = append(g.marbles, number)
	} else if pos == 0 {
		g.marbles = append([]int{number}, g.marbles...)
	} else {
		g.marbles = append(g.marbles[:pos], append([]int{number}, g.marbles[pos:]...)...)
	}
	g.currPos = pos
}

func (g Game) getNextPos() int {
	pos := g.currPos + 2
	return g.getRealPos(pos)
}

func(g Game) getRealPos(number int) int {
	if len(g.marbles) == 1 {
		return 1
	}
	if number > 0 && number <= len(g.marbles) {
		return number
	}
	pos := number
	pos = number % (len(g.marbles))
	if pos < 0 {
		pos = pos + len(g.marbles)
	}
	return pos
}

func (g Game) getRemovalPos() int {
	pos := g.currPos - 7
	return g.getRealPos(pos)
}

func main() {
	numPlayers := 428
	numMarbles := 7206100

	scores := make([]int, numPlayers)
	g := newGame()
	player := 0
	maxScore := 0

	for i := 1; i<=numMarbles; i++ {
		scores[player] += g.PlayMarble(i)
		if scores[player] > maxScore {
			maxScore = scores[player]
		}
		player = (player + 1)%numPlayers
	}

	fmt.Print(maxScore)
}
