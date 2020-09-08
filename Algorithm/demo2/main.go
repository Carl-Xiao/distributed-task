package main

func main() {
	c, quit := make(chan int), make(chan int)

	go func() {
		c <- 1    // c通道的数据没有被其他goroutine读取走，堵塞当前goroutine
		quit <- 0 // quit始终没有办法写入数据
	}()
	<-quit // quit 等待数据的写
}
