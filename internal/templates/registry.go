package templates

import (
	"text/template"
)

type Registry struct {
	langTemplates map[string][]*template.Template
}

func New() *Registry {
	registry := &Registry{
		langTemplates: make(map[string][]*template.Template),
	}

	return registry
}

func (r *Registry) Load(lang string) []*template.Template {
	if templates, ok := r.langTemplates[lang]; ok {
		return templates
	}

	return nil
}

func (r *Registry) Register(lang string, tpls ...*template.Template) {
	if _, ok := r.langTemplates[lang]; !ok {
		r.langTemplates[lang] = tpls
		return
	}

	r.langTemplates[lang] = append(r.langTemplates[lang], tpls...)
}
