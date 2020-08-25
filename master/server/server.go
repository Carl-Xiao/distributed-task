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

func handlerJobSavae(w http.ResponseWriter, r *http.Request) {

}

func InitServer() (err error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/Job/Save", handlerJobSavae)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", common.PORT),
		Handler:      mux,
		ReadTimeout:  common.ReadTimeout,
		WriteTimeout: common.WriteTimeout,
	}
	common.Info("server init")
	if err = server.ListenAndServe(); err != nil {
		return
	}
	G_apiServer = &ApiServer{server: server}
	return
}
