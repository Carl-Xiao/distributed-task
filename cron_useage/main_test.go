package cron_useage

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"testing"
	"time"
)

//测试单个
func TestCron01(t *testing.T) {
	var (
		expr *cronexpr.Expression
		err  error
	)
	if expr, err = cronexpr.Parse("*/3 * * * * * *"); err != nil {
		fmt.Println(err)
	}
	now := time.Now()
	nextTime := expr.Next(now)
	fmt.Println(nextTime)
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("调度处理")
	})
	time.Sleep(10 * time.Second)
}

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

//测试多个Cron表达式同时执行
func TestCron02(t *testing.T) {
	var (
		expr    *cronexpr.Expression
		taskJob map[string]*CronJob
	)
	taskJob = make(map[string]*CronJob)
	now := time.Now()
	expr = cronexpr.MustParse("*/3 * * * * * *")
	nextTime := expr.Next(now)
	cronJob1 := &CronJob{
		expr:     expr,
		nextTime: nextTime,
	}
	taskJob["job1"] = cronJob1

	expr = cronexpr.MustParse("*/5 * * * * * *")
	nextTime = expr.Next(now)
	cronJob2 := &CronJob{
		expr:     expr,
		nextTime: nextTime,
	}
	taskJob["job2"] = cronJob2

	//启动一个协程
	go func() {
		for {
			now := time.Now()
			for name, job := range taskJob {
				//执行函数
				if job.nextTime.Before(now) || job.nextTime.Equal(now) {
					go func(name string) {
						fmt.Println("执行任务 ：", name)
					}(name)
				}
				job.nextTime = job.expr.Next(now)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	time.Sleep(100 * time.Second)
}
