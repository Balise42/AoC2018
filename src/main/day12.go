package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Automaton struct {
	State []rune
	Start int64
	Rules []string
}

func NewAutomaton(initial string, rules1 []string) Automaton {
	res := Automaton{}
	res.State = []rune("....." + initial + ".....")
	res.Start = -5
	res.Rules = rules1
	return res
}

func (a *Automaton) evolve() {
	newState := make([]rune, len(a.State))
	for i := 2; i < len(newState)-2; i++ {
		newState[i] = a.getState(a.State[i-2 : i+3])
	}
	newState[0] = '.'
	newState[1] = '.'
	newState[len(newState)-2] = '.'
	newState[len(newState)-1] = '.'
	if string(newState[:5]) != "....." {
		newState = append([]rune("....."), newState...)
		a.Start = a.Start - 5
	}
	if string(newState[len(newState)-5:]) != "....." {
		newState = append(newState, []rune(".....")...)
	}
	for string(newState[:10]) == ".........." {
		newState = newState[5:]
		a.Start = a.Start + 5
	}
	a.State = newState
}

func (a Automaton) getState(s []rune) rune {
	for _, rule := range a.Rules {
		if rule == string(s) {
			return '#'
		}
	}
	return '.'
}

func (a Automaton) sumPots() int64 {
	sum := int64(0)
	for i, c := range a.State {
		if c == '#' {
			sum = sum + int64(i) + a.Start
		}
	}
	return sum
}

func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	rules := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line[len(line)-1] == '#' {
			rules = append(rules, strings.Split(line, " ")[0])
		}
	}

	a := NewAutomaton("##.##.#.#...#......#..#.###..##...##.#####..#..###.########.##.....#...#...##....##.#...#.###...#.##", rules)

	seen := make(map[string]int64)
	states := make([]Automaton, 0)
	states = append(states, a)

	it := int64(0)
	for {
		it++
		a.evolve()
		if seen[string(a.State)] != 0 {
			break
		}
		seen[string(a.State)] = it
		states = append(states, a)
	}

	fmt.Println(a)
	numGens := int64(50000000000)
	//numGens := int64(365)

	startRepeat := seen[string(a.State)]
	period := it - seen[string(a.State)]
	delta := a.Start - states[startRepeat].Start

	iter := (numGens - startRepeat) % period
	fmt.Println(iter)

	finalAutomaton := NewAutomaton(string(states[startRepeat+iter].State), rules)
	finalAutomaton.Start = states[startRepeat+iter].Start + delta*((numGens-startRepeat)/period)
	if string(finalAutomaton.State[:10]) == ".........." {
		finalAutomaton.State = finalAutomaton.State[5:]
	}
	if string(finalAutomaton.State[len(finalAutomaton.State) - 10:]) == ".........." {
		finalAutomaton.State = finalAutomaton.State[:len(finalAutomaton.State) - 5]
	}
	fmt.Println(finalAutomaton.sumPots())

	/*testAutomaton := NewAutomaton("##.##.#.#...#......#..#.###..##...##.#####..#..###.########.##.....#...#...##....##.#...#.###...#.##", rules)
	for i := int64(0); i < numGens; i++ {
		testAutomaton.evolve()
	}
	fmt.Println(testAutomaton)*/
}
