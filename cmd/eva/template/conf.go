package template

import (
	"html/template"
	"os"
	"path"
)

const (
	def string = `
{
  "name": "{{.Name}}",
  "port": {{.Port}},
  "version": "0.0.1",
  "mongo": {
    "node": {
      "url": "mongodb://192.168.3.3:27017/demo",
      "maxPoolSize": 20,
      "minPoolSize": 20
    }
  },
  "redis": {
    "node": {
      "addr": "192.168.3.3:6379",
      "password": "",
      "DB": 0,
      "PoolSize": 20,
      "MinIdleConns": 20,
      "DialTimeout": 5,
      "ReadTimeout": 3,
      "WriteTimeout": 3
    }
  },
  "trace": {
    "endpoint": "http://192.168.3.23:14268/api/traces",
    "ratio": 1
  },
  "grpc": {
    "SayHelloService": {
      "endpoint": ":50001",
      "timeout": 5,
      "mode": "dns"
    }
  },
  "http": {
    "demo-svc": {
      "endpoint": "http://127.0.0.1:8080",
      "timeout": 5,
      "maxConn": 10
    }
  },
  "log": {
    "server": true,
    "grpcClient": true,
    "httpClient": true
  }
}`
	loc  string = `{}`
	dev  string = `{}`
	test string = `{}`
	prod string = `{}`
)

func Conf(d Data) {
	f, err := os.Create(path.Join(d.Name, "conf", "config.default.json"))
	CheckErr(err)
	tmp, err := template.New("test").Parse(def)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)

	f, err = os.Create(path.Join(d.Name, "conf", "config.local.json"))
	CheckErr(err)
	tmp, err = template.New("test").Parse(loc)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)

	f, err = os.Create(path.Join(d.Name, "conf", "config.dev.json"))
	CheckErr(err)
	tmp, err = template.New("test").Parse(dev)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)

	f, err = os.Create(path.Join(d.Name, "conf", "config.test.json"))
	CheckErr(err)
	tmp, err = template.New("test").Parse(test)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)

	f, err = os.Create(path.Join(d.Name, "conf", "config.prod.json"))
	CheckErr(err)
	tmp, err = template.New("test").Parse(prod)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}
