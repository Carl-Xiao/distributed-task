package server

import (
	"encoding/json"
	"fmt"
	"github.com/Carl-Xiao/distributed-task/common"
	"net/http"
)

var (
	G_apiServer *ApiServer
)

type ApiServer struct {
	server *http.Server
}

func handlerJobSave(w http.ResponseWriter, r *http.Request) {
	common.Info("save jonb")
	var (
		err    error
		jobStr string
		job    common.Job
		oldJob *common.Job
		result []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	jobStr = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(jobStr), &job); err != nil {
		common.Error(err.Error())
		goto ERR
	}

	if oldJob, err = G_jobMgr.SaveManager(&job); err != nil {
		goto ERR
	}

	if result, err = common.BuildResultResponse(0, "SUCCESS", oldJob); err == nil {
		_, _ = w.Write(result)
	}
	return
ERR:
	if result, err = common.BuildResultResponse(-1, "FAIL", err.Error()); err != nil {
		_, _ = w.Write(result)
	}
}

func handlerJobDelete(w http.ResponseWriter, r *http.Request) {
	common.Info("delete jonb")
	var (
		err    error
		name   string
		oldJob *common.Job
		result []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.PostForm.Get("name")
	common.Info("name:" + name)
	if oldJob, err = G_jobMgr.DeleteManager(name); err != nil {
		goto ERR
	}

	if result, err = common.BuildResultResponse(0, "SUCCESS", oldJob); err == nil {
		_, _ = w.Write(result)
	}
	return
ERR:
	if result, err = common.BuildResultResponse(-1, "FAIL", nil); err != nil {
		_, _ = w.Write(result)
	}
}

func handlerJobList(w http.ResponseWriter, r *http.Request) {
	common.Info("Job List")
	var (
		err    error
		oldJob []*common.Job
		result []byte
	)

	if oldJob, err = G_jobMgr.ListManager(); err != nil {
		goto ERR
	}
	if result, err = common.BuildResultResponse(0, "SUCCESS", oldJob); err == nil {
		_, _ = w.Write(result)
	}
	return
ERR:
	if result, err = common.BuildResultResponse(-1, "FAIL", nil); err != nil {
		_, _ = w.Write(result)
	}
}

func handlerJobKiller(w http.ResponseWriter, r *http.Request) {
	common.Info("kill jonb")
	var (
		err    error
		name   string
		result []byte
	)
	name = r.PostForm.Get("name")
	if err = G_jobMgr.KillManager(name); err != nil {
		goto ERR
	}
	if result, err = common.BuildResultResponse(0, "SUCCESS", nil); err == nil {
		_, _ = w.Write(result)
	}
	return
ERR:
	if result, err = common.BuildResultResponse(-1, "FAIL", nil); err != nil {
		_, _ = w.Write(result)
	}
}

func InitServer() (err error) {
	var handler http.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/Job/save", handlerJobSave)
	mux.HandleFunc("/Job/delete", handlerJobDelete)
	mux.HandleFunc("/Job/list", handlerJobList)
	mux.HandleFunc("/Job/kill", handlerJobKiller)

	dir := http.Dir("./master/web")
	handler = http.FileServer(dir)
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", common.PORT),
		Handler:      mux,
		ReadTimeout:  common.ReadTimeout,
		WriteTimeout: common.WriteTimeout,
	}
	common.Info(fmt.Sprintf(":%d", common.PORT) + ":启动")

	if err = server.ListenAndServe(); err != nil {
		return
	}
	G_apiServer = &ApiServer{server: server}
	return
}
