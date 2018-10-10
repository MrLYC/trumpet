package http

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/fatih/structs"
	"github.com/oliveagle/jsonpath"
	"github.com/sirupsen/logrus"
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

// FLookup :
func (t *Template) FLookup(object interface{}, path string, defaults interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		result, err := jsonpath.JsonPathLookup(object, path)
		if err != nil {
			logrus.Warnf("look up path[%s] for %+v failed", path, object)
			return defaults
		}
		return result
	default:
		s := structs.New(object)
		return t.FLookup(s.Map(), path, defaults)
	}
}

// NewTemplateWithParser :
func NewTemplateWithParser(name string, parser func(*template.Template) (*template.Template, error)) (*Template, error) {
	var err error
	tmpl := &Template{}
	t := template.New(name).Funcs(template.FuncMap{
		"hasattr": tmpl.FHasAttr,
		"lookup":  tmpl.FLookup,
	})

	tmpl.Template, err = parser(t)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// NewTemplate :
func NewTemplate(path string) (*Template, error) {
	return NewTemplateWithParser(path, func(t *template.Template) (*template.Template, error) {
		file, err := FS.Open(path)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return t.Parse(string(content))
	})
}
