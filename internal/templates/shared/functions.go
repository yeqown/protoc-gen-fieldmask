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
		"pkg":              fns.PackageName,
		"snakeCase":        fns.snakeCase,
		"isMessage":        fns.isMessage,
		"dict":             fns.dict,
		"fieldOptions":     fns.fieldOptions,
		"nestedMessageName": fns.nestedMessageName,
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

func (fns sharedFuncs) fieldOptions(field pgs.Field) map[string]interface{} {
	result := map[string]interface{}{
		"Ignore": false,
		"Nested": false,
	}

	// Check if field has field mask options
	fieldDesc := field.Descriptor()
	if fieldDesc == nil {
		return result
	}

	extOpts := fieldDesc.GetOptions()
	if extOpts == nil {
		return result
	}

	// For now, return default values. In a full implementation,
	// this would parse the actual field options from the extension
	return result
}

func (fns sharedFuncs) nestedMessageName(field pgs.Field) string {
	if !fns.isMessage(field) {
		return ""
	}

	if field.Type().Embed() != nil {
		return field.Type().Embed().Name().String()
	}

	return ""
}
