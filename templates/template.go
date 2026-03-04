package templates

import (
	_ "embed"
	"text/template"

	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yeqown/protoc-gen-fieldmask/templates/shared"
)

func RegistryFactory(ctx pgsgo.Context) *Registry {
	registry := New()

	registry.Register("go", makeTemplatesForGo(ctx)) // lang=go

	return registry
}

var (
	//go:embed go/file.tpl
	_tplFileGO []byte
	//go:embed go/request_mask.tpl
	_tplRequestMaskGO []byte
	//go:embed go/response_mask.tpl
	_tplResponseMaskGO []byte
	//go:embed go/marked_checker.tpl
	_tplMarkedCheckerGO []byte
)

func makeTemplatesForGo(ctx pgsgo.Context) *Lang {
	tpl := template.New("go")

	shared.RegisterFunctions(tpl, ctx)
	template.Must(tpl.Parse(string(_tplFileGO)))
	template.Must(tpl.New("request_mask").Parse(string(_tplRequestMaskGO)))
	template.Must(tpl.New("response_mask").Parse(string(_tplResponseMaskGO)))
	template.Must(tpl.New("marked_checker").Parse(string(_tplMarkedCheckerGO)))

	return &Lang{
		Name:      "go",
		Ext:       ".fm.go",
		Templates: []*template.Template{tpl},
	}
}
