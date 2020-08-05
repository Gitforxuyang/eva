package etcd

import (
	"google.golang.org/grpc/resolver"
)

func NewEtcdResolver(target resolver.Target, cc resolver.ClientConn) *etcdResolver {
	return &etcdResolver{target: target, cc: cc, serviceMap: make(map[string]*etcdServiceNode, 20)}
}

type etcdResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	serviceMap map[string]*etcdServiceNode
}

func (m *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (m *etcdResolver) Close() {
}

func (m *etcdResolver) Run() error {
	addrs := make([]resolver.Address, 0, len(m.serviceMap))
	for _, v := range m.serviceMap {
		addrs = append(addrs, resolver.Address{Addr: v.ip})
	}
	m.cc.UpdateState(resolver.State{Addresses: addrs})

	return nil
}

type etcdServiceNode struct {
	ip string
}
