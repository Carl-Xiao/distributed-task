package server

import (
	"context"
	"fmt"
	"github.com/Carl-Xiao/distributed-task/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

type JobMgr struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		lease   clientv3.Lease
		kv      clientv3.KV
		watcher clientv3.Watcher
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   common.EndPoint, //集群地址
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		common.Error(err.Error())
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		Client:  client,
		Kv:      kv,
		Lease:   lease,
		Watcher: watcher,
	}
	err = G_jobMgr.JobWatch()
	return
}

//JobWatch 监听定时任务
func (manager *JobMgr) JobWatch() (err error) {
	var (
		response      *clientv3.GetResponse
		job           *common.Job
		jobEvent      *common.JobEvent
		revision      int64
		whathChan     clientv3.WatchChan
		watchResponse clientv3.WatchResponse
		event         *clientv3.Event
		jobName       string
	)
	if response, err = manager.Kv.Get(context.TODO(), common.JOB_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	//获取所有Keys
	for _, value := range response.Kvs {
		//正常情况
		if job, err = common.UnPackResponse(value.Value); err == nil {
			//TODO 开启任务调度协程
			common.Info(job.ToString())
		}
	}
	//启动一个监听协程
	go func() {
		//获取当前的历史版本
		revision = response.Header.Revision + 1
		//开启watcher
		whathChan = manager.Watcher.Watch(context.TODO(), common.JOB_DIR, clientv3.WithRev(revision), clientv3.WithPrefix())
		for watchResponse = range whathChan {
			for _, event = range watchResponse.Events {
				switch event.Type {
				case mvccpb.PUT:
					if job, err = common.UnPackResponse(event.Kv.Value); err != nil {
						common.Error(err.Error())
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(event.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				}
				//TODO 推送JobEvent事件
				fmt.Println(jobEvent.EventType)
			}
		}
	}()
	return
}
