package template

import (
	"html/template"
	"os"
	"path"
	"unicode"
)

const (
	makefile string = `
proto:
	mkdir ./proto/{{.Name}} || true
	protoc --eva_out=plugins=all:./proto/{{.Name}} -I=./proto {{.Name}}.proto
	protoc --go_out=plugins=grpc:./proto/{{.Name}} -I=./proto {{.Name}}.proto

.PHONY: proto
`
	gomod string = `
module {{.Name}}

go 1.12

require (
	github.com/Gitforxuyang/eva v0.0.0-20200812113758-270624ede9a9
	github.com/golang/protobuf v1.4.2
	google.golang.org/grpc v1.26.0
	google.golang.org/protobuf v1.25.0
)

`
	main string = `
package main

import (
	"{{.Name}}/handler"
	"{{.Name}}/proto/{{.Name}}"
	"github.com/Gitforxuyang/eva/server"
	"google.golang.org/grpc"
)

func main(){
	server.Init()
	server.RegisterGRpcService(func(server *grpc.Server) {
		{{.Name}}.Register{{.Service}}ServiceServer(server,&handler.HandlerService{})
	})
	server.Run()
}
`
	git string = `


# Created by https://www.gitignore.io/api/go

### Go ###
# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Test binary, build with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Project-local glide cache, RE: https://github.com/Masterminds/glide/issues/736
.glide/

# End of https://www.gitignore.io/api/go

# Created by https://www.gitignore.io/api/vim

### Vim ###
# swap
[._]*.s[a-v][a-z]
[._]*.sw[a-p]
[._]s[a-v][a-z]
[._]sw[a-p]
# session
Session.vim
# temporary
.netrwhist
*~
# auto-generated tag files
tags

# End of https://www.gitignore.io/api/vim

# Created by https://www.gitignore.io/api/emacs

### Emacs ###
# -*- mode: gitignore; -*-
*~
\#*\#
/.emacs.desktop
/.emacs.desktop.lock
*.elc
auto-save-list
tramp
.\#*

# Org-mode
.org-id-locations
*_archive

# flymake-mode
*_flymake.*

# eshell files
/eshell/history
/eshell/lastdir

# elpa packages
/elpa/

# reftex files
*.rel

# AUCTeX auto folder
/auto/

# cask packages
.cask/
dist/

# Flycheck
flycheck_*.el

# projectiles files
.projectile

# directory configuration
.dir-locals.el

# End of https://www.gitignore.io/api/emacs

bin/*

`
)

func Makefile(d Data) {
	f, err := os.Create(path.Join(d.Name, "Makefile"))
	CheckErr(err)
	tmp, err := template.New("test").Parse(makefile)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}

type Data struct {
	Name    string
	Port    int
	Service string
}

func GoMod(d Data) {
	f, err := os.Create(path.Join(d.Name, "go.mod"))
	CheckErr(err)
	tmp, err := template.New("test").Parse(gomod)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}

func Main(d Data) {
	f, err := os.Create(path.Join(d.Name, "main.go"))
	CheckErr(err)
	tmp, err := template.New("test").Parse(main)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}
func Git(d Data) {
	f, err := os.Create(path.Join(d.Name, ".gitignore"))
	CheckErr(err)
	tmp, err := template.New("test").Parse(git)
	CheckErr(err)
	err = tmp.Execute(f, d)
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
