package server

import (
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
