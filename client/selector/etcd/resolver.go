package etcd

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/registory/etcd"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/naming"
)

//
//func NewEtcdResolver(target resolver.Target, cc resolver.ClientConn) *etcdResolver {
//	return &etcdResolver{target: target, cc: cc, serviceMap: make(map[string]*etcd.ServiceNode, 20)}
//}
//
//type etcdResolver struct {
//	target     resolver.Target
//	cc         resolver.ClientConn
//	serviceMap map[string]*etcd.ServiceNode
//}
//
//func (m *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
//}
//
//func (m *etcdResolver) Close() {
//}
//
//func (m *etcdResolver) Run() error {
//	addrs := make([]resolver.Address, 0, len(m.serviceMap))
//	for _, v := range m.serviceMap {
//		addrs = append(addrs, resolver.Address{Addr: v.Endpoint})
//	}
//	m.cc.UpdateState(resolver.State{Addresses: addrs})
//
//	return nil
//}

// resolver is the implementaion of grpc.naming.Resolver

type resolver struct {
	serviceName string // service name to resolve
}

// NewResolver return resolver with service name
func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		panic("grpclb: no service name provided")
	}
	client := etcd.GetClient()
	return &watcher{re: re, client: client}, nil
}

type watcher struct {
	re            *resolver
	client        *clientv3.Client
	isInitialized bool
}

// Close do nothing
func (w *watcher) Close() {
}

// Next to return the updates
func (w *watcher) Next() ([]*naming.Update, error) {
	// prefix is the etcd prefix/value to watch
	prefix := fmt.Sprintf("%s%s/", etcd.ETCD_SERVICE_PREFIX, w.re.serviceName)
	// check if is initialized
	if !w.isInitialized {
		resp, err := w.client.Get(context.Background(), prefix, clientv3.WithPrefix())
		w.isInitialized = true
		if err == nil {
			addrs := extractAddrs(resp)
			if l := len(addrs); l != 0 {
				updates := make([]*naming.Update, l)
				for i := range addrs {
					updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
				}
				return updates, nil
			}
		}
	}
	// generate etcd Watcher
	rch := w.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			node := etcd.ServiceNode{}
			err := utils.JsonToStruct(string(ev.Kv.Value), &node)
			utils.Must(err)
			switch ev.Type {
			case mvccpb.PUT:
				return []*naming.Update{{Op: naming.Add, Addr: node.Endpoint}}, nil
			case mvccpb.DELETE:
				return []*naming.Update{{Op: naming.Delete, Addr: node.Endpoint}}, nil
			}
		}
	}
	return nil, nil
}
func extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := []string{}
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			node := etcd.ServiceNode{}
			err := utils.JsonToStruct(string(v), &node)
			utils.Must(err)
			addrs = append(addrs, node.Endpoint)
		}
	}
	return addrs
}
