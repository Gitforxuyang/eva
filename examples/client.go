package main

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/registory/etcd"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"strconv"
	"time"
)

func main() {
	config.Init()
	conf := config.GetConfig()
	logger.Init(conf.GetName())
	trace.Init(fmt.Sprintf("%s_%s", conf.GetName(), conf.GetENV()),
		conf.GetTraceConfig().Endpoint, conf.GetTraceConfig().Ratio)
	etcd.Init()
	client := hello.GetGRpcSayHelloServiceClient()
	client.Hello(context.TODO(), &hello.String{Name: strconv.Itoa(int(time.Now().Unix()))})
	//client := http.GetHttpClient("demo-svc")
	//data := make(map[string]interface{})
	//data["name"] = "123"
	//resp, err := client.Do(context.TODO(), "/ping1?age=123", http.METHOD_GET, http.Headers{}, data)
	//if err != nil {
	//	time.Sleep(time.Second * 3)
	//	panic(err)
	//}
	////fmt.Println(resp)
	//fmt.Print(string(resp))
	time.Sleep(time.Second * 3)
}
