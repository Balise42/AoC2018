package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type elems []int

func addr(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] + reg[instr[1]]
}

func addi(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] + instr[1]
}

func mulr(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] * reg[instr[1]]
}

func muli(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] * instr[1]
}

func borr(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] | reg[instr[1]]
}

func bori(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] | instr[1]
}

func banr(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] & reg[instr[1]]
}

func bani(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]] & instr[1]
}

func setr(instr elems, reg elems) {
	reg[instr[2]] = reg[instr[0]]
}

func seti(instr elems, reg elems) {
	reg[instr[2]] = instr[0]
}

func gtir(instr elems, reg elems) {
	if instr[0] > reg[instr[1]] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func gtri(instr elems, reg elems) {
	if reg[instr[0]] > instr[1] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func gtrr(instr elems, reg elems) {
	if reg[instr[0]] > reg[instr[1]] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func eqir(instr elems, reg elems) {
	if instr[0] == reg[instr[1]] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func eqri(instr elems, reg elems) {
	if reg[instr[0]] == instr[1] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func eqrr(instr elems, reg elems) {
	if reg[instr[0]] == reg[instr[1]] {
		reg[instr[2]] = 1
	} else {
		reg[instr[2]] = 0
	}
}

func main() {
	functions := map[string]func(elems, elems){
		"addr": addr,
		"addi": addi,
		"mulr": mulr,
		"muli": muli,
		"banr": banr,
		"bani": bani,
		"borr": borr,
		"bori": bori,
		"setr": setr,
		"seti": seti,
		"gtir": gtir,
		"gtri": gtri,
		"gtrr": gtrr,
		"eqir": eqir,
		"eqri": eqri,
		"eqrr": eqrr,
	}

	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	resA := 0

	candidates := make([][]string, 16)
	for i := 0; i<16; i++ {
		candidates[i] = []string {
			"addr",
			"addi",
			"mulr",
			"muli",
			"banr",
			"bani",
			"borr",
			"bori",
			"setr",
			"seti",
			"gtir",
			"gtri",
			"gtrr",
			"eqir",
			"eqri",
			"eqrr",
		}
	}
	ioRegex := regexp.MustCompile(`(Before|After):\s+\[(\d+),\s(\d+),\s(\d+),\s(\d+)]`)

	startProg := 0

	for i := 0; i < len(lines); i += 3 {
		if !strings.HasPrefix(lines[i], "Before") {
			startProg = i
			break
		}

		inputStr := ioRegex.FindStringSubmatch(lines[i])
		outputStr := ioRegex.FindStringSubmatch(lines[i+2])
		input := convertIoStrToElems(inputStr)
		output := convertIoStrToElems(outputStr)
		instr := convertStrToElems(strings.Split(lines[i+1], " "))
		currCandidates := getCandidates(input, output, instr[1:], functions)

		if len(currCandidates) >= 3 {
			resA++
		}
		candidates[instr[0]] = intersect(candidates[instr[0]], currCandidates)
	}

	fmt.Println("Part A: ", resA)

	candidates = resolveCandidates(candidates)

	registers := elems{0,0,0,0}
	executeProgram(lines[startProg:], functions, candidates, registers)
	fmt.Println("Part B: ", registers[0])
}

func executeProgram(codeLines []string, functions map[string]func(elems, elems), funcMapping [][]string, registers elems) {
	for _, line := range codeLines {
		instr := convertStrToElems(strings.Split(line, " "))
		function := functions[funcMapping[instr[0]][0]]
		function(instr[1:], registers)
	}
}

func resolveCandidates(candidates [][]string) [][]string {
	var removed bool
	changed := true
	for changed {
		changed = false
		for i := 0; i<16; i++ {
			if len(candidates[i]) == 1 {
				for j := 0; j<16; j++ {
					if i != j {
						candidates[j], removed = removeElem(candidates[j], candidates[i][0])
						if removed {
							changed = true
						}
					}
				}
			}
		}
	}
	return candidates
}

func removeElem(str []string, s string) ([]string, bool) {
	i := index(str, s)
	if i != -1 {
		return append(str[:i], str[i+1:]...), true
	} else {
		return str, false
	}
}


func intersect(candidates []string, currCandidates []string) []string {
	res := make([]string, 0)
	for _, c := range candidates {
		if containsFunc(currCandidates, c) {
			res = append(res, c)
		}
	}
	return res
}

func index(functions []string, f string) int {
	for i, function := range functions {
		if function == f {
			return i
		}
	}
	return -1
}

func containsFunc(functions []string, f string) bool {
	return index(functions, f) != -1
}

func convertIoStrToElems(toks []string) elems {
	elems := make(elems, 4)
	for i := 2; i < 6; i++ {
		elems[i-2], _ = strconv.Atoi(toks[i])
	}
	return elems
}

func convertStrToElems(toks []string) elems {
	elems := make(elems, 4)
	for i := 0; i < 4; i++ {
		elems[i], _ = strconv.Atoi(toks[i])
	}
	return elems
}

func getCandidates(input elems, output elems, instr elems, functions map[string]func(elems, elems)) []string {
	res := make([]string, 0)
	for name, f := range functions {
		inputCopy := make(elems, 4)
		copy(inputCopy, input)
		f(instr, inputCopy)
		if inputCopy[instr[2]] == output[instr[2]] {
			res = append(res, name)
		}
	}
	return res
}
