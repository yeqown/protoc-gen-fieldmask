package templates

import (
	"fmt"
	"text/template"
)

type Registry struct {
	settings map[string]*Lang
}

type Lang struct {
	Name      string               // Name of the language
	Ext       string               // file extension name, such as 'fm.go'
	Templates []*template.Template // templates for this extension
}

func New() *Registry {
	registry := &Registry{
		settings: make(map[string]*Lang, 2),
	}

	return registry
}

func (r *Registry) Load(lang string) *Lang {
	if ls, ok := r.settings[lang]; ok {
		return ls
	}

	return nil
}

func (r *Registry) Register(lang string, s *Lang) {
	if s == nil || s.Ext == "" || len(s.Templates) == 0 {
		fmt.Printf("Invalid lang setting: %#v\n", s)
		return
	}

	if _, ok := r.settings[lang]; ok {
		fmt.Printf("Duplicate lang=%s setting: %#v\n", lang, *s)
		return
	}

	r.settings[lang] = s
}
