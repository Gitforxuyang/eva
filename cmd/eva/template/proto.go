package template

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

const (
	proto string = `
syntax = "proto3";

package {{.Package}};

message Nil {
}
service {{.Name}}Service {
    rpc Ping (Nil) returns (Nil);
}
	`
)

func Proto(d Data) {
	f, err := os.Create(path.Join(d.Name, "proto", fmt.Sprintf("%s.proto", Lcfirst(d.Name))))
	CheckErr(err)
	tmp, err := template.New("test").Parse(proto)
	CheckErr(err)
	err = tmp.Execute(f, map[string]string{"Name": Ucfirst(d.Name), "Package": Lcfirst(d.Name)})
	CheckErr(err)
}
