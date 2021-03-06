package eva

import (
	"fmt"
	"github.com/Gitforxuyang/eva/cmd/protoc-gen-eva/generator"
	"github.com/Gitforxuyang/eva/util/utils"
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
	//if m.mode==MODE_EVA||m.mode==MODE_SERVER{
	//	m.genMessage(files)
	//}
	for _, svc := range files.Service {
		if m.mode == MODE_EVA {
			m.genServerCode(svc, files)
			m.genClientCode(svc)
		}
		if m.mode == MODE_SERVER {
			m.genServerCode(svc, files)
		}
		if m.mode == MODE_CLIENT {
			m.genClientCode(svc)
		}
	}
}
//func (m *EvaPlugin) genMessage(file *generator.FileDescriptor) {
//	for _, message := range file.MessageType {
//		fmt.Println(message.)
//	}
//}

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
	m.g.P(`"github.com/Gitforxuyang/eva/registory/etcd"`)
	m.g.P(`"google.golang.org/grpc"`)
	m.g.P(`"google.golang.org/grpc/keepalive"`)
	m.g.P(`"time"`)
	m.g.P(")")
}

var (
	typeNameMap map[int32]string= map[int32]string{
		1:  "Double",
		2:  "Float",
		3:  "Int64",
		4:  "UInt64",
		5:  "Int32",
		6:  "FIXED64",
		7:  "FIXED32",
		8:  "Bool",
		9:  "String",
		10: "GROUP",
		11: "MESSAGE",
		12: "BYTES",
		13: "UInt32",
		14: "ENUM",
		15: "SFIXED32",
		16: "SFIXED64",
		17: "SINT32",
		18: "SINT64",
		}
)
func (m *EvaPlugin) genServerCode(svc *descriptor.ServiceDescriptorProto, file *generator.FileDescriptor) {
	m.g.P("//获取服务的描述信息")
	m.g.P("func GetServerDesc() *etcd.Service{")
	m.g.P(fmt.Sprintf("messageMap:=make(map[string]map[string]string,%d)", len(file.MessageType)))
	for _, message := range file.MessageType {
		m.g.P(fmt.Sprintf("message%s:=make(map[string]string)", message.GetName()))
		for _, field := range message.Field {
			typeName:=field.GetTypeName()
			//如果typeName不为空则说明不是一个基本类型，则需要继续往下递归
			if typeName==""{
				typeName=typeNameMap[int32(field.GetType())]
				if field.GetLabel()== descriptor.FieldDescriptorProto_LABEL_REPEATED{
					typeName="[]"+typeName
				}
			}else{
				if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					desc := m.g.ObjectNamed(field.GetTypeName())
					if d, ok := desc.(*generator.Descriptor); ok && d.GetOptions().GetMapEntry() {
						// Figure out the Go types and tags for the key and value types.
						keyField, valField := d.Field[0], d.Field[1]
						keyType, _ := m.g.GoType(d, keyField)
						valType, _ := m.g.GoType(d, valField)

						// We don't use stars, except for message-typed values.
						// Message and enum types are the only two possibly foreign types used in maps,
						// so record their use. They are not permitted as map keys.
						keyType = strings.TrimPrefix(keyType, "*")
						switch *valField.Type {
						case descriptor.FieldDescriptorProto_TYPE_ENUM:
							valType = strings.TrimPrefix(valType, "*")
							m.g.RecordTypeUse(valField.GetTypeName())
						case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
							m.g.RecordTypeUse(valField.GetTypeName())
						default:
							valType = strings.TrimPrefix(valType, "*")
						}
						typeName = fmt.Sprintf("map[%s]%s", keyType, valType)
					}
				}
			}
			//fmt.Println(field.get)
			//typeName=field.GetOptions().Ctype.String()
				//fmt.Println(111)
				//fmt.Println(field.GetExtendee())
				//fmt.Println(111)
			m.g.P(fmt.Sprintf(`message%s["%s"]="%s"`, message.GetName(), field.GetName(),typeName))
		}
		m.g.P(fmt.Sprintf(`messageMap["%s"]=message%s`,message.GetName(),message.GetName()))
	}
	m.g.P("service:=new(etcd.Service)")
	m.g.P(fmt.Sprintf(`service.Name="%s"`, *svc.Name))
	m.g.P(fmt.Sprintf(`service.Package="%s"`, utils.StrFirstToLower(svc.GetName())))
	m.g.P(fmt.Sprintf(`service.AppId="%s"`, utils.StrFirstToLower(svc.GetName())))
	m.g.P(fmt.Sprintf("service.Methods=make(map[string]etcd.Method,%d)", len(svc.Method)))

	//
	//m.g.P(fmt.Sprintf("type I%sServer interface {",*svc.Name))
	for _, v := range svc.Method {
		var input string = *v.InputType
		inputArr := strings.Split(input, ".")
		inputArr = inputArr[:copy(inputArr, inputArr[2:])]
		req := strings.Join(inputArr, "_")

		var output string = *v.OutputType
		outputArr := strings.Split(output, ".")
		outputArr = outputArr[:copy(outputArr, outputArr[2:])]
		res := strings.Join(outputArr, "_")
		m.g.P(fmt.Sprintf("method%s:=etcd.Method{}", v.GetName()))
		m.g.P(fmt.Sprintf(`method%s.Req=messageMap["%s"]`, v.GetName(), req))
		m.g.P(fmt.Sprintf(`method%s.Resp=messageMap["%s"]`, v.GetName(), res))
		m.g.P(fmt.Sprintf(`service.Methods["%s"]=method%s`, v.GetName(), v.GetName()))
	}
	m.g.P("return service")
	m.g.P("}")
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
	//m.g.P("grpc.WithBlock(),")
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
