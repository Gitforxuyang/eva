package main

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/plugin/http"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"time"
)

func main() {
	config.Init()
	conf := config.GetConfig()
	logger.Init(conf.GetName())
	trace.Init(fmt.Sprintf("%s_%s", conf.GetName(), conf.GetENV()),
		conf.GetTraceConfig().Endpoint, conf.GetTraceConfig().Ratio)
	//client := hello.GetGRpcSayHelloServiceClient()
	//client.Hello(context.TODO(), &hello.String{Name: strconv.Itoa(int(time.Now().Unix()))})
	client := http.GetHttpClient("demo-svc")
	data := make(map[string]interface{})
	resp, err := client.Do(context.TODO(), "/ping1", http.METHOD_GET, http.Headers{}, data)
	if err != nil {
		time.Sleep(time.Second * 3)
		panic(err)
	}
	//fmt.Println(resp)
	fmt.Print(string(resp))
	time.Sleep(time.Second * 3)
}
