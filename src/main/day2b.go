package main

import (
	"os"
	"fmt"
	"bufio"
)

func main() {
	filename := os.Args[1]
	var ids []string

	f, _ := os.Open(filename)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	for i, id1 := range ids {
		for j, id2 := range ids {
			if i != j {
				res, correct := getMatching(id1, id2)
				if correct {
					fmt.Print(res)
					return
				}
			}
		}
	}
}

func getMatching(id1 string, id2 string) (string, bool) {
	currStr := ""
	currNum := 0

	for i := range id1 {
		if id1[i] != id2[i] {
			currNum++
			if currNum > 1 {
				break
			}
		} else {
			currStr = currStr + string(id1[i])
		}
	}
	return currStr, currNum == 1
}
