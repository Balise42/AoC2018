package main

import (
	"LineProcessor"
	"os"
	"strconv"
	"fmt"
)

func main() {
	filename := os.Args[1]
	processor := func(line string) int64 {
		num, _ := strconv.ParseInt(line, 10, 64)
		return num
	}
	fmt.Println(LineProcessor.SumLines(filename, processor))
}
