package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"sort"
)

var guardIds map[string]int
var times map[string][]int


func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	guardIds = make(map[string]int)
	times = make(map[string][]int)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parse(scanner.Text())
	}

	awake := make(map[string][]bool)

	for date, minutes := range times {
		awake[date] = make([]bool, 60)
		for i := 0; i<60; i++ {
			awake[date][i] = true
		}
		sort.Ints(minutes)

		prev := 0
		for i, m := range minutes {
			if i%2 == 0 {
				prev = m
			} else {
				for j := prev; j<m; j++ {
					awake[date][j] = false
				}
			}
			if len(minutes) % 2 == 1 {
				for j := prev; j<60; j++ {
					awake[date][j] = false
				}
			}
		}
	}

	fmt.Println(awake)
}

func parse(line string) {
	tokens := strings.Split(line, " ")
	date := tokens[0]
	time := tokens[1]
	action := tokens[2]

	dateKey := getDateKey(date)
	timeValue := getTimeValue(time)

	if action == "falls" || action == "wakes" {
		times[dateKey] = append(times[dateKey], timeValue)
	} else {
		guardId := getGuardId(tokens[3])
		guardIds[dateKey] = guardId
	}
}

func getTimeValue(time string) int {
	minuteStr := strings.Split(time,":")[1]
	minuteStr = minuteStr[:len(minuteStr)-1]
	minute, _ := strconv.ParseInt(minuteStr, 10, 64)
	return int(minute)
}

func getDateKey(date string) string {
	return strings.SplitN(date, "-", 2)[1]
}

func getGuardId(guard string) int {
	guardId, _ := strconv.ParseInt(guard[1:], 10, 64)
	return int(guardId)
}