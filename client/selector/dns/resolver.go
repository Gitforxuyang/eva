package dns

import (
	"google.golang.org/grpc/resolver"
)

func NewDNSResolver(target resolver.Target, cc resolver.ClientConn) *dnsResolver {
	r := &dnsResolver{target: target, cc: cc, serviceMap: make(map[string]*dnsServiceNode, 20)}
	r.serviceMap["default"] = &dnsServiceNode{ip: target.Endpoint}
	return r
}

type dnsResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	serviceMap map[string]*dnsServiceNode
}

func (m *dnsResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (m *dnsResolver) Close() {
}

func (m *dnsResolver) Run() error {
	addrs := make([]resolver.Address, 0, len(m.serviceMap))
	for _, v := range m.serviceMap {
		addrs = append(addrs, resolver.Address{Addr: v.ip})
	}
	m.cc.UpdateState(resolver.State{Addresses: addrs})
	return nil
}

type dnsServiceNode struct {
	ip string
}
