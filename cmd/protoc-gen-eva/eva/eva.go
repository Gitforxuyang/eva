package eva

import (
	"fmt"
	"github.com/Gitforxuyang/eva/cmd/protoc-gen-eva/generator"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

const (
	MODE_EVA    = "all"
	MODE_SERVER = "server"
	MODE_CLIENT = "client"
)

type EvaPlugin struct {
	g    *generator.Generator
	mode string //模式  eva默认模式 server-只生成server client-只生成client
}

//func (m *EvaPlugin) GenerateImports(file *generator.FileDescriptor, imports map[generator.GoImportPath]generator.GoPackageName) {
//	m.genImportCode(file)
//}

func init() {
	generator.RegisterPlugin(&EvaPlugin{mode: MODE_EVA})
	generator.RegisterPlugin(&EvaPlugin{mode: MODE_SERVER})
	generator.RegisterPlugin(&EvaPlugin{mode: MODE_CLIENT})
}
func (m *EvaPlugin) Name() string {
	return m.mode
}
func (m *EvaPlugin) Init(g *generator.Generator) {
	m.g = g
}
func (m *EvaPlugin) GenerateImports(file *generator.FileDescriptor, imports map[generator.GoImportPath]generator.GoPackageName) {
	if len(file.Service) == 0 {
		return
	}
	m.genImportCode(file)
}

func (m *EvaPlugin) Generate(files *generator.FileDescriptor) {
	if len(files.Service) == 0 {
		return
	}
	for _, svc := range files.Service {
		if m.mode == MODE_EVA {
			m.genServerCode(svc)
			m.genClientCode(svc)
		}
		if m.mode == MODE_SERVER {
			m.genServerCode(svc)
		}
		if m.mode == MODE_CLIENT {
			m.genServerCode(svc)
		}
	}
}

func (m *EvaPlugin) genImportCode(file *generator.FileDescriptor) {
	m.g.P("import (")
	m.g.P(`"context"`)
	m.g.P(`"fmt"`)
	m.g.P(`"github.com/Gitforxuyang/eva/client/selector"`)
	m.g.P(`"github.com/Gitforxuyang/eva/config"`)
	m.g.P(`"github.com/Gitforxuyang/eva/server"`)
	m.g.P(`"github.com/Gitforxuyang/eva/util/utils"`)
	m.g.P(`trace2 "github.com/Gitforxuyang/eva/util/trace"`)
	m.g.P(`"github.com/Gitforxuyang/eva/wrapper/catch"`)
	m.g.P(`"github.com/Gitforxuyang/eva/wrapper/log"`)
	m.g.P(`"github.com/Gitforxuyang/eva/wrapper/trace"`)
	m.g.P(`"google.golang.org/grpc"`)
	m.g.P(`"google.golang.org/grpc/keepalive"`)
	m.g.P(`"time"`)
	m.g.P(")")
}

func (m *EvaPlugin) genServerCode(svc *descriptor.ServiceDescriptorProto) {
	//
	//m.g.P(fmt.Sprintf("type I%sServer interface {",*svc.Name))
	//for _,v:=range svc.Method{
	//	var input string=*v.InputType
	//	inputArr:=strings.Split(input,".")
	//	inputArr=inputArr[:copy(inputArr,inputArr[2:])]
	//
	//	var output string=*v.OutputType
	//	outputArr:=strings.Split(output,".")
	//	outputArr=outputArr[:copy(outputArr,outputArr[2:])]
	//	m.g.P(fmt.Sprintf("%s(ctx context.Context,req *%s) (resp *%s,err error)",
	//		*v.Name,
	//		strings.Join(inputArr,"_"),
	//		strings.Join(outputArr,"_"),
	//	))
	//}
	//m.g.P(fmt.Sprintf("}"))
	//m.g.P("")
}

func (m *EvaPlugin) genClientCode(svc *descriptor.ServiceDescriptorProto) {

	//生成interface
	m.g.P(fmt.Sprintf("type GRpc%sClient interface {", *svc.Name))
	for _, v := range svc.Method {
		var input string = *v.InputType
		inputArr := strings.Split(input, ".")
		inputArr = inputArr[:copy(inputArr, inputArr[2:])]

		var output string = *v.OutputType
		outputArr := strings.Split(output, ".")
		outputArr = outputArr[:copy(outputArr, outputArr[2:])]
		m.g.P(fmt.Sprintf("%s(ctx context.Context,req *%s) (resp *%s,err error)",
			*v.Name,
			strings.Join(inputArr, "_"),
			strings.Join(outputArr, "_"),
		))
	}
	m.g.P("}")

	//生成struct
	m.g.P(fmt.Sprintf("type grpc%sClient struct {", *svc.Name))
	m.g.P(fmt.Sprintf("client %sClient", *svc.Name))
	m.g.P(fmt.Sprintf("}"))

	for _, v := range svc.Method {
		var input string = *v.InputType
		inputArr := strings.Split(input, ".")
		inputArr = inputArr[:copy(inputArr, inputArr[2:])]

		var output string = *v.OutputType
		outputArr := strings.Split(output, ".")
		outputArr = outputArr[:copy(outputArr, outputArr[2:])]

		m.g.P(fmt.Sprintf("func (m *grpc%sClient) %s(ctx context.Context, req *%s) (resp *%s, err error){",
			*svc.Name,
			*v.Name,
			strings.Join(inputArr, "_"),
			strings.Join(outputArr, "_"),
		))
		m.g.P(fmt.Sprintf("resp, err = m.client.%s(ctx, req)", *v.Name))
		m.g.P("return resp, err")
		m.g.P("}")
	}
	//生成get方法
	m.g.P(fmt.Sprintf("func GetGRpc%sClient() GRpc%sClient {", *svc.Name, *svc.Name))
	m.g.P("tracer := trace2.GetTracer()")
	m.g.P(fmt.Sprintf(`grpcClientConfig := config.GetConfig().GetGRpc("%s")`, *svc.Name))
	m.g.P(`conn, err := grpc.Dial(fmt.Sprintf("%s", grpcClientConfig.Endpoint),`)
	m.g.P("grpc.WithInsecure(),")
	m.g.P("grpc.WithBlock(),")
	//m.g.P("grpc.WithBalancerName(roundrobin.Name),")
	m.g.P("grpc.WithBalancer(grpc.RoundRobin(selector.NewCustomResolverBuilder(grpcClientConfig.Mode).GetResolver(grpcClientConfig.Endpoint))),")
	m.g.P("grpc.WithKeepaliveParams(")
	m.g.P("keepalive.ClientParameters{")
	m.g.P("Time:                time.Second * 10,")
	m.g.P("Timeout:             time.Second * 1,")
	m.g.P("PermitWithoutStream: true,")
	m.g.P("}),")
	m.g.P("grpc.WithChainUnaryInterceptor(")
	m.g.P("trace.NewClientWrapper(tracer),")
	m.g.P("log.NewClientWrapper(),")
	m.g.P("catch.NewClientWrapper(grpcClientConfig.Timeout),")
	m.g.P("),")
	m.g.P(")")
	m.g.P(fmt.Sprintf("c := &grpc%sClient{}", *svc.Name))
	m.g.P(fmt.Sprintf("c.client = New%sClient(conn)", *svc.Name))
	m.g.P("utils.Must(err)")
	m.g.P("server.RegisterShutdownFunc(func() {")
	m.g.P("conn.Close()")
	m.g.P("})")
	m.g.P("return c")
	m.g.P("}")
}
