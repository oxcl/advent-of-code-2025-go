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
		for entryNumber := start; entryNumber <= end; entryNumber++ {
			entry := strconv.Itoa(entryNumber)
			doesMatch := false
			for divisor := 2; divisor <= len(entry); divisor++ {
				if len(entry)%divisor != 0 {
					continue
				}
				chunkSize := len(entry) / divisor
				areChunksEqual := true
				for i := 1; i < divisor; i++ {
					for j := 0; j < chunkSize; j++ {
						if entry[j] != entry[chunkSize*i+j] {
							areChunksEqual = false
							break
						}
					}
					if !areChunksEqual {
						break
					}
				}
				if areChunksEqual {
					doesMatch = true
					break
				}
			}
			if doesMatch {
				Result += entryNumber
			}
		}
	}
	println("Result:", Result)
}
