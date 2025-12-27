package main

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func FindNodesGCD(nodes []Node) (int, int) {
	currentX, currentY := nodes[0].X, nodes[0].Y
	for _, node := range nodes[1:] {
		currentX = GCD(currentX, node.X)
		currentY = GCD(currentY, node.Y)
		if currentX == 1 && currentY == 1 {
			return 1, 1
		}
	}
	return currentX, currentY
}
