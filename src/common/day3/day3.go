package day3

import (
	"strings"
	"strconv"
	"os"
)

type Claim struct {
	X int
	Y int
	Width int
	Height int
}

type Pos struct {
	X int
	Y int
}

func CreateClaim(s string) Claim {
	tokens := strings.Split(s, " ")
	pos := strings.Split(tokens[2], ",")
	pos[1] = strings.Split(pos[1], ":")[0]
	size := strings.Split(tokens[3], "x")

	return Claim{parseInt(pos[0]), parseInt(pos[1]), parseInt(size[0]), parseInt(size[1])}
}

func parseInt(s string) int {
	a, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		os.Exit(-1)
	}
	return int(a)
}

