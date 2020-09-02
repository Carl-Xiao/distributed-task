package common

import (
	"encoding/json"
	"fmt"
	"strings"
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
