package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
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
