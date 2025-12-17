package main

import (
	"bufio"
	"os"
)

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		println("there was an error opening the input file.", err)
		return
	}

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		battery := scanner.Text()
		largest := 0
		for i := 0; i < len(battery)-1; i++ {
			a := int(battery[i] - '0')
			for j := i + 1; j < len(battery); j++ {
				b := int(battery[j] - '0')
				largest = max(largest, a*10+b)
			}
		}
		sum += largest
	}
	println("result: ", sum)
}
