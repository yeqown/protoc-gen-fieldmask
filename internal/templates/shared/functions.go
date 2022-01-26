package shared

import (
	"errors"
	"text/template"

	"github.com/iancoleman/strcase"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/descriptorpb"
)

func RegisterFunctions(tpl *template.Template, ctx pgsgo.Context) {
	fns := sharedFuncs{Context: ctx}

	tpl.Funcs(template.FuncMap{
		"pkg":       fns.PackageName,
		"snakeCase": fns.snakeCase,
		"isMessage": fns.isMessage,
		"dict":      fns.dict,
	})
}

type sharedFuncs struct {
	pgsgo.Context
}

func (fns sharedFuncs) snakeCase(name string) string {
	return strcase.ToSnake(name)
}

func (fns sharedFuncs) isMessage(f pgs.Field) bool {
	f.Message()

	return f.Descriptor().GetType().Number() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Number()
}

func (fns sharedFuncs) dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
