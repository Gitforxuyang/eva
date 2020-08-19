package etcd

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

const (
	//服务注册前缀
	ETCD_SERVICE_PREFIX = "/eva/service/"
)

type ServiceNode struct {
	Name     string `json:"name"`     //服务名
	Id       string `json:"id"`       //节点id 服务启动时随机生成的唯一id
	Endpoint string `json:"endpoint"` //服务的访问地址
}

var (
	client   *clientv3.Client
	etcdOnce sync.Once
)

func GetClient() *clientv3.Client {
	if client == nil {
		panic("client不存在")
	}
	return client
}
func Init() {
	etcdOnce.Do(func() {
		var err error
		client, err = clientv3.New(clientv3.Config{
			Endpoints:   config.GetConfig().GetEtcd(),
			DialTimeout: time.Second * 3,
		})
		utils.Must(err)
	})
}

func Registry(name string, endpoint string, id string) {
	node := ServiceNode{Name: name, Endpoint: endpoint, Id: id}
	lease := clientv3.NewLease(client)
	leaseGrantResp, err := lease.Grant(context.TODO(), 5)
	utils.Must(err)
	leaseId := leaseGrantResp.ID
	_, err = lease.KeepAlive(context.TODO(), leaseId)
	utils.Must(err)
	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(context.TODO(),
		fmt.Sprintf("%s%s/%s", ETCD_SERVICE_PREFIX, name, id),
		utils.StructToJson(node), clientv3.WithLease(leaseId))
	utils.Must(err)
	logger.GetLogger().Info(context.TODO(), "服务注册成功", logger.Fields{
		"revision": putResp.Header.Revision,
		"key":      fmt.Sprintf("%s%s/%s", ETCD_SERVICE_PREFIX, name, id),
		"node":     node,
	})
}

func UnRegistry(name string, id string) {
	kv := clientv3.NewKV(client)
	putResp, err := kv.Delete(context.TODO(),
		fmt.Sprintf("%s%s/%s", ETCD_SERVICE_PREFIX, name, id))
	utils.Must(err)
	logger.GetLogger().Info(context.TODO(), "服务注销成功", logger.Fields{
		"revision": putResp.Header.Revision,
		"key":      fmt.Sprintf("%s%s/%s", ETCD_SERVICE_PREFIX, name, id),
	})
}
