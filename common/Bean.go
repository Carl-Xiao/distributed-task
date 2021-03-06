package common

import (
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

const (
	JOB_DIR    = "/cron/jobs/"
	JOB_KILLER = "/cron/killer/"

	JOB_EVENT_SAVE   = 1
	JOB_EVENT_DELETE = 2
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//任务执行事件
type JobEvent struct {
	EventType int
	*Job
}

func (job Job) ToString() string {
	return fmt.Sprintf("Name:%s ,Command:%s,CronExpr:%s", job.Name, job.Command, job.CronExpr)
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResultResponse(code int, msg string, data interface{}) (byte []byte, err error) {
	var result Result

	result = Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	if byte, err = json.Marshal(result); err != nil {
		return
	}
	return
}

func UnPackResponse(byte []byte) (job *Job, err error) {
	var obj *Job
	if err = json.Unmarshal(byte, &obj); err != nil {
		return
	}
	job = obj
	return
}

//提取JobName
func ExtractJobName(name string) string {
	return strings.TrimPrefix(name, JOB_DIR)
}

// 封装JobEvent事件
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 任务调度器的执行事件
type JonSchedulerPlan struct {
	Job      *Job
	Express  *cronexpr.Expression
	NextTime time.Time
}

func BuildJobSchedulePlan(jobEvent *JobEvent) (plan *JonSchedulerPlan, err error) {
	var (
		cronExpr *cronexpr.Expression
		nextTime time.Time
	)
	if cronExpr, err = cronexpr.Parse(jobEvent.CronExpr); err != nil {
		Error(err.Error())
		return
	}
	//下次执行时间
	now := time.Now()
	nextTime = cronExpr.Next(now)
	plan = &JonSchedulerPlan{
		Job:      jobEvent.Job,
		Express:  cronExpr,
		NextTime: nextTime,
	}
	return
}

type JobExecuteInfo struct {
	Job      *Job
	PlanTime time.Time
	RealTime time.Time
}

func BuildJobExecuteInfo(plan *JonSchedulerPlan) (execute *JobExecuteInfo) {
	execute = &JobExecuteInfo{
		Job:      plan.Job,
		PlanTime: plan.NextTime,
		RealTime: time.Now(),
	}
	return
}
