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

const FILENAME = "../sample-input.txt"

type Node struct {
	X int
	Y int
}

type Walls struct {
	XMap map[int][2]int
	YMap map[int][2]int
}

type Area struct {
	Node1 *Node
	Node2 *Node
	Area  float64
}

func NewWalls() Walls {
	return Walls{
		XMap: make(map[int][2]int),
		YMap: make(map[int][2]int),
	}
}

func (walls Walls) isNodeOnWall(point Node) bool {
	if wallRange, exists := walls.XMap[point.X]; exists {
		if point.Y >= wallRange[0] && point.Y <= wallRange[1] {
			return true
		} else {
			return false
		}
	} else if wallRange, exists := walls.YMap[point.Y]; exists {
		if point.Y >= wallRange[0] && point.Y <= wallRange[1] {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
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
				Area:  (math.Abs(float64(n1.X)-float64(n2.X)) + 1) * (math.Abs(float64(n1.Y)-float64(n2.Y)) + 1),
			}
			index++
		}
	}
	return areas
}

func calculateWalls(nodes []Node) Walls {
	walls := NewWalls()
	wrappedNodes := append(nodes, nodes[0])
	for i := 0; i < len(wrappedNodes)-1; i++ {
		node := wrappedNodes[i]
		nextNode := wrappedNodes[i+1]
		if node.X == nextNode.X {
			min := int(math.Min(float64(node.Y), float64(nextNode.Y)))
			max := int(math.Min(float64(node.Y), float64(nextNode.Y)))
			walls.XMap[node.X] = [2]int{min, max}
		}
		if node.Y == nextNode.Y {
			min := int(math.Min(float64(node.X), float64(nextNode.X)))
			max := int(math.Min(float64(node.X), float64(nextNode.X)))
			walls.YMap[node.Y] = [2]int{min, max}
		}
	}
	return walls
}

func main() {
	nodes, err := parseNodes()
	if err != nil {
		println("there was an issue parsing the nodes. error:\n", err)
		return
	}
	walls := calculateWalls(nodes)
	areas := calculateAreas(nodes)
	slices.SortFunc(areas, func(d1, d2 Area) int {
		if d1.Area < d2.Area {
			return 1
		} else if d1.Area > d2.Area {
			return -1
		}
		return 0
	})

	for _, area := range areas {
		x1 := area.Node1.X
		y1 := area.Node1.Y
		x2 := area.Node2.X
		y2 := area.Node2.Y
		if walls.isNodeOnWall(Node{X: x1, Y: y1}) && walls.isNodeOnWall(Node{X: x2, Y: y2}) && walls.isNodeOnWall(Node{X: x1, Y: y2}) && walls.isNodeOnWall(Node{X: x2, Y: y1}) {
			println("Largest Area is:", int(area.Area))
			println("For Nodes: (", x1, ",", y1, ") and (", x2, ",", y2, ")")
			return
		}
	}
	println("couldn't find the largest area!")

	// width := math.Abs(float64(largestArea.Node1.X)-float64(largestArea.Node2.X)) + 1
	// height := math.Abs(float64(largestArea.Node1.Y)-float64(largestArea.Node2.Y)) + 1
	// println("Answer:", int(width*height))
}
