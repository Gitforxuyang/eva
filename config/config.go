package config

import (
	"fmt"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/spf13/viper"
	"strings"
)

//配置发送变更时的通知
type ChangeNotify func(config map[string]interface{})

type TraceConfig struct {
	Endpoint string
	Ratio    float64
}

type GRpcClientConfig struct {
	Mode     string //dns etcd
	Endpoint string
	Timeout  int64 //请求超时时间
}
type HttpClientConfig struct {
	Endpoint string
	Timeout  int64
	MaxConn  int //最大连接数
}
type LogConfig struct {
	Server     bool   //服务端日志是否打印
	GRpcClient bool   //grpc客户端日志
	HttpClient bool   //http客户端日志
	Level      string //日志打印级别
}
type EvaConfig struct {
	name              string
	port              int32
	env               string
	config            map[string]interface{}
	v                 *viper.Viper
	changeNotifyFuncs []ChangeNotify
	grpc              map[string]*GRpcClientConfig
	http              map[string]*HttpClientConfig
	log               *LogConfig
}

var (
	config *EvaConfig
)

func Init() {
	if config == nil {
		config = &EvaConfig{}
		config.config = make(map[string]interface{})
		config.changeNotifyFuncs = make([]ChangeNotify, 0)
		config.grpc = make(map[string]*GRpcClientConfig)
		config.log = &LogConfig{Server: false, GRpcClient: false, HttpClient: false, Level: "LOG"}
		v := viper.New()
		v.SetConfigName("config.default")
		v.AddConfigPath("./conf")
		v.SetConfigType("json")
		err := v.ReadInConfig()
		utils.Must(err)
		v.BindEnv("ENV")
		env := v.GetString("ENV")
		if env == "" {
			env = "local"
		}
		config.env = env
		v.SetConfigName(fmt.Sprintf("config.%s", env))
		err = v.MergeInConfig()
		utils.Must(err)
		config.name = v.GetString("name")
		if config.name == "" {
			panic("配置文件中name不能为空")
		}
		config.port = v.GetInt32("port")
		if config.port == 0 {
			panic("配置文件中port不能为空")
		}
		config.v = v
		err = v.UnmarshalKey("grpc", &config.grpc)
		utils.Must(err)
		err = v.UnmarshalKey("log", &config.log)
		utils.Must(err)
		err = v.UnmarshalKey("http", &config.http)
		utils.Must(err)
	}
}

func GetConfig() *EvaConfig {
	if config == nil {
		Init()
	}
	return config
}

func (m *EvaConfig) RegisterNotify(f ChangeNotify) {
	m.changeNotifyFuncs = append(m.changeNotifyFuncs, f)
}

func (m *EvaConfig) changeNotify(config map[string]interface{}) {
	for _, v := range m.changeNotifyFuncs {
		v(config)
	}
}
func (m *EvaConfig) GetName() string {
	return m.name
}
func (m *EvaConfig) GetPort() int32 {
	return m.port
}
func (m *EvaConfig) GetENV() string {
	return m.env
}

func (m *EvaConfig) GetTraceConfig() TraceConfig {
	c := TraceConfig{}
	if utils.IsNil(m.v.Get("trace")) {
		panic("trace设置不能为空")
	}
	err := m.v.UnmarshalKey("trace", &c)
	utils.Must(err)
	return c
}

func (m *EvaConfig) GetGRpc(app string) *GRpcClientConfig {
	c := m.grpc[strings.ToLower(app)]
	if c == nil {
		panic(fmt.Sprintf("grpc：%s配置未找到", app))
	}
	return c
}
func (m *EvaConfig) GetHttp(http string) *HttpClientConfig {
	c := m.http[strings.ToLower(http)]
	if c == nil {
		panic(fmt.Sprintf("http：%s配置未找到", http))
	}
	return c
}

func (m *EvaConfig) GetLogConfig() *LogConfig {
	if m.log == nil {
		panic(fmt.Sprintf("log配置未找到"))
	}
	return m.log
}