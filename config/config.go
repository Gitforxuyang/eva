package config

import (
	"fmt"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/spf13/viper"
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
type EvaConfig struct {
	name              string
	port              int32
	env               string
	config            map[string]interface{}
	v                 *viper.Viper
	changeNotifyFuncs []ChangeNotify
	apps              map[string]*GRpcClientConfig
}

var (
	config *EvaConfig
)

func Init() {
	if config == nil {
		config = &EvaConfig{}
		config.config = make(map[string]interface{})
		config.changeNotifyFuncs = make([]ChangeNotify, 0)
		config.apps = make(map[string]*GRpcClientConfig)
		v := viper.New()
		v.SetConfigName("config.default")
		v.AddConfigPath("./conf")
		v.SetConfigType("json")
		err := v.ReadInConfig()
		if err != nil {
			panic(err)
		}
		v.BindEnv("ENV")
		env := v.GetString("ENV")
		if env == "" {
			env = "local"
		}
		config.env = env
		v.SetConfigName(fmt.Sprintf("config.%s", env))
		err = v.MergeInConfig()
		if err != nil {
			panic(err)
		}
		config.name = v.GetString("name")
		if config.name == "" {
			panic("配置文件中name不能为空")
		}
		config.port = v.GetInt32("port")
		if config.port == 0 {
			panic("配置文件中port不能为空")
		}
		config.v = v
		err = v.UnmarshalKey("apps", &config.apps)
		if err != nil {
			panic(err)
		}
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
	if err != nil {
		panic(err)
	}
	return c
}

func (m *EvaConfig) GetApp(app string) *GRpcClientConfig {
	c := m.apps[app]
	if c == nil {
		panic(fmt.Sprintf("app：%s配置未找到", app))
	}
	return c
}
