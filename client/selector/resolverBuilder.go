package selector

import (
	"github.com/Gitforxuyang/eva/client/selector/dns"
	"github.com/Gitforxuyang/eva/client/selector/etcd"
	"google.golang.org/grpc/resolver"
)

type customResolver struct {
	name string
}

func (m *customResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var r CustomResolver
	switch m.name {
	case "etcd":
		r = etcd.NewEtcdResolver(target, cc)
	case "dns":
		r = dns.NewDNSResolver(target, cc)
	default:
		panic("不存在的resover类型")
	}
	err := r.Run()
	if err != nil {
		panic(err)
	}
	return r, nil
}

func (m *customResolver) Scheme() string {
	return m.name
}

func NewCustomResolverBuilder(name string) resolver.Builder {
	return &customResolver{name: name}
}

type CustomResolver interface {
	ResolveNow(options resolver.ResolveNowOptions)
	Close()
	Run() error
}
