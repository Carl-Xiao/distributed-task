package server

import (
	"context"
	"encoding/json"
	"github.com/Carl-Xiao/distributed-task/common"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type JobMgr struct {
	Client *clientv3.Client
	Kv     clientv3.KV
	Lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		lease  clientv3.Lease
		kv     clientv3.KV
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

	G_jobMgr = &JobMgr{
		Client: client,
		Kv:     kv,
		Lease:  lease,
	}
	return
}

//SaveManager 保存定时器任务到etcd中
func (manager *JobMgr) SaveManager(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey     string
		jobValue   []byte
		putRespone *clientv3.PutResponse
	)
	jobKey = "/cron/jobs/" + job.Name
	//序列化字符串存入到etcd
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	if putRespone, err = manager.Kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putRespone.PrevKv != nil {
		if err = json.Unmarshal(putRespone.PrevKv.Value, oldJob); err != nil {
			return
		}
	}
	return
}
