package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type elems []int64

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

type instr struct {
	instruction string
	args        elems
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

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()
	toks := strings.Split(firstLine, " ")

	instrptr, _ := strconv.Atoi(toks[1])
	registers := elems{0, 0, 0, 0, 0, 0}

	// okay it's just the sum of all integer factors.

	program := make([]*instr, 0)
	for scanner.Scan() {
		program = appendInstr(program, scanner.Text())
	}


	nums := make([]int64, 0)
	seen := make(map[int64]int)
	k := 0
	it := 0
	registers = elems{int64(k), 0, 0, 0, 0, 0}
	for {
		it++
		if (registers[instrptr] == 28) {
			myit, ok := seen[registers[4]]
			if ok {
				fmt.Println(myit, registers[4])
				break
			} else {
				nums = append(nums, registers[4])
				seen[registers[4]] = it
			}

		}


		execute(program[registers[instrptr]], registers, functions)
		registers[instrptr]++
	}
	fmt.Println(nums[len(nums)-1])

}

func printState(registers elems, instruction *instr) {
	for i := 0; i < 6; i++ {
		fmt.Printf("%15d ", registers[i])
	}
	fmt.Println(instruction)
}

func execute(instruction *instr, regs elems, functions map[string]func(elems, elems)) {
	functions[instruction.instruction](instruction.args, regs)
}

func appendInstr(program []*instr, line string) []*instr {
	instruction := new(instr)
	toks := strings.Split(line, " ")
	instruction.instruction = toks[0]
	arg1, _ := strconv.Atoi(toks[1])
	arg2, _ := strconv.Atoi(toks[2])
	arg3, _ := strconv.Atoi(toks[3])
	instruction.args = elems{int64(arg1), int64(arg2), int64(arg3)}
	return append(program, instruction)
}
