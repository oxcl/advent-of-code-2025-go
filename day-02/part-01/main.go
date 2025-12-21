package main

import (
	"os"
	"strconv"
	"strings"
)

const Filename = "../input.txt"

var Result = 0

func main() {
	bytes, err := os.ReadFile(Filename)
	if err != nil {
		println("there was an issue reading the file", Filename)
		return
	}
	input := string(bytes)
	for _, idRange := range strings.Split(input, ",") {
		idRangeArr := strings.Split(idRange, "-")
		start, err := strconv.Atoi(idRangeArr[0])
		if err != nil {
			println("there was an issue converting the range to numbers", idRange)
			return
		}
		end, err := strconv.Atoi(idRangeArr[1])
		if err != nil {
			println("there was an issue converting the range to numbers", idRange)
			return
		}
		for i := start; i <= end; i++ {
			entry := strconv.Itoa(i)
			if len(entry)%2 != 0 {
				continue
			}
			middle := len(entry) / 2
			firstChunk := entry[:middle]
			secondChunk := entry[middle:]
			if firstChunk == secondChunk {
				// println("Invalid ID detected", entry, "Range: ", idRange, "Result:", Result)
				Result += i
			}
		}
	}
	println("Result:", Result)
}
