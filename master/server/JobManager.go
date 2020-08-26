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

		jobObj common.Job
	)
	jobKey = common.JOB_DIR + job.Name
	//序列化字符串存入到etcd
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	if putRespone, err = manager.Kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putRespone.PrevKv != nil {
		if err = json.Unmarshal(putRespone.PrevKv.Value, &jobObj); err != nil {
			common.Error(err.Error())
			return
		}
		oldJob = &jobObj
	}
	return
}

//DeleteManager 删除定时器
func (manager *JobMgr) DeleteManager(name string) (oldJob *common.Job, err error) {
	var (
		jobKey        string
		deleteRespone *clientv3.DeleteResponse
		jobObj        common.Job
	)
	jobKey = common.JOB_DIR + name

	if deleteRespone, err = manager.Kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	if deleteRespone.Deleted > 0 {
		if err = json.Unmarshal(deleteRespone.PrevKvs[0].Value, &jobObj); err != nil {
			common.Error(err.Error())
			err = nil

		}
		oldJob = &jobObj
	}
	return
}

//ListManager 列出所有定时器任务
func (manager *JobMgr) ListManager() (jobList []*common.Job, err error) {
	var (
		jobKey      string
		listRespone *clientv3.GetResponse
		jobObj      common.Job
	)
	jobKey = common.JOB_DIR
	if listRespone, err = manager.Kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	if listRespone.Count > 0 {
		jobList = make([]*common.Job, 0)
		for _, value := range listRespone.Kvs {
			if err = json.Unmarshal(value.Value, &jobObj); err != nil {
				err = nil
				continue
			}
			jobList = append(jobList, &jobObj)
		}
	}
	return
}

// 通知manager杀死任务
func (manager *JobMgr) KillManager(name string) (err error) {
	var (
		jobKey             string
		leaseGrantResponse *clientv3.LeaseGrantResponse
		leaseID            clientv3.LeaseID
	)
	jobKey = common.JOB_KILLER + name

	if leaseGrantResponse, err = manager.Lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	leaseID = leaseGrantResponse.ID

	if _, err = manager.Client.Put(context.TODO(), jobKey, "", clientv3.WithLease(leaseID)); err != nil {
		return
	}
	return
}
