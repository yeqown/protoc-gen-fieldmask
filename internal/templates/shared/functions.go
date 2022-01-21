package shared

import (
	"text/template"

	"github.com/iancoleman/strcase"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func RegisterFunctions(tpl *template.Template, ctx pgsgo.Context) {
	fns := sharedFuncs{Context: ctx}

	tpl.Funcs(template.FuncMap{
		"pkg":       fns.PackageName,
		"snakeCase": fns.snakeCase,
	})
}

type sharedFuncs struct {
	pgsgo.Context
}

func (fns sharedFuncs) snakeCase(name string) string {
	return strcase.ToSnake(name)
}
