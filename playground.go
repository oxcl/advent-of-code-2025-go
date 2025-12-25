package main

func modifySlice(slice [3]int) {
	slice[0] = 1
	slice[1] = 2
	slice[2] = 3
}

func main() {

	slice := [3]int{0, 1, 2}
	modifySlice(slice)
	for _, num := range slice {
		println(num)
	}
}
