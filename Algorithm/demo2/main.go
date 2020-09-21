package main

import "fmt"

const (
	a = iota
	b = iota
)
const (
	xx = ""
	c  = iota
	d  = iota
)

func main() {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]
	str2[1] = "new"
	fmt.Println(str1)
	str2 = append(str2, "z", "x", "y")
	fmt.Println(str1)
}
