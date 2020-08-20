package template

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

const (
	handler string = `
package handler

import (
	"context"
	"{{.Name}}/proto/{{.Name}}"
)

type HandlerService struct {

}

func (m *HandlerService) Ping(context.Context, *{{.Name}}.Nil) (*{{.Name}}.Nil, error) {
	return &{{.Name}}.Nil{}, nil
}

`
)

func Handler(d Data) {
	f, err := os.Create(path.Join(d.Name, "handler", fmt.Sprintf("index.go")))
	CheckErr(err)
	tmp, err := template.New("test").Parse(handler)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}
