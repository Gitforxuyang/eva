package selector

import (
	"github.com/Gitforxuyang/eva/client/selector/dns"
	"github.com/Gitforxuyang/eva/client/selector/etcd"
	"google.golang.org/grpc/naming"
)

type CustomResolver struct {
	scheme string
}

func (m *CustomResolver) GetResolver(target string) naming.Resolver {
	switch m.scheme {
	case "dns":
		return dns.NewResolver(target)
	case "etcd":
		return etcd.NewResolver(target)
	default:
		panic("不存在的协议")
	}
}

func NewCustomResolverBuilder(scheme string) *CustomResolver {
	return &CustomResolver{scheme: scheme}
}
