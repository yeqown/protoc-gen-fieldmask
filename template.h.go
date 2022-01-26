package main

import (
	_ "embed"
	"text/template"

	"github.com/yeqown/protoc-gen-fieldmask/internal/templates"
	"github.com/yeqown/protoc-gen-fieldmask/internal/templates/shared"

	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func tplRegistryFactory(ctx pgsgo.Context) *templates.Registry {
	registry := templates.New()

	registry.Register("go", makeTemplatesForGo(ctx)) // lang=go

	return registry
}

var (
	//go:embed internal/templates/go/file.tpl
	_tplFileGO []byte
	//go:embed internal/templates/go/message.tpl
	_tplMessageGO []byte
	//go:embed internal/templates/go/fm.tpl
	_tplFieldMaskGO []byte
	//go:embed internal/templates/go/fm.in.tpl
	_tplFieldMaskInGO []byte
	//go:embed internal/templates/go/fm.out.tpl
	_tplFieldMaskOutGO []byte
)

func makeTemplatesForGo(ctx pgsgo.Context) *template.Template {
	tpl := template.New("go")

	shared.RegisterFunctions(tpl, ctx)
	template.Must(tpl.Parse(string(_tplFileGO)))
	template.Must(tpl.New("message").Parse(string(_tplMessageGO)))
	template.Must(tpl.New("fm").Parse(string(_tplFieldMaskGO)))
	template.Must(tpl.New("fm.in").Parse(string(_tplFieldMaskInGO)))
	template.Must(tpl.New("fm.out").Parse(string(_tplFieldMaskOutGO)))

	return tpl
}
