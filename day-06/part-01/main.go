package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const Filename = "../input.txt"

type NumberPair struct {
	mult int
	sum  int
}

func main() {
	file, err := os.Open(Filename)
	if err != nil {
		println("there was an issue opening the file", Filename)
		return
	}
	scanner := bufio.NewScanner(file)
	length := 0
	numberPairs := make([]NumberPair, 0)
	total := 0

	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("[ \t]+")
		items := re.Split(line, -1)
		if strings.Trim(items[len(items)-1], " \t\n\r") == "" {
			items = items[0 : len(items)-1]
		}
		if strings.Trim(items[0], " \t\n\r") == "" {
			items = items[1:]
		}
		if length == 0 {
			length = len(items)
			numberPairs = make([]NumberPair, length)
			for i := 0; i < length; i++ {
				numberPairs[i].sum = 0
				numberPairs[i].mult = 1
			}
		}
		if len(items) != length {
			println("number of columns in the input or not euqal", length, " != ", len(items))
			return
		}
		if items[0] == "*" || items[0] == "+" {
			for i := 0; i < length; i++ {
				switch items[i] {
				case "*":
					total += numberPairs[i].mult
				case "+":
					total += numberPairs[i].sum
				default:
					println("invalid operator", items[i])
					return
				}
			}
			break
		}
		for i := 0; i < length; i++ {
			number, err := strconv.Atoi(strings.Trim(items[i], " \t\r\n"))
			if err != nil {
				println("therre was an issue with converting the entry to number", items[i])
				return
			}
			numberPairs[i].sum += number
			numberPairs[i].mult *= number
		}
	}
	println("Total: ", total)
}
