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

type Walls struct {
	XMap   map[int][2]int
	YMap   map[int][2]int
	Width  int
	Height int
}

func calculateWalls(nodes []Node) Walls {
	walls := Walls{
		XMap: make(map[int][2]int),
		YMap: make(map[int][2]int),
	}
	wrappedNodes := append(nodes, nodes[0])

	for i := 0; i < len(wrappedNodes)-1; i++ {
		node := wrappedNodes[i]
		nextNode := wrappedNodes[i+1]
		if node.X == nextNode.X {
			min := int(math.Min(float64(node.Y), float64(nextNode.Y)))
			max := int(math.Max(float64(node.Y), float64(nextNode.Y)))
			walls.XMap[node.X] = [2]int{min, max}
		}
		if node.Y == nextNode.Y {
			min := int(math.Min(float64(node.X), float64(nextNode.X)))
			max := int(math.Max(float64(node.X), float64(nextNode.X)))
			walls.YMap[node.Y] = [2]int{min, max}
		}
	}
	maxX := 0
	maxY := 0
	for i := 0; i < len(nodes); i++ {
		maxX = int(math.Max(float64(nodes[i].X), float64(maxX)))
		maxY = int(math.Max(float64(nodes[i].Y), float64(maxY)))
	}
	walls.Width = maxX + 1
	walls.Height = maxY + 1
	return walls
}

func (walls Walls) isNodeInsideWalls(point Node) bool {
	// check if there's a wall at the top of the point
	hitWall := false
	for wallY, xRange := range walls.YMap {
		wallX1, wallX2 := xRange[0], xRange[1]
		if wallY <= point.Y && point.X >= wallX1 && point.X <= wallX2 {
			hitWall = true
			break
		}
	}
	if !hitWall {
		return false
	}
	hitWall = false
	// check if there's a wall at the bottom of the point
	for wallY, xRange := range walls.YMap {
		wallX1, wallX2 := xRange[0], xRange[1]
		if wallY >= point.Y && point.X >= wallX1 && point.X <= wallX2 {
			hitWall = true
			break
		}
	}
	if !hitWall {
		return false
	}
	// check if there's a wall on the right of the point
	hitWall = false
	for wallX, yRange := range walls.XMap {
		wallY1, wallY2 := yRange[0], yRange[1]
		if wallX >= point.X && point.Y >= wallY1 && point.Y <= wallY2 {
			hitWall = true
			break
		}
	}
	if !hitWall {
		return false
	}
	hitWall = false
	// check if there's a wall at the bottom of the point
	for wallX, yRange := range walls.XMap {
		wallY1, wallY2 := yRange[0], yRange[1]
		if wallX <= point.X && point.Y >= wallY1 && point.Y <= wallY2 {
			hitWall = true
			break
		}
	}
	if !hitWall {
		return false
	}
	return true
}

type Area struct {
	Node1 *Node
	Node2 *Node
	Area  float64
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
			width := math.Abs(float64(n1.X)-float64(n2.X)) + 1
			height := math.Abs(float64(n1.Y)-float64(n2.Y)) + 1
			areas[index] = Area{
				Node1: n1,
				Node2: n2,
				Area:  width * height,
			}
			index++
		}
	}
	return areas
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
		// if one of the corners of the rectangle is outisde the walls then the whole rectangle is not inside the walls (duh)
		if !walls.isNodeInsideWalls(Node{X: x1, Y: y2}) || !walls.isNodeInsideWalls(Node{X: x2, Y: y1}) {
			continue
		}

		// even if all the four corners are inside the walls we got to still check if no wall is carving inside the borders
		// meaning that we check to see if all the points on the borders of the rectangle are inside the walls
		// we do so by checking if any of the borders either have one of their points inside the rectangle
		// or their range extends the whole size of the rectangle
		minX := int(math.Min(float64(x1), float64(x2)))
		minY := int(math.Min(float64(y1), float64(y2)))
		maxX := int(math.Max(float64(x1), float64(x2)))
		maxY := int(math.Max(float64(y1), float64(y2)))
		didFail := false
		for wallY, xRange := range walls.YMap {
			wallX1, wallX2 := xRange[0], xRange[1]
			if wallY <= minY || wallY >= maxY {
				// the wall is outside of the rectangle. early exit.
				continue
			}
			if (wallX1 > minX && wallX1 < maxX) || (wallX2 > minX && wallX2 < maxX) {
				// this means that one of the edges of a wall is inside the rectangle
				didFail = true
				break
			}
			if wallX1 <= minX && wallX2 >= maxX {
				// this means that a wall is bigger than the rectangle and it's going through
				didFail = true
				break
			}
		}
		if didFail {
			continue
		}
		for wallX, yRange := range walls.XMap {
			wallY1, wallY2 := yRange[0], yRange[1]
			if wallX <= minX || wallX >= maxX {
				// the wall is outside of the rectangle. early exit.
				continue
			}
			if (wallY1 > minY && wallY1 < maxY) || (wallY2 > minY && wallY2 < maxY) {
				// this means that one of the edges of a wall is inside the rectangle
				didFail = true
				break
			}
			if wallY1 <= minY && wallY2 >= maxY {
				// this means that a wall is bigger than the rectangle and it's going through
				didFail = true
				break
			}
		}
		if didFail {
			continue
		}

		println("Found it:", int(area.Area))
		return
	}
	println("couldn't find the largest area!")
}
