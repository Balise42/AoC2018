package main

import (
	"os"
	"bufio"
	"fmt"
	"day24/main/game"
)

func main() {
	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	g := game.New(lines, 0)

	for !g.IsFinished() {
		g.PlayRound()
	}

	fmt.Println("Part A: ", g.WinnerScore())

	boost := 1
	for {
		g = game.New(lines, boost)
		for {
			if g.IsFinished() || g.Stalemate {
				break
			}
			g.PlayRound()
		}
		if g.ImmuneWins() {
			break
		}
		boost++
	}
	fmt.Println("Part B: ", g.WinnerScore())
}
