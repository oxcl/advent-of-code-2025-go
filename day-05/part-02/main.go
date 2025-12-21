package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const Filename = "../input.txt"

func reduceRangeUntilUnique(theRange []int, ranges [][]int) [][]int {
	start := theRange[0]
	end := theRange[1]
	for _, aRange := range ranges {
		if start >= aRange[0] && end <= aRange[1] {
			// the range is redundant
			return ranges
		}
		if start >= aRange[0] && start <= aRange[1] {
			start = aRange[1] + 1
		}
		if end >= aRange[0] && end <= aRange[1] {
			end = aRange[0] - 1
		}
		if start <= aRange[0] && end >= aRange[1] {
			// the range must be splitted into two parts
			ranges = reduceRangeUntilUnique([]int{start, aRange[0] - 1}, ranges)
			start = aRange[1] + 1
		}
	}
	return append(ranges, []int{start, end})
}

func main() {
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

	// analyze ranges to convert them into a list of non intersecting ranges
	newRanges := make([][]int, 0)
	for _, theRange := range ranges {
		newRanges = reduceRangeUntilUnique(theRange, newRanges)
	}

	result := 0
	for _, theRange := range newRanges {
		result += theRange[1] - theRange[0] + 1
	}

	println("Result:", result)
}
