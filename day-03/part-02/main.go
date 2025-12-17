package main

import (
	"bufio"
	"math"
	"os"
)

func main() {
	file, err := os.Open("../input.txt")

	if err != nil {
		println("there was an error opening the input file.", err)
		return
	}

	scanner := bufio.NewScanner(file)
	sum := 0.0
	for scanner.Scan() {
		battery := scanner.Text()
		battery_sum := 0.0
		min_j := 0
		for i := 1; i <= 12; i++ {
			largest_digit := -1
			for j := min_j; j < len(battery)-(12-i); j++ {
				digit := int(battery[j] - '0')
				if digit > largest_digit {
					largest_digit = digit
					min_j = j + 1
				}
			}
			battery_sum += float64(largest_digit) * math.Pow10(12-i)
		}
		sum += battery_sum
	}
	println("result: ", int(sum))
}
