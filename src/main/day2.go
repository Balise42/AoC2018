package main

import (
	"fmt"
	"LineProcessor"
	"os"
)

func main() {
	filename := os.Args[1]
	processor1 := counting(2)
	processor2 := counting(3)
	fmt.Println(LineProcessor.SumLines(filename, processor1) * LineProcessor.SumLines(filename, processor2))
}


func counting(num int) func(string) int64 {
	f := func (line string) int64 {
		count := make(map[rune]int)
		for _, r := range line {
			curr, ok := count[r]
			if ok {
				count[r] = curr + 1
			} else {
				count[r] = 1
			}
		}

		for _, v := range count {
			if v == num {
				return 1
			}
		}
		return 0
	}
	return f
}