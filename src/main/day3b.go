package main

import (
	"os"
	"common/day3"
	"bufio"
	"fmt"
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

	for i := 0; i < len(claims); i++ {
		intersected := false
		for j  := 0; j <len(claims); j++ {
			if i != j && intersects(claims[i], claims[j]) {
				intersected = true
				break
			}
		}
		if !intersected {
			fmt.Println(i + 1)
		}
	}
}

func intersects(c1 day3.Claim, c2 day3.Claim) bool {
	return intersects1D(c1.X, c2.X, c1.Width, c2.Width) && intersects1D(c1.Y, c2.Y, c1.Height, c2.Height)
}

func intersects1D(z1 int, z2 int, size1 int, size2 int) bool {
	var x1, x2, w1 int
	if z1 < z2 {
		x1 = z1
		x2 = z2
		w1 = size1
	} else {
		x1 = z2
		x2 = z1
		w1 = size2
	}

	if x1 + w1 - 1 < x2 {
		return false
	}
	return true
}





