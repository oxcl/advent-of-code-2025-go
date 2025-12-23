package main

import (
	"bufio"
	"os"
	"strings"
)

const FILENAME = "../sample-input.txt"

func main() {
	file, err := os.Open(FILENAME)
	if err != nil {
		println("there was an issue opening the file", FILENAME)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	width := len(firstLine)
	beamArray := make([]bool, width)
	indexOfS := strings.Index(firstLine, "S")
	if indexOfS == -1 {
		println("couldn't find S in the first line")
		return
	}
	beamArray[indexOfS] = true
	hitCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x := 0; x < len(line); x++ {
			if line[x] != '^' {
				continue
			}
			if beamArray[x] == false {
				// the beam will never hit this split so it's irrelevant
				continue
			}
			// here the beam WILL hit the split
			beamArray[x] = false
			hitCount++
			if x > 0 {
				beamArray[x-1] = true
			}
			if x < width-1 {
				beamArray[x+1] = true
			}
		}
	}
	println("Total: ", hitCount)
}
