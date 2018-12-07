package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)



func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	graph := make(map[string][]string)
	parents := make(map[string][]string)
	vertices := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		v1, v2 := parseEdge(scanner.Text())
		graph[v1] = append(graph[v1], v2)
		parents[v2] = append(parents[v2], v1)
		vertices[v1] = true
		vertices[v2] = true
	}

	fmt.Println(topologicalSort(graph, vertices))
}

func topologicalSort(graph map[string][]string, vertices map[string]bool) string {
	candidates := make([]string, 0)
	for v, _ := range vertices {
		candidate := true
		for _, neighbors := range graph {
			for _, neighbor := range neighbors {
				if v == neighbor {
					candidate = false
					break
				}
			}
			if !candidate {
				break
			}
		}
		if candidate {
			candidates = append(candidates, v)
		}
	}

	topOrder := ""
	var v string
	for len(candidates) != 0 {
		fmt.Println(candidates)
		v, candidates = extractNext(candidates)
		candidates = appendListIfNotExists(candidates, graph[v])

		topOrder = topOrder + v
	}
	return topOrder
}

func appendListIfNotExists(list []string, elements []string) []string {
	for _, element := range elements {
		if !containsStr(list, element) {
			list = append(list, element)
		}
	}
	return list
}

func containsStr(list []string, element string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

func extractNext(candidates []string) (string, []string) {
	minStrIndex := 0
	minStr := candidates[0]
	for i, candidate := range candidates {
		if candidate < minStr {
			minStrIndex = i
			minStr = candidate
		}
	}
	return minStr, append(candidates[:minStrIndex], candidates[minStrIndex + 1:]...)
}

func parseEdge(line string) (string, string) {
	tokens := strings.Split(line, " ")
	return tokens[1], tokens[7]
}







