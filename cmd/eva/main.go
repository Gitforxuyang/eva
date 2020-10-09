package main

import (
	"flag"
	"fmt"
	"github.com/Gitforxuyang/eva/cmd/eva/template"
	"os"
	"path"
	"regexp"
	"unicode"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	tmp string = "eva -name=xxx -port=xxx"
)

func main() {
	name := flag.String("name", "", "app name")
	port := flag.Int("port", 0, "app port")
	flag.Parse()
	if *name == "" {
		panic(fmt.Sprintf("name不能为空 \n示例：%s", tmp))
	}
	if *port == 0 {
		panic(fmt.Sprintf("port不能为0 \n示例：%s", tmp))
	}
	match, _ := regexp.MatchString("^[a-zA-Z]+$", *name)
	if !match {
		panic("app name只能是大小写字母")
	}
	//首字母不允许大写
	if unicode.IsUpper([]rune(*name)[0]) {
		panic("首字母必须小写")
	}
	if isExist(*name) {
		panic("期望创建的服务已存在")
	}
	//创建文件夹
	err := os.MkdirAll(path.Join(*name), 0777)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(path.Join(*name, "app", "service"), 0777)
	os.MkdirAll(path.Join(*name, "app", "assembler"), 0777)
	os.MkdirAll(path.Join(*name, "conf"), 0777)
	os.MkdirAll(path.Join(*name, "handler"), 0777)
	os.MkdirAll(path.Join(*name, "domain", "entity"), 0777)
	os.MkdirAll(path.Join(*name, "domain", "event"), 0777)
	os.MkdirAll(path.Join(*name, "domain", "repo"), 0777)
	os.MkdirAll(path.Join(*name, "domain", "service"), 0777)
	os.MkdirAll(path.Join(*name, "domain", "util"), 0777)
	os.MkdirAll(path.Join(*name, "infra", "err"), 0777)
	os.MkdirAll(path.Join(*name, "infra", "helper"), 0777)
	os.MkdirAll(path.Join(*name, "infra", "repo"), 0777)
	os.MkdirAll(path.Join(*name, "infra", "util"), 0777)
	os.MkdirAll(path.Join(*name, "proto"), 0777)
	d := template.Data{Name: *name, Port: *port, Service: template.Ucfirst(*name)}
	template.Makefile(d)
	template.GoMod(d)
	template.Git(d)
	template.Main(d)
	template.Dockerfile(d)
	template.Conf(d)
	template.Proto(d)
	template.Handler(d)

}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
