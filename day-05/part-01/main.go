package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const Filename = "../input.txt"

func main() {
	result := 0
	file, err := os.Open(Filename)
	if err != nil {
		println("there was an issue opening the file", Filename)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ranges := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Trim(line, " \r\n\t") == "" {
			break
		}
		numberRange := strings.Split(line, "-")
		start, err := strconv.Atoi(numberRange[0])
		if err != nil {
			println("there was an issue parsing this range to int", line)
			return
		}
		end, err := strconv.Atoi(numberRange[1])
		if err != nil {
			println("there was an issue parsing this range to int", line)
			return
		}
		ranges = append(ranges, []int{start, end})
	}

	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)
		if err != nil {
			println("there was an issue parsing this entry number to int", line)
			return
		}
		for _, theRange := range ranges {
			if number >= theRange[0] && number <= theRange[1] {
				result++
				break
			}
		}
	}

	println("Result: ", result)
}
