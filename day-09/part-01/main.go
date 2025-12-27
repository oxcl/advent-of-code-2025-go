package main

import (
	"bufio"
	"errors"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const FILENAME = "../input.txt"

type Node struct {
	X int
	Y int
}

type Area struct {
	Node1 *Node
	Node2 *Node
	Area  float64
}

func parseNodes() ([]Node, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		println()
		return nil, errors.New("there was a problem opening file: " + FILENAME)
	}
	defer file.Close()

	nodes := make([]Node, 0, 500)

	scanner := bufio.NewScanner(file)

	// parse the cordinates of the nodes from the file
	for scanner.Scan() {
		line := scanner.Text()
		cords := strings.Split(line, ",")
		if len(cords) != 2 {
			return nil, errors.New("cords length must be 2. recieved: " + line)
		}
		x, err := strconv.Atoi(cords[0])
		if err != nil {
			return nil, errors.New("failed to parse the cordinates in this string: " + line)
		}
		y, err := strconv.Atoi(cords[1])
		if err != nil {
			return nil, errors.New("failed to parse the cordinates in this string: " + line)
		}
		node := Node{X: x, Y: y}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func calculateAreas(nodes []Node) []Area {
	numberOfAreas := len(nodes) * (len(nodes) - 1) / 2
	areas := make([]Area, numberOfAreas)
	index := 0
	var n1, n2 *Node
	for i := 0; i < len(nodes)-1; i++ {
		n1 = &nodes[i]
		for j := i + 1; j < len(nodes); j++ {
			n2 = &nodes[j]
			areas[index] = Area{
				Node1: n1,
				Node2: n2,
				Area:  math.Abs(float64(n1.X)-float64(n2.X)) * math.Abs(float64(n1.Y)-float64(n2.Y)),
			}
			index++
		}
	}
	return areas
}

type Circuit []*Node

func main() {
	nodes, err := parseNodes()
	if err != nil {
		println("there was an issue parsing the nodes. error:\n", err)
		return
	}
	areas := calculateAreas(nodes)
	largestArea := slices.MaxFunc(areas, func(d1, d2 Area) int {
		if d1.Area > d2.Area {
			return 1
		} else if d1.Area < d2.Area {
			return -1
		}
		return 0
	})

	width := math.Abs(float64(largestArea.Node1.X)-float64(largestArea.Node2.X)) + 1
	height := math.Abs(float64(largestArea.Node1.Y)-float64(largestArea.Node2.Y)) + 1
	println("Answer:", int(width*height))
}
