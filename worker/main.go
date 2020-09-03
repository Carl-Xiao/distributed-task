package main

import (
	"fmt"
	"github.com/Carl-Xiao/distributed-task/common"
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

	//初始化配置文件
	if err = common.InitBase(); err != nil {
		goto ERR
	}

	if err = server.InitJobMgr(); err != nil {
		goto ERR
	}

	//初始化协程调度器
	if err = server.InitScheduler(); err != nil {
		goto ERR
	}

	//初始化执行器
	if err = server.InitExecutor(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(time.Second * 1)
	}
ERR:
	fmt.Println(err)
}
