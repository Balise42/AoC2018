package main

import (
	"os"
	"bufio"
	"fmt"
	"math"
	"strings"
)

var offset = 60
var workers = 5

func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	graph := make(map[string][]string)
	parents := make(map[string][]string)
	vertices := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		v1, v2 := parseEdgeb(scanner.Text())
		graph[v1] = append(graph[v1], v2)
		parents[v2] = append(parents[v2], v1)
		vertices[v1] = true
		vertices[v2] = true
	}

	fmt.Println(processDuration(graph, vertices, parents))
}
func processDuration(graph map[string][]string, vertices map[string]bool, parents map[string][]string) int {
	dependencies := make(map[string][]string)
	for k, v := range parents {
		dependencies[k] = make([]string, 0)
		for _, u := range v {
			dependencies[k] = append(dependencies[k], u)
		}
	}

	completionDate := make(map[string]int)
	startDate := make(map[string]int)
	workerNextSpot := make([]int, workers)

	candidates := make(map[string]int)
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
			candidates[v] = 0
			startDate[v] = -1
		}
	}


	for len(candidates) != 0 {
		fmt.Println(candidates)
		w, date := getWorkerStartTime(workerNextSpot)
		var v string
		var startTask int
		v, candidates, startTask = getNextInOrder(candidates, date)
		if startTask < date {
			startDate[v] = date
		} else {
			startDate[v] = startTask
		}
		completionDate[v] = startDate[v] + offset + int(byte(v[0]) - 'A')
		workerNextSpot[w] = completionDate[v] + 1

		for _, next := range graph[v] {
			parents[next] = removeEdge(parents[next], v)
			if len(parents[next]) == 0 {
				candidates[next] = getTaskStartTime(completionDate, dependencies[next])
			}
		}


	}

	res := 0
	for _, time := range completionDate {
		if time > res {
			res = time
		}
	}
	return res + 1
}

func parseEdgeb(line string) (string, string) {
	tokens := strings.Split(line, " ")
	return tokens[1], tokens[7]
}

func getTaskStartTime(completionDates map[string]int, parents []string) int {
	maxTime := 0
	for _, t := range parents {
		if completionDates[t] > maxTime {
			maxTime = completionDates[t]
		}
	}
	return maxTime + 1
}

func getWorkerStartTime(workers []int) (int, int) {
	minTime := math.MaxInt32
	worker := -1
	for i, t := range workers {
		if t < minTime {
			minTime = t
			worker = i
		}
	}
	return worker, minTime
}
func getNextInOrder(candidates map[string]int, date int) (string, map[string]int, int) {
	var candidate = "a"
	var startTime = math.MaxInt64

	for k, v := range candidates {
		if v <= date && k < candidate {
			candidate = k
			startTime = v
		}
	}

	if candidate != "a" {
		delete(candidates, candidate)
		return candidate, candidates, startTime
	}

	for k, v := range candidates {
		if v < startTime {
			candidate = k
			startTime = v
		} else if v == startTime && k < candidate {
			candidate = k
		}
	}
	delete(candidates, candidate)
	return candidate, candidates, startTime
}

func removeEdge(edges []string, v string) []string {
	for i, u := range edges {
		if u == v {
			return append(edges[:i], edges[i+1:]...)
		}
	}
	return edges
}



