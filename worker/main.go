package main

import (
	"fmt"
	"github.com/Carl-Xiao/distributed-task/worker/server"
	"runtime"
	"time"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var err error
	//初始化CPU
	initEnv()
	if err = server.InitJobMgr(); err != nil {
		goto ERR
	}
	for {
		time.Sleep(time.Second * 1)
	}
ERR:
	fmt.Println(err)
}
