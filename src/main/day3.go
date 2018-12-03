package main

import (
	"os"
	"bufio"
	"fmt"
	"common/day3"
)

func main() {
	filename := os.Args[1]
	var claims []day3.Claim

	f, _ := os.Open(filename)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		claims = append(claims, day3.CreateClaim(scanner.Text()))
	}

	sum := 0

	fabric := make(map[day3.Pos]int)

	for _, c := range claims {
		for i := 0; i<c.Width; i++ {
			for j := 0; j<c.Height; j++ {
				used, ok := fabric[day3.Pos{c.X+i, c.Y+j}]
				if ok {
					fabric[day3.Pos{c.X+i, c.Y+j}] = used + 1
				} else {
					fabric[day3.Pos{c.X+i, c.Y+j}] = 1
				}
			}
		}

	}

	for _, v := range fabric {
		if v > 1 {
			sum = sum + 1
		}
	}

	fmt.Println(sum)

}



