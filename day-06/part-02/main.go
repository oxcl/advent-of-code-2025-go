package main

import (
	"bufio"
	"math"
	"os"
)

const Filename = "../input.txt"

func main() {
	file, err := os.Open(Filename)
	if err != nil {
		println("there was an issue opening the file", Filename)
		return
	}
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	length := 0
	for scanner.Scan() {
		line := scanner.Text()
		if length == 0 {
			length = len(line)
		}

		if length != len(line) {
			println("number of clumns in each line are not equal!")
			return
		}
		lines = append(lines, line)
	}
	var operator byte = 0
	total := 0
	totalColumn := 0
	for column := 0; column < length; column++ {
		newOperator := lines[len(lines)-1][column]
		if newOperator == '*' || newOperator == '+' {
			operator = newOperator
		} else if operator == 0 {
			println("failed to get the operator for the first column of numbers!")
			return
		}
		number := 0
		pow := 0
		for row := len(lines) - 2; row >= 0; row-- {
			item := lines[row][column]
			if item == ' ' {
				continue
			}
			digit := int(item - '0')
			number += digit * int(math.Pow10(pow))
			pow++
		}
		if number == 0 {
			total += totalColumn
			totalColumn = 0
		} else if operator == '*' {
			if totalColumn == 0 {
				totalColumn = 1
			}
			totalColumn *= number
		} else if operator == '+' {
			totalColumn += number
		} else {
			println("invalid operator!", operator)
			return
		}
	}
	total += totalColumn
	println("Total:", total)
}
