package LineProcessor

import (
	"bufio"
	"os"
)

func SumLines(filename string, processor func(string) int64) int64 {
	file, _ := os.Open(filename)
	defer file.Close()

	var result int64
	result = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = result + processor(scanner.Text())
	}
	return result
}
