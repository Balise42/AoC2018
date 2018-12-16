package main

import (
	"bufio"
	"fmt"
	"main/day15"
	"os"
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

	grid := day15.NewGrid(lines, 3)
	for {
		cont := grid.PlayRound()
		if !cont {
			break
		}
	}

	o, _ := grid.Outcome()
	fmt.Println("Part A: ", o)

	for i := 3; i < 200; i++ {

		grid := day15.NewGrid(lines, i)
		for {
			cont := grid.PlayRound()
			if !cont {
				break
			}
		}
		o, valid := grid.Outcome()
		if valid {
			fmt.Println("Part B: ", o)
			break
		}
	}
}
