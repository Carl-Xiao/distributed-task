package common

import (
	"encoding/json"
	"fmt"
)

const (
	JOB_DIR    = "/cron/jobs/"
	JOB_KILLER = "/cron/killer/"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
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

func UnPackResponse(byte []byte) (job Job, err error) {
	var obj Job
	if err = json.Unmarshal(byte, &obj); err != nil {
		return
	}
	job = obj
	return
}
