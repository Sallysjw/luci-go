// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"text/template"

	"github.com/luci/luci-go/grpc/internal/svctool"
)

var (
	tmpl = template.Must(template.New("").Parse(
		`// Code generated by svcdec; DO NOT EDIT

package {{.PackageName}}

import (
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	{{range .ExtraImports}}
	{{.Name}} "{{.Path}}"{{end}}
)

{{range .Services}}
{{$StructName := .StructName}}
type {{$StructName}} struct {
	// Service is the service to decorate.
	Service {{.Service.TypeName}}
	// Prelude is called in each method before forwarding the call to Service.
	// If Prelude returns an error, it is returned without forwarding the call.
	Prelude func(c context.Context, methodName string, req proto.Message) (context.Context, error)
}

{{range .Methods}}
func (s *{{$StructName}}) {{.Name}}(c context.Context, req {{.InputType}}) ({{.OutputType}}, error) {
	c, err := s.Prelude(c, "{{.Name}}", req)
	if err != nil {
		return nil, err
	}
	return s.Service.{{.Name}}(c, req)
}
{{end}}
{{end}}
`))
)

type (
	templateArgs struct {
		PackageName  string
		Services     []*service
		ExtraImports []svctool.Import
	}

	service struct {
		*svctool.Service
		StructName string
	}
)