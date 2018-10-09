package http

import (
	"bytes"
	"text/template"

	"github.com/fatih/structs"
)

// Template :
type Template struct {
	*template.Template
}

// Render :
func (t *Template) Render(data interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)
	return buf.String(), err
}

// FHasAttr :
func (t *Template) FHasAttr(object interface{}, name string) bool {
	switch value := object.(type) {
	case map[string]interface{}:
		_, ok := value[name]
		return ok
	default:
		s := structs.New(object)
		return t.FHasAttr(s.Map(), name)
	}
}

// NewTemplateWithParser :
func NewTemplateWithParser(name string, parser func(*template.Template) (*template.Template, error)) (*Template, error) {
	var err error
	tmpl := &Template{}
	t := template.New(name).Funcs(template.FuncMap{
		"hasAttr": tmpl.FHasAttr,
	})

	tmpl.Template, err = parser(t)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// NewTemplateFromString :
func NewTemplateFromString(name string, content string) (*Template, error) {
	return NewTemplateWithParser(name, func(t *template.Template) (*template.Template, error) {
		return t.Parse(content)
	})
}

// NewTemplate :
func NewTemplate(path string) (*Template, error) {
	return NewTemplateWithParser(path, func(t *template.Template) (*template.Template, error) {
		return t.ParseFiles(path)
	})
}
