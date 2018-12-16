package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	recipes := make([]int, 2, 1000000)
	recipes[0] = 3
	recipes[1] = 7

	e1 := 0
	e2 := 1
	numRecipes := "190221"
	recipesStr := make([]byte, 2, 1000000)
	recipesStr[0] = '3'
	recipesStr[1] = '7'

	iter := 0

	for i := 0; i<20000000; i++{
		sumRec := recipes[e1] + recipes[e2]
		if sumRec < 10 {
			recipes = append(recipes, sumRec)
		} else {
			recipes = append(recipes, 1, sumRec-10)
		}
		recipesStr = append(recipesStr, []byte(strconv.FormatInt(int64(sumRec), 10))...)
		e1 = (e1 + 1 + recipes[e1])%len(recipes)
		e2 = (e2 + 1 + recipes[e2])%len(recipes)
		iter++
	}

	fmt.Println(strings.Index(string(recipesStr), numRecipes))

}
