package server

import "github.com/Carl-Xiao/distributed-task/common"

var (
	G_scheduler *Scheduler
)

type Scheduler struct {
	Events chan *common.JobEvent
}

func InitScheduler() {
	G_scheduler = &Scheduler{
		Events: make(chan *common.JobEvent, 100),
	}
}

//推送任务到调度器中
func (scheduler *Scheduler) PushEvent(event *common.JobEvent) {
	scheduler.Events <- event
}

func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:

	case common.JOB_EVENT_DELETE:

	}

}

//循环处理事件
func (scheduler *Scheduler) LoopEvent() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case jobEvent = <-scheduler.Events:
			scheduler.handleJobEvent(jobEvent)
		}

	}
}
