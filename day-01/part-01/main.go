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

		if distance >= 100 {
			Password += distance / 100
			distance = distance % 100
		}

		prevDial := Dial

		if turn == 'L' {
			Dial = (((Dial - distance) % 100) + 100) % 100
		} else {
			Dial = (Dial + distance) % 100
		}

		if Dial == 0 {
			Password++
		} else if prevDial != 0 && turn == 'L' && Dial > prevDial {
			Password++
		} else if prevDial != 0 && turn == 'R' && Dial < prevDial {
			Password++
		}

		// println(prevDial, " + ", line, " -> ", Dial, ", Password:", Password)
	}
	println("Password: ", Password)
}
