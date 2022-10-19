package main

import "fmt"

func main() {
	x := make(map[int][]int, 0)
	x[1] = make([]int, 0)
	x[1] = append(x[1], 2)
	fmt.Println(x)
}
