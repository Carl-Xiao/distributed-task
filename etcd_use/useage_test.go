package etcd_use

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

func getClient() *clientv3.Client {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return nil
	}

	return client
}

//CRUD操作
func TestGet(t *testing.T) {
	var (
		getResponse *clientv3.GetResponse
		err         error
	)

	client := getClient()
	con := context.Background()

	kv := clientv3.NewKV(client)
	if _, err = kv.Put(con, "test1", "CD"); err != nil {
		fmt.Println(err)
		return
	}
	getResponse, err = kv.Get(con, "test1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(getResponse.Kvs[0].Value))
}

//申请租约
func TestLease(t *testing.T) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		leaseId        clientv3.LeaseID
		kv             clientv3.KV
		getResp        *clientv3.GetResponse
	)
	client := getClient()
	// 申请一个租约
	lease := clientv3.NewLease(client)
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = leaseGrantResp.ID

	// 获得kv API子集
	kv = clientv3.NewKV(client)

	if _, err = kv.Put(context.TODO(), "/test2", "HAHA", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	for {
		if getResp, err = kv.Get(context.TODO(), "/test2"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}

//监听器
func TestWatch(t *testing.T) {
	client := getClient()
	var (
		err                error
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchRespChan      clientv3.WatchChan
	)

	kv = clientv3.NewKV(client)

	// 模拟KV的变化
	go func() {
		for {
			_, err = kv.Put(context.TODO(), "/school/class/students", "helios1")
			_, err = kv.Delete(context.TODO(), "/school/class/students")
			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/school/class/students"); err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}

	// 获得当前revision
	watchStartRevision = getResp.Header.Revision + 1
	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)
	fmt.Println("从该版本向后监听:", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	watchRespChan = watcher.Watch(ctx, "/school/class/students", clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
