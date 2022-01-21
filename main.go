package main

import (
	_ "embed"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yeqown/protoc-gen-fieldmask/internal/module"
	"github.com/yeqown/protoc-gen-fieldmask/internal/templates"
	"github.com/yeqown/protoc-gen-fieldmask/internal/templates/shared"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(pgs.DebugEnv("DEBUG_PGFM"), pgs.SupportedFeatures(&optional)).
		RegisterModule(module.FieldMask(tplRegistryFactory)).
		RegisterPostProcessor(
			pgsgo.GoFmt(),
		).
		Render()
}

var (
	//go:embed internal/templates/go/message.tpl
	_tplMessageGO string
	//go:embed internal/templates/go/file.tpl
	_tplFileGO string
)

func tplRegistryFactory(ctx pgsgo.Context) *templates.Registry {
	registry := templates.New()

	registry.Register("go", makeTemplatesForGo(ctx)) // lang=go

	return registry
}

func makeTemplatesForGo(ctx pgsgo.Context) *template.Template {
	tpl := template.New("file.tpl")

	shared.RegisterFunctions(tpl, ctx)
	template.Must(tpl.Parse(_tplFileGO))
	template.Must(tpl.New("message.tpl").Parse(_tplMessageGO))

	return tpl
}
