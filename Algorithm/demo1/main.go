package main

import (
	"fmt"
	"time"
)

/*
交替打印数字和字母
问题描述

使用两个 goroutine 交替打印序列，一个 goroutine 打印数字， 另外一个 goroutine 打印字母， 最终效果如下：

12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
*/
func main() {
	//使用管道chan
	//1 需要使用两个chan，控制两个不同的输出
	number := make(chan bool)
	letter := make(chan bool)
	//2 打印数字 和字母

	number <- true

	//2.1 打印数字
	go func() {
		i := 0
		for {
			select {
			case <-number:
				i++
				fmt.Print(i)
				i++
				fmt.Print(i)
				letter <- true
				break
			default:
				break
			}
		}
	}()
	//2.2 打印字母
	go func() {
		str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		for {
			select {
			case <-letter:
				if len(str)-1 < 0 {
					break
				}
				temp := str[0:1]
				fmt.Print(temp)
				str = str[1:]
				temp = str[0:1]
				fmt.Print(temp)
				str = str[1:]
				number <- true
				break
			default:
				break
			}
		}
	}()

	//3 组合打印队列
	//  在两个chan依次传递数据
	time.Sleep(time.Second * 3)
}
