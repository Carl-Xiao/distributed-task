package server

import "github.com/Carl-Xiao/distributed-task/common"

var (
	G_executor *Executor
)

type Executor struct {
}

func (executor *Executor) Execute(info *common.JobExecuteInfo) {
	go func() {
		//TODO shell直接执行bash命令即可

	}()
}

// 初始化任务调度器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
