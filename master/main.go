package main

import (
	"fmt"
	"github.com/Carl-Xiao/distributed-task/common"
	"github.com/Carl-Xiao/distributed-task/master/server"
	"runtime"
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
		fmt.Println(err)
		goto ERR
	}

	//初始化服务
	if err = server.InitServer(); err != nil {
		goto ERR
	}
ERR:
	fmt.Println(err)
}
