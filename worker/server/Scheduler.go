package server

import (
	"fmt"
	"github.com/Carl-Xiao/distributed-task/common"
	"time"
)

var (
	G_scheduler *Scheduler
)

type Scheduler struct {
	Events        chan *common.JobEvent
	JobMap        map[string]*common.JonSchedulerPlan
	JobExecuteMap map[string]*common.JobExecuteInfo
}

func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		Events:        make(chan *common.JobEvent, 100),
		JobMap:        make(map[string]*common.JonSchedulerPlan, 100),
		JobExecuteMap: make(map[string]*common.JobExecuteInfo, 10),
	}
	err = G_jobMgr.JobWatch()
	go G_scheduler.LoopEvent()
	return
}

//推送任务到调度器中
func (scheduler *Scheduler) PushEvent(event *common.JobEvent) {
	scheduler.Events <- event
}

func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		plan *common.JonSchedulerPlan
		err  error
		ok   bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		common.Info("保存事件")
		if plan, err = common.BuildJobSchedulePlan(jobEvent); err == nil {
			scheduler.JobMap[plan.Job.Name] = plan
		}
	case common.JOB_EVENT_DELETE:
		common.Info("删除事件")
		if _, ok = scheduler.JobMap[jobEvent.Name]; ok {
			delete(scheduler.JobMap, jobEvent.Name)
		}
	}
}

// 启动任务
func (scheduler *Scheduler) StartJob(plan *common.JonSchedulerPlan) {
	var (
		execute *common.JobExecuteInfo
		exist   bool
	)
	//记录任务执行状态，如果当前任务正在执行,跳过执行任务
	if _, exist = scheduler.JobExecuteMap[plan.Job.Name]; exist {
		return
	}

	//调度执行任务
	execute = common.BuildJobExecuteInfo(plan)
	scheduler.JobExecuteMap[plan.Job.Name] = execute

	//执行任务
	G_executor.Execute(execute)
}

//计算任务调度器 寻找合理的休息时间,避免内存消耗
func (scheduler *Scheduler) CalculationScheduler() (timeAfter time.Duration) {
	//1 遍历所有任务
	var (
		job        *common.JonSchedulerPlan
		now        time.Time
		recentTime *time.Time
	)
	//没有任务就休息1s
	if len(scheduler.JobMap) == 0 {
		timeAfter = time.Second * 1
		return
	}
	now = time.Now()
	for _, job = range scheduler.JobMap {
		if job.NextTime.Before(now) || job.NextTime.Equal(now) {
			//执行当前任务,不管之前是否执行过与否
			scheduler.StartJob(job)
			//更新下次执行时间
			job.NextTime = job.Express.Next(now)
		}
		if recentTime == nil || recentTime.Before(job.NextTime) {
			recentTime = &job.NextTime
		}
	}
	timeAfter = (*recentTime).Sub(now)

	return
}

//循环处理事件
func (scheduler *Scheduler) LoopEvent() {
	var (
		jobEvent  *common.JobEvent
		sleepTime time.Duration
		timeTr    *time.Timer
	)
	sleepTime = scheduler.CalculationScheduler()
	timeTr = time.NewTimer(sleepTime)
	fmt.Println(sleepTime)
	for {
		fmt.Println("调度一次")
		select {
		case jobEvent = <-scheduler.Events:
			scheduler.handleJobEvent(jobEvent)
		case <-timeTr.C:
			common.Info("时间已到")
		}
		sleepTime = scheduler.CalculationScheduler()
		fmt.Println(sleepTime)
		timeTr.Reset(sleepTime)
	}
}
