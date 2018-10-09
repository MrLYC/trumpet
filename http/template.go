package http

import (
	"bytes"
	"text/template"
)

// Template :
type Template struct {
	*template.Template
}

// init :
func (t *Template) init() {
	t.Template.Funcs(template.FuncMap{})
}

// Render :
func (t *Template) Render(data interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)
	return buf.String(), err
}

// NewTemplateWithParser :
func NewTemplateWithParser(parser func() (*template.Template, error)) (*Template, error) {
	t, err := parser()
	if err != nil {
		return nil, err
	}
	tmpl := &Template{
		Template: t,
	}
	tmpl.init()
	return tmpl, nil
}

// NewTemplateFromString :
func NewTemplateFromString(name string, content string) (*Template, error) {
	return NewTemplateWithParser(func() (*template.Template, error) {
		return template.New(name).Parse(content)
	})
}

// NewTemplate :
func NewTemplate(path string) (*Template, error) {
	return NewTemplateWithParser(func() (*template.Template, error) {
		return template.New(path).ParseFiles(path)
	})
}
