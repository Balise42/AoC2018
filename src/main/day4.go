package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"sort"
	"time"
)

var guardIds map[time.Time]int
var times map[time.Time][]int
var asleep map[time.Time][]int


func main() {
	filename := os.Args[1]
	f, _ := os.Open(filename)
	defer f.Close()

	guardIds = make(map[time.Time]int)
	times = make(map[time.Time][]int)
    asleep = make(map[time.Time][]int)



	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parse(scanner.Text())
	}

	computeAsleep()

	sleepy := getMostAsleepGuard()
	asleepTime := getMostAsleepTime(sleepy)

	fmt.Println(sleepy * asleepTime)

	sleepyMinute := getMostSleepyMinute()
	asleepGuard := getMostAsleepGuardAtMinute(sleepyMinute)

	fmt.Println(sleepyMinute * asleepGuard)
}

func getMostAsleepGuardAtMinute(sleepyMinute int) int {
	countAsleep := make(map[int]int)

	for k, v := range asleep {
		if guardIds[k] == 0 {
			fmt.Println(k)
		}
		countAsleep[guardIds[k]] = countAsleep[guardIds[k]] + v[sleepyMinute]
	}

	fmt.Println(countAsleep)

	maxTime := 0
	sleepyId := -1

	for guard, time := range countAsleep {
		if time > maxTime {
			maxTime = time
			sleepyId = guard
		}
	}
	return sleepyId
}

func getMostSleepyMinute() int {
	countAsleep := make([]int, 60)
	for _, v := range asleep {
		addArrays(countAsleep, v)
	}
	return getMaxIndex(countAsleep)
}

func getMostAsleepTime(guard int) int {
	countAsleep := make([]int, 60)
	for k, v := range guardIds {
		if v == guard {
			addArrays(countAsleep, asleep[k])
		}
	}
	return getMaxIndex(countAsleep)
}

func getMaxIndex(countAsleep []int) int {
	maxMinutes := 0
	maxIndex := -1
	for i, val := range countAsleep {
		if val > maxMinutes {
			maxIndex = i
			maxMinutes = val
		}
	}
	return maxIndex
}

func addArrays(dest []int, increment []int) {
	for i, _ := range dest {
		dest[i] = dest[i] + increment[i]
	}
}



func getMostAsleepGuard() int {
	asleepAmount := make(map[int]int)
	for date, guard := range guardIds {
		asleepAmount[guard] += sum(asleep[date])
	}

	maxTime := 0
	sleepyId := -1

	for guard, time := range asleepAmount {
		if time > maxTime {
			maxTime = time
			sleepyId = guard
		}
	}
	return sleepyId
}
func sum(ints []int) int {
	sum := 0
	for _, v := range ints {
		sum = sum + v
	}
	return sum
}

func parse(line string) {
	tokens := strings.Split(line, " ")
	date := tokens[0]
	hour := tokens[1]
	action := tokens[2]

	dateKey := getDateKey(date[1:])
	timeValue := getTimeValue(hour)

	if action == "falls" || action == "wakes" {
		times[dateKey] = append(times[dateKey], timeValue)
	} else {
		guardId := getGuardId(tokens[3])
		if(timeValue > 50) {
			dateKey = incrementDateKey(dateKey)
		}
		guardIds[dateKey] = guardId
	}
}
func incrementDateKey(t time.Time) time.Time {
	return t.AddDate(0,0,1)
}



func getTimeValue(time string) int {
	minuteStr := strings.Split(time,":")[1]
	minuteStr = minuteStr[:len(minuteStr)-1]
	minute, _ := strconv.ParseInt(minuteStr, 10, 64)
	return int(minute)
}

func getDateKey(date string) time.Time {
	parsedDate, _ := time.Parse("2006-01-02", date)
	return parsedDate
}

func getGuardId(guard string) int {
	guardId, _ := strconv.ParseInt(guard[1:], 10, 64)
	return int(guardId)
}

func computeAsleep() {
	for date, minutes := range times {
		asleep[date] = make([]int, 60)
		for i := 0; i<60; i++ {
			asleep[date][i] = 0
		}
		sort.Ints(minutes)

		prev := 0
		for i, m := range minutes {
			if i%2 == 0 {
				prev = m
			} else {
				for j := prev; j<m; j++ {
					asleep[date][j] = 1
				}
			}
			if len(minutes) % 2 == 1 {
				for j := prev; j<60; j++ {
					asleep[date][j] = 1
				}
			}
		}
	}
}