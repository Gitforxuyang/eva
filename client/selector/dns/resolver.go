package dns

import (
	"google.golang.org/grpc/naming"
	"os"
	"os/signal"
	"syscall"
)

//func NewDNSResolver(target resolver.Target, cc resolver.ClientConn) *dnsResolver {
//	r := &dnsResolver{target: target, cc: cc, serviceMap: make(map[string]*dnsServiceNode, 20)}
//	r.serviceMap["default"] = &dnsServiceNode{ip: target.Endpoint}
//	return r
//}
//
//type dnsResolver struct {
//	target     resolver.Target
//	cc         resolver.ClientConn
//	serviceMap map[string]*dnsServiceNode
//}
//
//func (m *dnsResolver) ResolveNow(options resolver.ResolveNowOptions) {
//}
//
//func (m *dnsResolver) Close() {
//}
//
//func (m *dnsResolver) Run() error {
//	addrs := make([]resolver.Address, 0, len(m.serviceMap))
//	for _, v := range m.serviceMap {
//		addrs = append(addrs, resolver.Address{Addr: v.ip})
//	}
//	m.cc.UpdateState(resolver.State{Addresses: addrs})
//	return nil
//}
//
//type dnsServiceNode struct {
//	ip string
//}

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
	return &watcher{re: re}, nil
}

type watcher struct {
	re            *resolver
	isInitialized bool
}

// Close do nothing
func (w *watcher) Close() {
}

// Next to return the updates
func (w *watcher) Next() ([]*naming.Update, error) {
	if !w.isInitialized {
		updates := make([]*naming.Update, 1, 1)
		for i := 0; i < 1; i++ {
			updates[i] = &naming.Update{Op: naming.Add, Addr: w.re.serviceName}
		}
		return updates, nil
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGUSR2)
	<-sig
	return nil, nil
}
