package main

import (
	"bufio"
	"os"
	"strconv"
)

const Filename = "../input.txt"

var Password = 0
var Dial = 50

func main() {
	file, err := os.Open(Filename)
	if err != nil {
		println("there was an issue opening the file ", Filename)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		turn := line[0]
		distance, err := strconv.Atoi(line[1:])
		if err != nil {
			println("failed to parse distance number in line:", line)
			return
		}
		if turn == 'L' {
			Dial = (((Dial - distance) % 100) + 100) % 100
		} else {
			Dial = (Dial + distance) % 100
		}

		if Dial == 0 {
			Password++
		}
	}
	println("Password: ", Password)
}
