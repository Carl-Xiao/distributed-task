package main

import "fmt"

/*
问题描述

请实现一个算法，在不使用【额外数据结构和储存空间】的情况下，翻转一个给定的字符串(可以使用单个过程变量)。

给定一个string，请返回一个string，为翻转后的字符串。保证字符串的长度小于等于5000。
*/

func main() {
	var a = 1
	var b = 2
	a, b = b, a

	fmt.Print(a, b)
	//str := "abcdefg"
	//
	//strs := []rune(str)
	//length := len(strs)
	//
	//for i := 0; i < len(strs)/2; i++ {
	//	strs[i], strs[length-i-1] = strs[length-i-1], strs[i]
	//}
	//
	//fmt.Print(string(strs))
}
