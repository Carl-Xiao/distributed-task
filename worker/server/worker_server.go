package server

import (
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

func InitServer() (err error) {
	var handler http.Handler
	mux := http.NewServeMux()

	dir := http.Dir("./master/web")
	handler = http.FileServer(dir)
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", common.App.WORKER_PORT),
		Handler:      mux,
		ReadTimeout:  common.ReadTimeout,
		WriteTimeout: common.WriteTimeout,
	}
	common.Info(fmt.Sprintf(":%d", common.App.WORKER_PORT) + ":启动")

	if err = server.ListenAndServe(); err != nil {
		return
	}
	G_apiServer = &ApiServer{server: server}
	return
}
