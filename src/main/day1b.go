package main

import (
	"os"
	"strconv"
	"fmt"
	"bufio"
)

func main() {
	filename := os.Args[1]
	file, _ := os.Open(filename)
	defer file.Close()

	seen := make(map[int64]bool)

	var sum int64
	sum = 0
	seen[0] = true

	finish := false

	var list []int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		list = append(list, num)
	}

	var i int
	i = 0
	for !finish {
		sum = sum + list[i%len(list)]
		_, in := seen[sum]
		if in {
			fmt.Println(sum)
			finish = true
		}
		seen[sum] = true
		i++
	}
}
