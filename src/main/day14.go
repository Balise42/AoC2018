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
	endchar := numRecipes[len(numRecipes)-1]
//	numRecipes := "59414"
	recipesStr := make([]byte, 2, 1000000)
	recipesStr[0] = '3'
	recipesStr[1] = '7'

	iter := 0

	for {
		sumRec := recipes[e1] + recipes[e2]
		if sumRec < 10 {
			recipes = append(recipes, sumRec)
		} else {
			recipes = append(recipes, 1, sumRec-10)
		}
		recipesStr = append(recipesStr, []byte(strconv.FormatInt(int64(sumRec), 10))...)
		if recipesStr[len(recipesStr) - 1] == endchar {
			if strings.HasSuffix(string(recipesStr), numRecipes) || strings.HasSuffix(string(recipesStr[:len(recipesStr)-1]), numRecipes) {
				break
			}
		}
		e1 = (e1 + 1 + recipes[e1])%len(recipes)
		e2 = (e2 + 1 + recipes[e2])%len(recipes)
		iter++
		if iter % 10000 == 0 {
			fmt.Println(iter, len(recipesStr))
		}
	}

	fmt.Println(len(recipesStr) - len(numRecipes))

}
