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

func reduceOffset(nodes []Node) (offsetX int, offsetY int, reducedNodes []Node) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	for _, node := range nodes {
		minX = math.Min(minX, float64(node.X))
		minY = math.Min(minY, float64(node.Y))
	}
	offsetX = int(minX)
	offsetY = int(minY)
	reducedNodes = make([]Node, len(nodes))
	for index, node := range nodes {
		reducedNodes[index] = Node{X: node.X - offsetX, Y: node.Y - offsetY}
	}
	return offsetX, offsetY, reducedNodes
}

func compressNodes(nodes []Node) (float64, float64, []Node) {
	gcdX, gcdY := FindNodesGCD(nodes)
	if gcdX == 1 && gcdY == 1 {
		return 1, 1, nodes
	}
	newNodes := make([]Node, len(nodes), len(nodes))
	for index, node := range nodes {
		newNodes[index] = Node{
			X: node.X / gcdX,
			Y: node.Y / gcdY,
		}
	}
	return 1.0 / float64(gcdX), 1.0 / float64(gcdY), newNodes
}

type Walls struct {
	XMap   map[int][2]int
	YMap   map[int][2]int
	Width  int
	Height int
	cache  map[int]bool
}

func calculateWalls(nodes []Node) Walls {
	walls := Walls{
		XMap:  make(map[int][2]int),
		YMap:  make(map[int][2]int),
		cache: make(map[int]bool),
	}
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
	width := 0
	height := 0
	for i := 0; i < len(nodes); i++ {
		width = int(math.Max(float64(nodes[i].X), float64(width)))
		height = int(math.Max(float64(nodes[i].Y), float64(height)))
	}
	walls.Width = width
	walls.Height = height
	return walls
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

var directions = [4][2]int{
	{0, -1},
	{0, 1},
	{1, 0},
	{-1, 0},
}

func (walls Walls) isNodeInsideWalls(point Node) bool {
	if cacheHit, exists := walls.cache[point.Y*walls.Width+point.X]; exists {
		return cacheHit
	}
	for _, direction := range directions {
		dx := direction[0]
		dy := direction[1]
		didHitWall := false
		for x, y := point.X, point.Y; x >= 0 && x <= walls.Width && y >= 0 && y <= walls.Height; x, y = x+dx, y+dy {
			if cacheHit, exists := walls.cache[y*walls.Width+x]; exists {
				for cx, cy := point.X, point.Y; cx != x && cy != y; cx, cy = cx+dx, cy+dy {
					walls.cache[cy*walls.Width+cx] = cacheHit
				}
				return cacheHit
			}
			if walls.isNodeOnWall(Node{x, y}) {
				didHitWall = true
				break
			}
		}
		if !didHitWall {
			walls.cache[point.Y*walls.Width+point.X] = false
			return false
		}
	}
	walls.cache[point.Y*walls.Width+point.X] = true
	return true
}

type Area struct {
	Node1 *Node
	Node2 *Node
	Area  float64
}

func calculateAreas(nodes []Node, compressX float64, compressY float64) []Area {
	numberOfAreas := len(nodes) * (len(nodes) - 1) / 2
	areas := make([]Area, numberOfAreas)
	index := 0
	var n1, n2 *Node
	for i := 0; i < len(nodes)-1; i++ {
		n1 = &nodes[i]
		for j := i + 1; j < len(nodes); j++ {
			n2 = &nodes[j]
			width := (math.Abs(float64(n1.X)-float64(n2.X)) / compressX) + 1
			height := (math.Abs(float64(n1.Y)-float64(n2.Y)) / compressY) + 1
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
	_, _, nodes = reduceOffset(nodes)
	compressX, compressY, nodes := compressNodes(nodes)
	walls := calculateWalls(nodes)
	// width, height, nodePairs := pairNodes(nodes)
	// grid := generateGrid(width, height, nodes)
	areas := calculateAreas(nodes, compressX, compressY)
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
		var dx int
		var dy int
		if x1 < x2 {
			dx = 1
		} else {
			dx = -1
		}
		if y1 < y2 {
			dy = 1
		} else {
			dy = -1
		}
		didFail := false
		for _, x := range [2]int{x1, x2} {
			for y := y1; y != y2; y += dy {
				if !walls.isNodeInsideWalls(Node{x, y}) {
					didFail = true
					break
				}
			}
			if didFail {
				break
			}
		}
		if didFail {
			continue
		}
		for _, y := range [2]int{y1, y2} {
			for x := x1; y != x2; x += dx {
				if !walls.isNodeInsideWalls(Node{x, y}) {
					didFail = true
					break
				}
			}
			if didFail {
				break
			}
		}
		if didFail {
			continue
		}

		print("Found it!:", area.Area)
	}
	println("couldn't find the largest area!")
}
