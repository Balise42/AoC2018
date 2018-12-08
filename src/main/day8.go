package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[1]
	var tree []int

	f, _ := os.Open(filename)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	for _, tok := range strings.Split(line, " ") {
		v, _ := strconv.ParseInt(tok, 10, 64)
		tree = append(tree, int(v))
	}

	res, _ := getMetadataSum(tree, 1, 0)
	fmt.Println(res)
	values, _ := getValues(tree, 1, 0)
	fmt.Println(values)
}

func getMetadataSum(tree []int, elems int, start int) (int, int) {
	sumMetadata := 0

	currentElementStart := start
	for i := 0; i < elems; i++ {
		numChildren := tree[currentElementStart]
		numMetaData := tree[currentElementStart+1]
		if numChildren == 0 {
			for j := 0; j < numMetaData; j++ {
				sumMetadata = sumMetadata + tree[currentElementStart+2+j]
			}
			currentElementStart = currentElementStart + 2 + numMetaData
		} else {
			var tmpSum int
			tmpSum, currentElementStart = getMetadataSum(tree, numChildren, currentElementStart+2)
			sumMetadata = sumMetadata + tmpSum
			for j := 0; j < numMetaData; j++ {
				sumMetadata = sumMetadata + tree[currentElementStart+j]
			}
			currentElementStart = currentElementStart + numMetaData
		}
	}
	return sumMetadata, currentElementStart
}

func getValues(tree []int, elems int, start int) ([]int, int) {
	values := make([]int, elems)
	for i := 0; i<elems; i++ {
		values[i] = 0
	}

	currentElementStart := start
	for i := 0; i<elems; i++ {
		numChildren := tree[currentElementStart]
		numMetaData := tree[currentElementStart+1]

		if numChildren == 0 {
			for j := 0; j < numMetaData; j++ {
				values[i] = values[i] + tree[currentElementStart+2+j]
			}
			currentElementStart = currentElementStart + 2 + numMetaData
		} else {
			var tmpValues []int
			tmpValues, currentElementStart = getValues(tree, numChildren, currentElementStart+2)
			for j := 0; j < numMetaData; j++ {
				index := tree[currentElementStart+j]
				if index > 0 && index <= len(tmpValues) {
					values[i] = values[i] + tmpValues[index-1]
				}
			}
			currentElementStart = currentElementStart + numMetaData
		}
	}
	return values, currentElementStart
}