package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Print("1")
	select {
	case <-ctx.Done():
		fmt.Fprintf(os.Stdout, "Task cancelded! \n")
		return
	default:
	}

	fmt.Print("2")
	cancel()
	time.Sleep(time.Second * 10)
}
