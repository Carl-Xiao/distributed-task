package main

import (
	"fmt"
	"strings"
	"sync"
)

/*
交替打印数字和字母
问题描述

使用两个 goroutine 交替打印序列，一个 goroutine 打印数字， 另外一个 goroutine 打印字母， 最终效果如下：

12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
*/
func main() {
	number := make(chan bool)
	letter := make(chan bool)
	character := "ABCDEFGHIJKLMNOPQRSTUVWSYZ"
	wait := sync.WaitGroup{}

	go func() {
		var i = 1
		for {
			select {
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- true
				break
			default:
				break
			}
		}
	}()

	wait.Add(1)
	go func(wait *sync.WaitGroup) {
		i := 0
		for {
			select {
			case <-letter:
				if i >= strings.Count(character, "")-1 {
					wait.Done()
					return
				}
				fmt.Print(character[i : i+1])
				i++
				if i >= strings.Count(character, "") {
					i = 0
				}
				fmt.Print(character[i : i+1])
				i++
				number <- true
				break
			default:
				break
			}
		}
	}(&wait)
	number <- true

	wait.Wait()
}
