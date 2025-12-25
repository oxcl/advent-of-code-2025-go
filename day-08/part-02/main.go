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
const MAX_CONNECTIONS = 1000

type Node struct {
	X int
	Y int
	Z int
}

type Distance struct {
	Node1    *Node
	Node2    *Node
	Distance float64
}

func parseNodes() ([]Node, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		println()
		return nil, errors.New("there was a problem opening file: " + FILENAME)
	}
	defer file.Close()

	nodes := make([]Node, 0, 1000)

	scanner := bufio.NewScanner(file)

	// parse the cordinates of the nodes from the file
	for scanner.Scan() {
		line := scanner.Text()
		cords := strings.Split(line, ",")
		if len(cords) != 3 {
			return nil, errors.New("cords length must be 3. recieved: " + line)
		}
		x, err := strconv.Atoi(cords[0])
		if err != nil {
			return nil, errors.New("failed to parse the cordinates in this string: " + line)
		}
		y, err := strconv.Atoi(cords[1])
		if err != nil {
			return nil, errors.New("failed to parse the cordinates in this string: " + line)
		}
		z, err := strconv.Atoi(cords[2])
		if err != nil {
			return nil, errors.New("failed to parse the cordinates in this string: " + line)
		}
		node := Node{x, y, z}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func calculateDistances(nodes []Node) []Distance {
	numberOfDistances := len(nodes) * (len(nodes) - 1) / 2
	distances := make([]Distance, numberOfDistances)
	index := 0
	var n1, n2 *Node
	for i := 0; i < len(nodes)-1; i++ {
		n1 = &nodes[i]
		for j := i + 1; j < len(nodes); j++ {
			n2 = &nodes[j]
			distances[index] = Distance{
				Node1:    n1,
				Node2:    n2,
				Distance: math.Sqrt(float64((n1.X-n2.X)*(n1.X-n2.X) + (n1.Y-n2.Y)*(n1.Y-n2.Y) + (n1.Z-n2.Z)*(n1.Z-n2.Z))),
			}
			index++
		}
	}
	return distances
}

type Circuit []*Node

func main() {
	nodes, err := parseNodes()
	if err != nil {
		println("there was an issue parsing the nodes. error:\n", err)
		return
	}
	distances := calculateDistances(nodes)
	slices.SortFunc(distances, func(d1, d2 Distance) int {
		if d1.Distance > d2.Distance {
			return 1
		} else if d1.Distance < d2.Distance {
			return -1
		}
		return 0
	})
	connections := 0
	var circuits []*Circuit
	connectedNodes := make(map[*Node]*Circuit)
	last_index := 0
	for index, distance := range distances {
		if len(circuits) > 0 && len(*circuits[0]) >= len(nodes) {
			break
		}
		last_index = index
		n1Circuit, isN1Connected := connectedNodes[distance.Node1]
		n2Circuit, isN2Connected := connectedNodes[distance.Node2]
		if !isN1Connected && !isN2Connected {
			// create a new circuit and add the two nodes
			circuit := Circuit{distance.Node1, distance.Node2}
			connectedNodes[distance.Node1] = &circuit
			connectedNodes[distance.Node2] = &circuit
			circuits = append(circuits, &circuit)
			// connections++
		} else if isN1Connected && !isN2Connected {
			// connect n2 to n1's circuit
			*n1Circuit = append(*n1Circuit, distance.Node2)
			connectedNodes[distance.Node2] = n1Circuit
			// connections++
		} else if !isN1Connected && isN2Connected {
			// connect n1 to n2's circuit
			*n2Circuit = append(*n2Circuit, distance.Node1)
			connectedNodes[distance.Node1] = n2Circuit
			// connections++
		} else if isN1Connected && isN2Connected && n1Circuit != n2Circuit {
			// connect the two circuits together
			mergedCircuit := append(*n1Circuit, *n2Circuit...)
			for _, node := range mergedCircuit {
				connectedNodes[node] = &mergedCircuit
			}
			// replace the old circuits with the merged circuit
			newCircuits := []*Circuit{&mergedCircuit}
			for _, circuit := range circuits {
				if circuit != n1Circuit && circuit != n2Circuit {
					newCircuits = append(newCircuits, circuit)
				}
			}
			circuits = newCircuits
			// connections++
		}
		connections++
	}
	println("Answer: ", distances[last_index].Node1.X*distances[last_index].Node2.X)
}
